# Chemistry Design Changes

> Implementation roadmap for the molecular chemistry UX described in
> `molecular-chemistry.md`

## Summary of Changes

The chemistry metaphor requires the following changes to Beads and Gas Town:

### Beads Changes

| Change | Priority | Issue |
|--------|----------|-------|
| Add `bd pour` command (alias for `bd mol spawn --pour`) | P0 | Create |
| Add `bd wisp` command (alias for `bd mol spawn`) | P0 | Create |
| Add `bd pin` command for agent attachment | P1 | Create |
| Add `bd hook` command for hook inspection | P2 | Create |
| Rename `--persistent` to `--pour` in `bd mol spawn` | P0 | Update |
| Add `--pour` flag to `bd mol bond` | P1 | Update |
| Implement digest ID reservation for wisps | P1 | Create |

### Gas Town Changes

| Change | Priority | Issue |
|--------|----------|-------|
| Update daemon: remove permanent attachment for patrol | P0 | gt-3x0z.9 |
| Update deacon.md.tmpl: use wisp-based patrol | P0 | gt-3x0z.9 |
| Update witness.md.tmpl: use wisp-based patrol | P1 | Create |
| Add `gt hook` command (thin wrapper around `bd hook`) | P2 | Create |

---

## Detailed Specifications

### 1. `bd pour` Command

**Purpose:** Instantiate a proto as a persistent mol (liquid phase).

**Syntax:**
```bash
bd pour <proto-id> [flags]

Flags:
  --var strings    Variable substitution (key=value)
  --assignee       Assign the root issue to this agent
  --dry-run        Preview what would be created
```

**Implementation:**
- Alias/wrapper for `bd mol spawn <proto-id> --pour`
- Default behavior: creates mol in permanent `.beads/` storage
- Returns the head bead ID of the created mol

**Example:**
```bash
bd pour mol-feature --var name=auth
# Output: Created mol bd-abc123 from mol-feature
```

---

### 2. `bd wisp` Command

**Purpose:** Instantiate a proto as an ephemeral wisp (vapor phase).

**Syntax:**
```bash
bd wisp <proto-id> [flags]

Flags:
  --var strings    Variable substitution (key=value)
  --dry-run        Preview what would be created
```

**Implementation:**
- Alias/wrapper for `bd mol spawn <proto-id>` (wisp is default)
- Creates wisp in `.beads-wisp/` storage
- Reserves digest ID in permanent storage (placeholder)

**Example:**
```bash
bd wisp mol-patrol
# Output: Created wisp bd-xyz789 from mol-patrol
```

---

### 3. `bd pin` Command

**Purpose:** Attach a mol to an agent's hook (work assignment).

**Syntax:**
```bash
bd pin <mol-id> [flags]

Flags:
  --for string     Agent to pin work for (default: current agent)
```

**Implementation:**
1. Look up the mol by ID
2. Set `pinned: true` on the mol's head bead
3. Set `assignee` to the target agent
4. Update `status` to `in_progress` if not already

**Example:**
```bash
# Pin to myself
bd pin bd-abc123

# Pin to specific agent (Witness assigning work)
bd pin bd-abc123 --for polecat-ace
```

**Unpin:**
```bash
bd unpin [mol-id]
# Clears pinned flag, optionally releases assignee
```

---

### 4. `bd hook` Command

**Purpose:** Inspect what's on an agent's hook.

**Syntax:**
```bash
bd hook [flags]

Flags:
  --agent string   Agent to inspect (default: current agent)
  --json           Output in JSON format
```

**Implementation:**
- Query beads for issues where `pinned: true` AND `assignee: <agent>`
- Display the mol(s) attached to the hook

**Example:**
```bash
bd hook
# Output:
# Hook: polecat-ace
# Pinned: bd-abc123 (mol-feature) - in_progress
#   Step: implement (2 of 5)

bd hook --agent deacon
# Output:
# Hook: deacon
# (empty - patrol uses wisps, no persistent attachment)
```

---

### 5. Rename `--persistent` to `--pour`

**Current:**
```bash
bd mol spawn mol-feature --persistent
```

**New:**
```bash
bd mol spawn mol-feature --pour
# or simply:
bd pour mol-feature
```

**Migration:**
- Keep `--persistent` as deprecated alias
- Log warning when `--persistent` is used
- Remove in next major version

---

### 6. Add `--pour` flag to `bd mol bond`

**Purpose:** Override phase when spawning protos during bond.

**Current behavior:**
- Phase follows target (mol → liquid, wisp → vapor)
- `--wisp` forces vapor

**New:**
- Add `--pour` to force liquid even when target is vapor

```bash
# Found important bug during patrol, make it a real issue
bd mol bond mol-critical-bug wisp-patrol-123 --pour
```

