package rig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/steveyegge/gastown/internal/config"
	"github.com/steveyegge/gastown/internal/git"
	"github.com/steveyegge/gastown/internal/templates"
)

// Common errors
var (
	ErrRigNotFound = errors.New("rig not found")
	ErrRigExists   = errors.New("rig already exists")
)

// RigConfig represents the rig-level configuration (config.json at rig root).
type RigConfig struct {
	Type      string       `json:"type"`       // "rig"
	Version   int          `json:"version"`    // schema version
	Name      string       `json:"name"`       // rig name
	GitURL    string       `json:"git_url"`    // repository URL
	CreatedAt time.Time    `json:"created_at"` // when rig was created
	Beads     *BeadsConfig `json:"beads,omitempty"`
}

// BeadsConfig represents beads configuration for the rig.
type BeadsConfig struct {
	Prefix     string `json:"prefix"`                // issue prefix (e.g., "gt")
	SyncRemote string `json:"sync_remote,omitempty"` // git remote for bd sync
}

// CurrentRigConfigVersion is the current schema version.
const CurrentRigConfigVersion = 1

// Manager handles rig discovery, loading, and creation.
type Manager struct {
	townRoot string
	config   *config.RigsConfig
	git      *git.Git
}

// NewManager creates a new rig manager.
func NewManager(townRoot string, rigsConfig *config.RigsConfig, g *git.Git) *Manager {
	return &Manager{
		townRoot: townRoot,
		config:   rigsConfig,
		git:      g,
	}
}

// DiscoverRigs returns all rigs registered in the workspace.
func (m *Manager) DiscoverRigs() ([]*Rig, error) {
	var rigs []*Rig

	for name, entry := range m.config.Rigs {
		rig, err := m.loadRig(name, entry)
		if err != nil {
			// Log error but continue with other rigs
			continue
		}
		rigs = append(rigs, rig)
	}

	return rigs, nil
}

// GetRig returns a specific rig by name.
func (m *Manager) GetRig(name string) (*Rig, error) {
	entry, ok := m.config.Rigs[name]
	if !ok {
		return nil, ErrRigNotFound
	}

	return m.loadRig(name, entry)
}

// RigExists checks if a rig is registered.
func (m *Manager) RigExists(name string) bool {
	_, ok := m.config.Rigs[name]
	return ok
}

// loadRig loads rig details from the filesystem.
func (m *Manager) loadRig(name string, entry config.RigEntry) (*Rig, error) {
	rigPath := filepath.Join(m.townRoot, name)

	// Verify directory exists
	info, err := os.Stat(rigPath)
	if err != nil {
		return nil, fmt.Errorf("rig directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", rigPath)
	}

	rig := &Rig{
		Name:   name,
		Path:   rigPath,
		GitURL: entry.GitURL,
		Config: entry.BeadsConfig,
	}

	// Scan for polecats
	polecatsDir := filepath.Join(rigPath, "polecats")
	if entries, err := os.ReadDir(polecatsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				rig.Polecats = append(rig.Polecats, e.Name())
			}
		}
	}

	// Scan for crew workers
	crewDir := filepath.Join(rigPath, "crew")
	if entries, err := os.ReadDir(crewDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				rig.Crew = append(rig.Crew, e.Name())
			}
		}
	}

	// Check for witness
	witnessPath := filepath.Join(rigPath, "witness", "rig")
	if _, err := os.Stat(witnessPath); err == nil {
		rig.HasWitness = true
	}

	// Check for refinery
	refineryPath := filepath.Join(rigPath, "refinery", "rig")
	if _, err := os.Stat(refineryPath); err == nil {
		rig.HasRefinery = true
	}

	// Check for mayor clone
	mayorPath := filepath.Join(rigPath, "mayor", "rig")
	if _, err := os.Stat(mayorPath); err == nil {
		rig.HasMayor = true
	}

	return rig, nil
}

// AddRigOptions configures rig creation.
type AddRigOptions struct {
	Name        string // Rig name (directory name)
	GitURL      string // Repository URL
	BeadsPrefix string // Beads issue prefix (defaults to derived from name)
	CrewName    string // Default crew workspace name (defaults to "main")
}

