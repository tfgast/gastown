# Molecular Chemistry: Work Composition in Gas Town

> *"Work is fractal. Money is crystallized labor. Blockchain was the mechanism
> searching for its purpose."*

Gas Town is a **work composition and execution engine**. This document describes
the chemistry-inspired system for expressing, instantiating, and executing work
at any scale - from single tasks to massive polymers that can grind through
weekends of autonomous operation.

## The Steam Engine Metaphor

Gas Town is an engine. Engines do work and generate steam.

```
Claude = Fire (the energy source)
Claude Code = Steam Engine (harnesses the fire)
Gas Town = Steam Train (coordinates engines on tracks)
Beads = Railroad Tracks (the persistent ledger of work)
```

In our chemistry:
- **Proto molecules** are the fuel (templates that define workflows)
- **Mols** are liquid work (flowing, dynamic, adapting as steps complete)
- **Wisps** are the steam (transient execution traces that rise and dissipate)
- **Digests** are the distillate (condensed permanent records of completed work)

## The Three Phases of Matter

Work in Gas Town exists in three phases, following the states of matter:

| Phase | Name | State | Storage | Behavior |
|-------|------|-------|---------|----------|
| **Solid** | Proto | Frozen template | `.beads/` (template label) | Crystallized, immutable, reusable |
| **Liquid** | Mol | Flowing instance | `.beads/` | Dynamic, adapting, persistent |
| **Vapor** | Wisp | Ephemeral trace | `.beads-wisp/` | Transient, dissipates, operational |

### Proto (Solid Phase)

Protos are **frozen workflow patterns** - crystallized templates that encode
reusable work structures. They're the "molds" from which instances are cast.

```markdown
## Molecule: engineer-in-box
Full workflow from design to merge.

## Step: design
Think carefully about architecture.
Needs: (none)

## Step: implement
Write the code.
Needs: design

## Step: test
Write and run tests.
Needs: implement

## Step: submit
Submit for merge.
Needs: test
```

**Properties:**
- Immutable once defined (frozen)
- Named (e.g., `mol-engineer-in-box`, `mol-code-review`)
- Stored in permanent beads with `template` label
- Can be composed into larger protos (polymers)

### Mol (Liquid Phase)

Mols are **flowing work instances** - live executions that adapt as steps
complete, status changes, and work evolves.

**Properties:**
- Identified by head bead ID (e.g., `bd-abc123`)
- Dynamic - steps transition through states
- Persistent - survives sessions, crashes, context compaction
- Auditable - full history in beads

### Wisp (Vapor Phase)

Wisps are **ephemeral execution traces** - the steam that rises during work
and dissipates when done.

**Properties:**
- Stored in `.beads-wisp/` (gitignored, never synced)
- Single-cycle lifetime
- Either evaporates (burn) or condenses to digest (squash)
- Used for patrol cycles, operational loops, routine work

## Phase Transition Operators

Work transitions between phases through specific operators:

```
                    ┌─────────────────┐
                    │     PROTO       │
                    │    (solid)      │
                    └────────┬────────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
            pour           wisp         distill
              │              │              ↑
              ▼              ▼              │
      ┌───────────────┐  ┌───────────────┐ │
      │      MOL      │  │     WISP      │ │
      │   (liquid)    │  │    (vapor)    │ │
      └───────┬───────┘  └───────┬───────┘ │
              │                  │         │
           squash             squash       │
              │                  │         │
              ▼                  ▼         │
      ┌───────────────┐  ┌───────────────┐ │
      │    DIGEST     │  │  (evaporates) │ │
      │  (condensed)  │  │    or burn    │ │
      └───────────────┘  └───────────────┘ │
              │                            │
              └────────────────────────────┘
                       (experience crystallizes)
```

### Pour: Solid → Liquid

**Pour** instantiates a proto into a persistent mol - like pouring molten metal
into a mold to create a solid casting that will flow through the workflow.

```bash
bd pour mol-feature                    # Create mol from proto
bd pour mol-feature --var version=1.0  # With variable substitution
```

**Use cases:**
- Feature work that spans sessions
- Important work needing audit trail
- Anything you might need to reference later

### Wisp: Solid → Vapor (Sublimation)

**Wisp** instantiates a proto directly into vapor - sublimation that skips
the liquid phase for ephemeral, operational work.

```bash
bd wisp mol-patrol                     # Create wisp from proto
bd wisp mol-health-check               # Ephemeral operational task
```

**Use cases:**
- Patrol cycles (deacon, witness)
- Health checks and monitoring
- One-shot orchestration runs
- Routine operations with no audit value

### Squash: Liquid/Vapor → Condensed

**Squash** condenses work into a permanent digest - the outcome crystallizes
while the execution trace compresses or evaporates.