---

### 7. Digest ID Reservation

**Problem:** When a wisp is created and later squashed, the digest should
have the same ID so cross-phase references remain valid.

**Solution:** Reserve the ID on wisp creation.

**Implementation:**

1. **On wisp creation (`bd wisp`):**
   - Generate the head bead ID
   - Write a placeholder to permanent beads:
     ```json
     {
       "id": "bd-xyz789",
       "title": "[Wisp Placeholder]",
       "status": "open",
       "labels": ["wisp-placeholder"],
       "description": "Reserved for wisp digest"
     }
     ```
   - Create actual wisp in `.beads-wisp/` with same ID

2. **On squash (`bd mol squash`):**
   - Replace placeholder with actual digest content
   - Delete wisp from `.beads-wisp/`

3. **On burn (`bd mol burn`):**
   - Delete placeholder from permanent beads
   - Delete wisp from `.beads-wisp/`

**Edge cases:**
- Crash before squash: Placeholder remains (orphan cleanup needed)
- Multiple wisps: Each has unique ID, no collision

---

### 8. Daemon Patrol Changes (Gas Town)

**Current behavior (`checkDeaconAttachment`):**
- Checks if Deacon has pinned mol
- If not, spawns `mol-deacon-patrol` and attaches permanently
- This is wrong for wisp-based patrol

**New behavior:**
- Remove `checkDeaconAttachment` entirely
- Deacon manages its own wisp lifecycle
- Daemon just ensures Deacon session is running and pokes it

**Code change in `daemon.go`:**
```go
// Remove this function entirely:
// func (d *Daemon) checkDeaconAttachment() error { ... }

// Or replace with a simpler check:
func (d *Daemon) ensureDeaconReady() error {
    // Just verify session is running, don't attach anything
    // Deacon self-spawns wisps for patrol
    return nil
}
```

---

### 9. Deacon Template Update

**Current (`deacon.md.tmpl`):**
```markdown
If no molecule (naked), **start a new patrol**:
```bash
bd mol run mol-deacon-patrol
```
```

**New:**
```markdown
## Patrol Cycle (Wisp-Based)

Each patrol cycle uses ephemeral wisps:

```bash
# 1. Spawn wisp for this cycle
bd wisp mol-deacon-patrol

# 2. Execute steps
bd close <step-1>
bd close <step-2>
# ...

# 3. Squash with summary
bd mol squash <wisp-id> --summary="Patrol complete: <findings>"

# 4. Loop
# Repeat from step 1
```

**Why wisps?**
- Patrol cycles are operational, not auditable work
- Each cycle is independent
- Only the digest matters (and only if notable)
- Keeps permanent beads clean
```

---

## Implementation Order

### Phase 1: Core Commands (P0)

1. [ ] Add `bd pour` command
2. [ ] Add `bd wisp` command
3. [ ] Rename `--persistent` to `--pour` (with deprecated alias)
4. [ ] Update daemon to remove `checkDeaconAttachment`
5. [ ] Update `deacon.md.tmpl` for wisp-based patrol

### Phase 2: Agent Attachment (P1)

1. [ ] Add `bd pin` command
2. [ ] Add `bd unpin` command
3. [ ] Add `--pour` flag to `bd mol bond`
4. [ ] Implement digest ID reservation for wisps
5. [ ] Update `witness.md.tmpl` for wisp-based patrol

### Phase 3: Inspection (P2)

1. [ ] Add `bd hook` command
2. [ ] Add `gt hook` command (thin wrapper)

---

## Testing Plan

### Manual Tests

```bash
# Test pour
bd pour mol-quick-fix
bd show <id>  # Verify in permanent beads

# Test wisp
bd wisp mol-patrol
ls .beads-wisp/  # Verify wisp created
bd show <id>  # Should work from permanent (placeholder)

# Test squash
bd mol squash <wisp-id> --summary="Test"
ls .beads-wisp/  # Wisp should be gone
bd show <id>  # Digest should exist

# Test pin
bd pour mol-feature
bd pin <id>
bd hook  # Should show pinned mol
```

### Integration Tests

- Deacon patrol cycle with wisps
- Cross-phase bonding (mol + wisp)
- Digest ID stability after squash

---

## Migration Notes

### Existing Code

- `bd mol spawn` defaults to wisp (vapor) now
- Code using `bd mol spawn` for permanent mols needs `--pour`
- `bd mol run` continues to work (creates mol, not wisp)

### Deprecation Path

| Old | New | Deprecation |
|-----|-----|-------------|
| `--persistent` | `--pour` | Warn in 0.x, remove in 1.0 |
| `bd mol spawn` (for mols) | `bd pour` | Keep both, prefer new |

---

*This document tracks the implementation of chemistry UX changes.*