// AddRig creates a new rig as a container with clones for each agent.
// The rig structure is:
//
//	<name>/                    # Container (NOT a git clone)
//	├── config.json            # Rig configuration
//	├── .beads/                # Rig-level issue tracking
//	├── refinery/rig/          # Canonical main clone
//	├── mayor/rig/             # Mayor's working clone
//	├── witness/               # Witness agent (no clone)
//	├── polecats/              # Worker directories (empty)
//	└── crew/<crew>/           # Default human workspace
func (m *Manager) AddRig(opts AddRigOptions) (*Rig, error) {
	if m.RigExists(opts.Name) {
		return nil, ErrRigExists
	}

	rigPath := filepath.Join(m.townRoot, opts.Name)

	// Check if directory already exists
	if _, err := os.Stat(rigPath); err == nil {
		return nil, fmt.Errorf("directory already exists: %s", rigPath)
	}

	// Derive defaults
	if opts.BeadsPrefix == "" {
		opts.BeadsPrefix = deriveBeadsPrefix(opts.Name)
	}
	if opts.CrewName == "" {
		opts.CrewName = "main"
	}

	// Create container directory
	if err := os.MkdirAll(rigPath, 0755); err != nil {
		return nil, fmt.Errorf("creating rig directory: %w", err)
	}

	// Track cleanup on failure
	cleanup := func() { os.RemoveAll(rigPath) }
	success := false
	defer func() {
		if !success {
			cleanup()
		}
	}()

	// Create rig config
	rigConfig := &RigConfig{
		Type:      "rig",
		Version:   CurrentRigConfigVersion,
		Name:      opts.Name,
		GitURL:    opts.GitURL,
		CreatedAt: time.Now(),
		Beads: &BeadsConfig{
			Prefix: opts.BeadsPrefix,
		},
	}
	if err := m.saveRigConfig(rigPath, rigConfig); err != nil {
		return nil, fmt.Errorf("saving rig config: %w", err)
	}

	// Clone repository for refinery (canonical main)
	refineryRigPath := filepath.Join(rigPath, "refinery", "rig")
	if err := os.MkdirAll(filepath.Dir(refineryRigPath), 0755); err != nil {
		return nil, fmt.Errorf("creating refinery dir: %w", err)
	}
	if err := m.git.Clone(opts.GitURL, refineryRigPath); err != nil {
		return nil, fmt.Errorf("cloning for refinery: %w", err)
	}
	// Create refinery CLAUDE.md (overrides any from cloned repo)
	if err := m.createRoleCLAUDEmd(refineryRigPath, "refinery", opts.Name, ""); err != nil {
		return nil, fmt.Errorf("creating refinery CLAUDE.md: %w", err)
	}

	// Clone repository for mayor
	mayorRigPath := filepath.Join(rigPath, "mayor", "rig")
	if err := os.MkdirAll(filepath.Dir(mayorRigPath), 0755); err != nil {
		return nil, fmt.Errorf("creating mayor dir: %w", err)
	}
	if err := m.git.Clone(opts.GitURL, mayorRigPath); err != nil {
		return nil, fmt.Errorf("cloning for mayor: %w", err)
	}
	// Create mayor CLAUDE.md (overrides any from cloned repo)
	if err := m.createRoleCLAUDEmd(mayorRigPath, "mayor", opts.Name, ""); err != nil {
		return nil, fmt.Errorf("creating mayor CLAUDE.md: %w", err)
	}

	// Clone repository for default crew workspace
	crewPath := filepath.Join(rigPath, "crew", opts.CrewName)
	if err := os.MkdirAll(filepath.Dir(crewPath), 0755); err != nil {
		return nil, fmt.Errorf("creating crew dir: %w", err)
	}
	if err := m.git.Clone(opts.GitURL, crewPath); err != nil {
		return nil, fmt.Errorf("cloning for crew: %w", err)
	}
	// Create crew CLAUDE.md (overrides any from cloned repo)
	if err := m.createRoleCLAUDEmd(crewPath, "crew", opts.Name, opts.CrewName); err != nil {
		return nil, fmt.Errorf("creating crew CLAUDE.md: %w", err)
	}

	// Create witness directory (no clone needed)
	witnessPath := filepath.Join(rigPath, "witness")
	if err := os.MkdirAll(witnessPath, 0755); err != nil {
		return nil, fmt.Errorf("creating witness dir: %w", err)
	}

	// Create polecats directory (empty)
	polecatsPath := filepath.Join(rigPath, "polecats")
	if err := os.MkdirAll(polecatsPath, 0755); err != nil {
		return nil, fmt.Errorf("creating polecats dir: %w", err)
	}

	// Initialize agent state files
	if err := m.initAgentStates(rigPath); err != nil {
		return nil, fmt.Errorf("initializing agent states: %w", err)
	}

	// Initialize beads at rig level
	if err := m.initBeads(rigPath, opts.BeadsPrefix); err != nil {
		return nil, fmt.Errorf("initializing beads: %w", err)
	}

	// Register in town config
	m.config.Rigs[opts.Name] = config.RigEntry{
		GitURL:  opts.GitURL,
		AddedAt: time.Now(),
		BeadsConfig: &config.BeadsConfig{
			Prefix: opts.BeadsPrefix,
		},
	}

	success = true
	return m.loadRig(opts.Name, m.config.Rigs[opts.Name])
}