```bash
bd mol squash bd-abc123                              # Squash mol to digest
bd mol squash bd-abc123 --summary="Completed auth"   # With summary
```

**For mols:** Creates digest in permanent beads, preserves full outcome
**For wisps:** Creates digest, deletes wisp (vapor condenses to residue)

### Burn: Vapor → Nothing

**Burn** discards a wisp without creating a digest - the steam simply
evaporates with no residue.

```bash
bd mol burn wisp-123                   # Discard without digest
```

**Use cases:**
- Routine patrol cycles with nothing notable
- Failed attempts that don't need recording
- Test runs

### Distill: Liquid → Solid (Crystallization)

**Distill** extracts a reusable proto from an existing mol or epic - the
reverse of pour. Experience crystallizes into a template.

```bash
bd mol distill bd-abc123 --as "Release Workflow"
bd mol distill bd-abc123 --var feature=auth --var version=1.0
```

**Process:**
1. Analyze the existing work structure
2. Extract the pattern (steps, dependencies)
3. Replace concrete values with `{{variable}}` placeholders
4. Crystallize as a new proto

**Use cases:**
- Team develops a good workflow organically, wants to reuse it
- Capture tribal knowledge as executable templates
- Create starting points for similar future work

## The Polymorphic Bond Operator

**Bond** is Gas Town's polymorphic combiner - it handles different operand
types with different phase behaviors:

### Bond: Proto + Proto → Compound Proto

Two solid templates fuse into a larger solid template.

```bash
bd mol bond mol-review mol-deploy --as "Review and Deploy"
```

Creates a compound proto that includes both workflows. The result is a
reusable template (solid phase).

### Bond: Proto + Mol → Spawn + Attach

A solid template melts into an existing liquid workflow.

```bash
bd mol bond mol-hotfix bd-feature-123
```

The proto is instantiated (as liquid by default) and attached to the
existing mol. The new issues become part of the flowing work.

### Bond: Proto + Wisp → Spawn + Attach (Vapor)

A solid template sublimates into an existing vapor workflow.

```bash
bd mol bond mol-extra-check wisp-patrol-456
```