// saveRigConfig writes the rig configuration to config.json.
func (m *Manager) saveRigConfig(rigPath string, cfg *RigConfig) error {
	configPath := filepath.Join(rigPath, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// initAgentStates creates initial state.json files for agents.
func (m *Manager) initAgentStates(rigPath string) error {
	agents := []struct {
		path string
		role string
	}{
		{filepath.Join(rigPath, "refinery", "state.json"), "refinery"},
		{filepath.Join(rigPath, "witness", "state.json"), "witness"},
		{filepath.Join(rigPath, "mayor", "state.json"), "mayor"},
	}

	for _, agent := range agents {
		state := &config.AgentState{
			Role:       agent.role,
			LastActive: time.Now(),
		}
		data, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(agent.path, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

// initBeads initializes the beads database at rig level.
func (m *Manager) initBeads(rigPath, prefix string) error {
	beadsDir := filepath.Join(rigPath, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		return err
	}

	// Run bd init if available, with --no-agents to skip AGENTS.md creation
	cmd := exec.Command("bd", "init", "--prefix", prefix, "--no-agents")
	cmd.Dir = rigPath
	if err := cmd.Run(); err != nil {
		// bd might not be installed or --no-agents not supported, create minimal structure
		configPath := filepath.Join(beadsDir, "config.yaml")
		configContent := fmt.Sprintf("prefix: %s\n", prefix)
		if writeErr := os.WriteFile(configPath, []byte(configContent), 0644); writeErr != nil {
			return writeErr
		}
	}
	return nil
}

// deriveBeadsPrefix generates a beads prefix from a rig name.
// Examples: "gastown" -> "gt", "my-project" -> "mp", "foo" -> "foo"
func deriveBeadsPrefix(name string) string {
	// Remove common suffixes
	name = strings.TrimSuffix(name, "-py")
	name = strings.TrimSuffix(name, "-go")

	// Split on hyphens/underscores
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '-' || r == '_'
	})

	if len(parts) >= 2 {
		// Take first letter of each part: "gas-town" -> "gt"
		prefix := ""
		for _, p := range parts {
			if len(p) > 0 {
				prefix += string(p[0])
			}
		}
		return strings.ToLower(prefix)
	}

	// Single word: use first 2-3 chars
	if len(name) <= 3 {
		return strings.ToLower(name)
	}
	return strings.ToLower(name[:2])
}

// RemoveRig unregisters a rig (does not delete files).
func (m *Manager) RemoveRig(name string) error {
	if !m.RigExists(name) {
		return ErrRigNotFound
	}

	delete(m.config.Rigs, name)
	return nil
}

// ListRigNames returns the names of all registered rigs.
func (m *Manager) ListRigNames() []string {
	names := make([]string, 0, len(m.config.Rigs))
	for name := range m.config.Rigs {
		names = append(names, name)
	}
	return names
}

// createRoleCLAUDEmd creates a CLAUDE.md file with role-specific context.
// This ensures each workspace (crew, refinery, mayor) gets the correct prompting,
// overriding any CLAUDE.md that may exist in the cloned repository.
func (m *Manager) createRoleCLAUDEmd(workspacePath string, role string, rigName string, workerName string) error {
	tmpl, err := templates.New()
	if err != nil {
		return err
	}

	data := templates.RoleData{
		Role:     role,
		RigName:  rigName,
		TownRoot: m.townRoot,
		WorkDir:  workspacePath,
		Polecat:  workerName, // Used for crew member name as well
	}

	content, err := tmpl.RenderRole(role, data)
	if err != nil {
		return err
	}

	claudePath := filepath.Join(workspacePath, "CLAUDE.md")
	return os.WriteFile(claudePath, []byte(content), 0644)
}