The proto spawns as vapor (following the wisp's phase) and attaches.

### Bond: Mol + Mol → Compound Mol

Two liquid workflows merge into a larger flowing structure.

```bash
bd mol bond bd-feature-123 bd-related-456
```

Links them via dependency edges. Both continue flowing.

### Bond: Wisp + Wisp → Compound Wisp

Two vapor traces merge into a larger ephemeral cloud.

```bash
bd mol bond wisp-123 wisp-456
```

### Phase Override Flags

Bond's spawning behavior can be overridden:

```bash
# Force liquid when attaching to wisp (found something important!)
bd mol bond mol-critical-bug wisp-patrol --pour

# Force vapor when attaching to mol (ephemeral diagnostic)
bd mol bond mol-temp-check bd-feature --wisp
```

| Flag | Effect | Use Case |
|------|--------|----------|
| `--pour` | Force spawn as liquid | "This matters, persist it" |
| `--wisp` | Force spawn as vapor | "This is ephemeral, let it evaporate" |

### Cross-Phase Bonding

What happens when you bond liquid and vapor directly?

```bash
bd mol bond bd-feature-123 wisp-456    # Mol + Wisp
```

**Answer: Reference-only linking.** They connect via dependency edges but
stay in their respective stores. No phase change occurs - you're linking
across the phase boundary without forcing conversion.

This enables patterns like:
- Patrol wisp discovers issue → creates liquid mol for the fix
- Feature mol needs diagnostic → spawns vapor wisp for the check
- The reference survives even when the wisp evaporates (ID stable)

## Agent Attachment: Hooks and Pins

Agents need work attached to them. In Gas Town, this uses **hooks** and **pins**.

### The Hook

Each agent has a **hook** - an anchor point where work hangs. It's the
agent's "pinned bead" - the top of their inbox, the work they're focused on.

```bash
bd hook                                # Show what's on my hook
bd hook --agent deacon                 # Show deacon's hook
```

**Hook states:**
- **Empty (naked)**: Agent awaiting work assignment
- **Occupied**: Agent has work to execute
- **Multiple**: Agent managing several concurrent mols (rare)

### Pin: Attaching Work to Agents

**Pin** attaches a mol to an agent's hook - the action of assigning work.

```bash
bd pin bd-feature-123                  # Pin to my hook
bd pin bd-feature-123 --for witness    # Pin to specific agent's hook
bd unpin                               # Detach current work
```

**The Witness → Polecat flow:**

```bash
# Witness assigns work to polecat
bd pour mol-feature                    # Create liquid mol
bd pin bd-abc123 --for polecat-ace     # Hang on polecat's hook
gt nudge polecat-ace                   # Wake the polecat
```

### Wisps Don't Need Pinning

Wisps are single-cycle and don't survive session boundaries in the
traditional sense. Agents hold them in working memory for one cycle:

```bash
# Deacon self-spawns patrol (no pin needed)
bd wisp mol-deacon-patrol              # Create vapor
# ... execute steps ...
bd mol squash <id> --summary="..."     # Condense and dissipate
# Loop
```

## The Epic-Mol Relationship

Epics and mols are **isomorphic** but represent different mental models.

### Epic: The Business View

An epic is a **simple mol shape** - essentially a TODO list:

```
epic-root
├── child.1
├── child.2
├── child.3
└── child.4
(flat list, no sibling dependencies, execution order implicit)
```

**Properties:**
- One level of children via `.N` numbering
- No explicit serial/parallel encoding
- Human-readable, business-oriented
- The natural shape most humans create

### Mol: The Chemistry View

A mol is the **general case** - arbitrary graphs with explicit workflow
semantics:

```
mol-root
├── phase-A (epic)
│   ├── task.1 ───blocks──→ task.2 (serial)
│   └── task.3 (parallel with task.1)
├── phase-B (epic) ←───blocked-by─── phase-A
│   └── ...
└── standalone (fanout)
(arbitrary DAG, explicit dependencies encode serial/parallel)
```

**Properties:**
- Can contain multiple epics as subgraphs
- Dependency edges encode execution order
- `blocks` = serial (bottleneck)
- No dep = parallel (fanout)
- `conditional` = if-fail path

### The Relationship

**All epics are mols. Not all mols are epics.**

| Aspect | Epic | Mol |
|--------|------|-----|
| Shape | Flat (root + children) | Arbitrary DAG |
| Dependencies | Implicit in ordering | Explicit edges |
| Parallelism | Assumed parallel | Encoded in structure |
| Mental model | TODO list | Workflow graph |
| Common use | Simple feature work | Complex orchestration |

When you `distill` an epic, you get a simple proto.
When you `distill` a complex mol, you get a complex proto (preserving structure).

## Thermodynamic Properties

Gas Town's chemistry has thermodynamic properties - work is energy flowing
through the system.

### Work as Energy

```
Proto (potential energy) → Pour/Wisp → Mol/Wisp (kinetic energy) → Squash → Digest (stored work)
```

- **Protos** store potential energy - the capability to do work
- **Mols/Wisps** are kinetic - work actively flowing
- **Digests** are stored energy - crystallized outcomes

### The Audit Trail as Entropy

Every execution increases entropy - creating more history, more records,
more state. Gas Town manages this through:

- **Wisps**: High entropy, but evaporates (entropy contained)
- **Squash**: Compresses entropy into minimal digest
- **Distill**: Reduces entropy by extracting reusable pattern

### The CV Chain

Every agent has a **chain** of work - their CV:

```
Agent CV = ∑(digests) = crystallized capability proof
```

Work completed → digest created → agent's chain grows → capability demonstrated.

This is the foundation for capability-based work matching: your work
history IS your resume. The ledger speaks.

## The Complete Lifecycle

### Feature Work (Liquid Path)

```bash
# 1. Create work from template
bd pour mol-feature --var name=auth

# 2. Pin to agent
bd pin bd-abc123 --for polecat-ace

# 3. Agent executes steps
bd update bd-abc123.design --status=in_progress
# ... cognition ...
bd close bd-abc123.design
bd update bd-abc123.implement --status=in_progress
# ... and so on

# 4. Squash when complete
bd mol squash bd-abc123 --summary="Implemented auth feature"

# 5. Digest remains in permanent ledger
```

### Patrol Work (Vapor Path)

```bash
# 1. Self-spawn wisp (no pin needed)
bd wisp mol-deacon-patrol

# 2. Execute cycle steps
bd close <step-1>
bd close <step-2>
# ...

# 3. Generate summary and squash
bd mol squash <wisp-id> --summary="Patrol complete, no issues"

# 4. Loop
bd wisp mol-deacon-patrol
# ...
```

### Template Creation (Distillation)

```bash
# 1. Complete some work organically
# ... team develops release workflow over several iterations ...

# 2. Distill the pattern
bd mol distill bd-release-v3 --as "Release Workflow"

# 3. Result: new proto available
bd pour mol-release-workflow --var version=2.0
```

## Polymers: Large-Scale Composition

Protos can compose into arbitrarily large **polymers** - chains of molecules
that encode complex multi-phase work.

```bash
# Create polymer from multiple protos
bd mol bond mol-design mol-implement --as "Design and Implement"
bd mol bond mol-design-implement mol-test --as "Full Dev Cycle"
bd mol bond mol-full-dev mol-deploy --as "End to End"
```

**Polymer properties:**
- Preserve phase relationships from constituent protos
- Can encode hours, days, or weeks of work
- Enable "weekend warrior" autonomous operation
- Beads tracks progress; agents execute; humans sleep

### The Cognition Sausage Machine

A large polymer is a **cognition sausage machine**:

```
Proto Polymer (input)
     │
     ▼
┌─────────────────────────────────────────┐
│           GAS TOWN ENGINE               │
│  ┌─────┐  ┌─────┐  ┌─────┐  ┌─────┐    │
│  │Pole │  │Pole │  │Pole │  │Pole │    │
│  │cat  │  │cat  │  │cat  │  │cat  │    │
│  └──┬──┘  └──┬──┘  └──┬──┘  └──┬──┘    │
│     │        │        │        │        │
│     └────────┴────────┴────────┘        │
│              ↓                          │
│         Merge Queue                     │
│              ↓                          │
│          Refinery                       │
└─────────────────────────────────────────┘
     │
     ▼
Completed Work + Digests (output)
```

Feed in a polymer. Get back completed features, merged PRs, and audit trail.

## The Proto Library

Gas Town maintains a **library of curated protos** - the fuel stockpile:

```
~/gt/molecules/
├── mol-engineer-in-box/       # Full quality workflow
├── mol-quick-fix/             # Fast path for small changes
├── mol-code-review/           # Pluggable review dimensions
├── mol-release/               # Release workflow
├── mol-deacon-patrol/         # Deacon monitoring cycle
├── mol-witness-patrol/        # Witness worker monitoring
└── mol-polecat-work/          # Standard polecat lifecycle
```

**Library operations:**

```bash
bd mol list                    # List available protos
bd mol show mol-code-review    # Show proto details
bd pour mol-code-review        # Instantiate for use
```

### Curated vs Organic

- **Curated protos**: Refined templates in the library, battle-tested
- **Organic protos**: Distilled from real work, may need refinement
- **Path**: Organic → refine → curate → library

## Digest ID Stability

When a wisp is created, its head bead ID is **reserved** in permanent storage.
This ensures cross-phase references remain valid:

```
1. bd wisp mol-patrol          → Creates wisp-123 (ID reserved)
2. bd mol bond ... wisp-123    → Reference created
3. bd mol squash wisp-123      → Digest takes same ID
4. Reference still valid       → Points to digest now
```

**Implementation:**
- On wisp creation: Write placeholder/tombstone to permanent beads
- On squash: Replace placeholder with actual digest
- Cross-phase references never break

## Summary of Operators

| Operator | From | To | Effect |
|----------|------|------|--------|
| `pour` | Proto | Mol | Instantiate as persistent liquid |
| `wisp` | Proto | Wisp | Instantiate as ephemeral vapor |
| `bond` | Any + Any | Compound | Combine (polymorphic) |
| `squash` | Mol/Wisp | Digest | Condense to permanent record |
| `burn` | Wisp | Nothing | Discard without record |
| `distill` | Mol/Epic | Proto | Extract reusable template |
| `pin` | Mol | Agent | Attach work to agent's hook |

## Design Implications

### For Beads

1. **New commands**: `bd pour`, `bd wisp`, `bd pin`
2. **Flag changes**: `--persistent` → `--pour` (or phase follows operand)
3. **Wisp storage**: `.beads-wisp/` directory, gitignored
4. **Digest ID reservation**: Placeholder in permanent store on wisp creation

### For Gas Town

1. **Daemon**: Don't attach permanent molecules for patrol roles
2. **Deacon template**: Use `bd wisp mol-deacon-patrol` pattern
3. **Polecat lifecycle**: Consider wisp-based with digest on completion
4. **Hook inspection**: `bd hook` command for debugging

### For Agents

1. **Polecats**: Receive pinned mols, execute, squash, request shutdown
2. **Patrol roles**: Self-spawn wisps, execute cycle, squash, loop
3. **Recovery**: Re-read beads state, continue from last completed step

---

## Appendix: The Vocabulary

| Term | Meaning |
|------|---------|
| **Proto** | Frozen template molecule (solid phase) |
| **Mol** | Flowing work instance (liquid phase) |
| **Wisp** | Ephemeral execution trace (vapor phase) |
| **Digest** | Condensed permanent record |
| **Pour** | Instantiate proto as mol (solid → liquid) |
| **Wisp** (verb) | Instantiate proto as wisp (solid → vapor) |
| **Bond** | Combine molecules (polymorphic) |
| **Squash** | Condense to digest |
| **Burn** | Discard wisp without digest |
| **Distill** | Extract proto from experience |
| **Hook** | Agent's attachment point for work |
| **Pin** | Attach mol to agent's hook |
| **Polymer** | Large composed proto chain |
| **Epic** | Simple mol shape (flat TODO list) |

---

*The chemistry is the interface. The ledger is the truth. The work gets done.*
