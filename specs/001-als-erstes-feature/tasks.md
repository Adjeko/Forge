# Tasks: Forge TUI Task Runner Core Interface & Workflow Execution

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/internal-exec.md
**Prerequisites**: plan.md (complete), spec.md (with FR-022..024), research decisions finalized

## Format
`[ID] [P?] [Story] Description`

---
## Phase 1: Setup (Shared Infrastructure)

Purpose: Initialize repository structure, dependencies, and base tooling.

- [ ] T001 Create base Go module (`go mod init forge`) at repo root
- [ ] T002 Add dependency imports (bubbletea, lipgloss, bubblezone, huh) in a placeholder `main.go` (build succeeds)
- [ ] T003 [P] Create directory scaffold (`src/ui/...`, `src/exec`, `src/output`, `src/logging`, `tests/...`)
- [ ] T004 [P] Add `Makefile` (or simple build script) with targets: build, test, lint (placeholder lint step)
- [ ] T005 Configure basic `.gitignore` (Go build artifacts, temp files)

Checkpoint: Repository builds (`go build ./...`).

---
## Phase 2: Foundational (Blocking Prerequisites)

Purpose: Core abstractions all stories depend on.

- [ ] T006 Define `PrimitiveCommand` and whitelist metadata in `src/exec/primitives.go`
- [ ] T007 [P] Implement color palette + deterministic assignment in `src/output/color.go`
- [ ] T008 [P] Implement `OutputBuffer` struct (append, viewport slice) in `src/output/buffer.go`
- [ ] T009 Implement logging setup (`slog` wrapper) in `src/logging/logger.go`
- [ ] T010 [P] Define hotkey action registry types in `src/ui/accessibility/actions.go`
- [ ] T011 [P] Define focus manager skeleton in `src/ui/accessibility/focus.go`
- [ ] T012 Create root Bubbletea model skeleton in `src/ui/model/root.go` (init sub-model placeholders)
- [ ] T013 Implement gradient header & static footer layout skeleton in `src/ui/components/chrome.go`
- [ ] T014 Add BubbleZone registration helper utilities in `src/ui/zones/zones.go`
- [ ] T015 Wire main loop `main.go` minimal run (quit on 'q')
### Composition Attempt Documentation (Moved Earlier per Constitution)
// Each new UI component requires a composition attempt doc BEFORE non-trivial custom code.
- [ ] T086A [P] Create composition attempt template (`docs/composition/attempt-template.md`) with fields: Component, Intended Behavior, Attempted Charm Primitives, Gaps, Decision Rationale.
- [ ] T086B [P] Add composition attempt doc for CommandList component.
- [ ] T086C [P] Add composition attempt doc for OutputViewport component.
- [ ] T086D [P] Add composition attempt doc for ProgressBar component.
- [ ] T086E [P] Add composition attempt doc for MonitorsPanel component.
- [ ] T086F [P] Add composition attempt doc for HelpOverlay component.

Checkpoint: App launches with header & footer, quits via 'q'; no command execution yet.

### Added Pre-Implementation Test Harness (Remediation)
- [ ] T068 [P] Foundational Test: Viewport retention (>5000 lines) (FR-010a)
- [ ] T069 [P] Foundational Test: Input latency harness scaffold (FR-010b) (define measurement abstraction)
- [ ] T070 [P] Foundational Test: Layout degrade rendering at width 99 (FR-013)
- [ ] T071 [P] Foundational Test: Hanging detection simulation (virtual clock) (FR-017)
- [ ] T072 [P] Foundational Test: Non-whitelist error message pattern (FR-021)
- [ ] T073 [P] Foundational Test: Advisory threshold warning trigger at 100k lines (Assumption)
- [ ] T074 [P] Foundational Test: Composition-first UI assembly smoke (constitution)
- [ ] T075 [P] Foundational Test: Accessibility parity snapshot (hotkey vs zone) (Principle VI)
- [ ] T086G [P] Meta-test: Verify composition attempt docs exist for all declared components (CommandList, OutputViewport, ProgressBar, MonitorsPanel, HelpOverlay) BEFORE their implementation tasks run.

---
## Phase 3: User Story 1 - Einzelnes primitives Kommando ausführen (P1) 🎯 MVP

Goal: Execute single whitelisted primitive command and display colored output.
Independent Test: Start app → run `git status` → colored output block with exit status.

### Tests (TDD)
- [ ] T016 [P] [US1] Test: whitelist validation rejects non-listed command (exec layer)
- [ ] T017 [P] [US1] Test: color assignment deterministic for given primitive ID
- [ ] T018 [P] [US1] Test: output buffer append & viewport integrity (lines retained)

### Implementation
- [ ] T019 [P] [US1] Implement single command execution function in `src/exec/workflow.go` (single-step path)
- [ ] T020 [P] [US1] Integrate execution into root model (trigger on Enter) in `src/ui/model/root.go`
- [ ] T021 [P] [US1] Implement command list component in `src/ui/components/command_list.go` with BubbleZones per command
- [ ] T022 [US1] Hook output streaming into buffer + color application (update view) in `src/ui/model/root.go`
- [ ] T023 [US1] Render output viewport component `src/ui/components/output_view.go` (scroll placeholder)
- [ ] T024 [US1] Display exit status + colored section delimiter in footer
- [ ] T025 [US1] Register hotkeys (focus command list, run command, quit) in actions registry

Checkpoint: Single command execution fully functional & test suite green.

---
## Phase 4: User Story 2 - Sequenzieller Workflow (P2)

Goal: Execute multi-step workflow sequentially, fail-fast, show progress bar updates.
Independent Test: 3-step workflow stops on second step failure; progress halts accordingly.

### Tests
- [ ] T026 [P] [US2] Test: workflow stops on first failing step (fail-fast)
- [ ] T027 [P] [US2] Test: progress percentage matches completed steps

### Implementation
- [ ] T028 [P] [US2] Extend workflow orchestration (multi-step) in `src/exec/workflow.go`
- [ ] T029 [P] [US2] Implement progress bar component in `src/ui/components/progress_bar.go`
- [ ] T030 [US2] Update footer rendering to show dynamic progress gradient fill
- [ ] T031 [US2] Add per-step color labeling in output (reuse palette) with separation markers
- [ ] T032 [US2] Record step exit codes & final workflow status in model
- [ ] T033 [US2] Register workflow-related hotkeys if distinct (reuse run key)

Checkpoint: Multi-step workflows run & fail-fast; progress accuracy validated.

---
## Phase 5: User Story 3 - Hintergrund Status-Monitore (P3)

Goal: Background monitors poll targets periodically; LEDs reflect states (Checking/OK/Failed/Disabled).
Independent Test: Reachable & unreachable ping monitors update appropriately within interval.

### Tests
- [ ] T034 [P] [US3] Test: monitor poll updates state within ≤6s (default interval 5s)
- [ ] T035 [P] [US3] Test: jitter prevents synchronized poll (timestamps variance) (allow ±500ms)
- [ ] T076 [P] [US3] Test: Monitor LED state accuracy classification (OK/Failed/Hanging) mapping
- [ ] T080 [P] [US3] Test: Configurable interval override (FR-018) via `--monitor-interval` (set to 2s, assert poll cadence)
- [ ] T081 [P] [US3] Test: Monitor type dispatch & script constraints (FR-019) (mock ping vs script)

### Implementation
- [ ] T036 [P] [US3] Implement monitor type structs & interface in `src/exec/monitors.go`
- [ ] T037 [P] [US3] Implement polling scheduler with jitter and non-blocking design
- [ ] T038 [US3] LED panel component in `src/ui/components/monitors_panel.go`
- [ ] T039 [US3] Integrate monitor updates into root model tick/update cycle
- [ ] T040 [US3] Register BubbleZones for each LED; click currently no-op (future details)

Checkpoint: Monitors visible & update concurrently without input lag regression.

---
## Phase 6: User Story 4 - Deterministische Farbzuteilung (P4)

Goal: Stable, deterministic palette reuse across sessions; errors override with red.
Independent Test: Re-run identical workflow; colors stable except error steps forced red.

### Tests
- [ ] T041 [P] [US4] Test: same primitive sequence yields identical color set across two runs
- [ ] T042 [P] [US4] Test: error step overrides assigned color with red
- [ ] T082 [P] [US4] Test: Palette exhaustion reuse strategy (FR-012) (more steps than palette size)

### Implementation
- [ ] T043 [US4] Refine color hashing function & collision handling in `src/output/color.go`
- [ ] T044 [US4] Add palette exhaustion fallback & reuse strategy docs inline

Checkpoint: Color determinism demonstrated via tests.

---
## Phase 7: Accessibility & Interaction (Cross-Cutting)

Goal: Enforce Principle VI parity (hotkeys + BubbleZone) and add help overlay & scrolling.
Independent Test: Press `?` → overlay appears ≤150 ms listing actions; Tab cycles regions; PgUp/PgDn scroll output; input latency within bounds.

### Tests
- [ ] T045 [P] [ACC] Test: help overlay toggles & renders within ≤150 ms
- [ ] T046 [P] [ACC] Test: focus traversal cycles all regions (CommandList→Output→Monitors→Help→CommandList)
- [ ] T047 [P] [ACC] Test: output scroll maintains ≤200 ms input latency under stream simulation
- [ ] T077 [P] [ACC] Test: Hotkey & BubbleZone action parity audit (no missing modality)
- [ ] T083 [P] [ACC] Test: Footer completion state visible after workflow end (FR-011)

### Implementation
- [ ] T048 [P] [ACC] Implement help overlay component in `src/ui/help/overlay.go`
- [ ] T049 [P] [ACC] Populate overlay from hotkey action registry (auto-generate table)
- [ ] T050 [ACC] Implement focus traversal logic & visual indicators (styles) in `src/ui/accessibility/focus.go`
- [ ] T051 [ACC] Implement scrolling logic (PgUp/PgDn + wheel) in `src/ui/components/output_view.go`
- [ ] T052 [ACC] Add parity verification routine (assert all actions have both modalities) logged at startup

Checkpoint: Accessibility parity demonstrable; overlay & focus fully functional.

---
## Phase 8: Logging & Observability (Cross-Cutting)

Goal: Structured events for execution, workflows, monitors, UI state.
Independent Test: Running workflow emits expected ordered events; monitor poll entries appear; help overlay toggle logged.

### Tests
- [ ] T053 [P] [OBS] Test: command start/end events emitted with correct fields
- [ ] T054 [P] [OBS] Test: monitor poll event sequence (start → result) captured
- [ ] T078 [P] [OBS] Test: progress bar update latency ≤250 ms after step completion (SC-003/004)
- [ ] T079 [P] [OBS] Test: combined load latency (<200 ms input under 50k stream + 5 monitors) (SC-009)
- [ ] T084 [P] [OBS] Test: Zero-step workflow validation (FR-009) (attempt start → error, no crash)
- [ ] T085 [P] [OBS] Test: Fixed chrome (header/footer immobile on scroll) (FR-014)
- [ ] T088 [P] [OBS] Test: Startup readiness time (SC-001) (measure time from process start to first command runnable ≤5s)
- [ ] T089 [P] [OBS] Test: High-volume stream generator (50k lines) maintains ≤200 ms input latency (SC-002, FR-010b)
- [ ] T090 [P] [OBS] Test: Monitor/workflow state separation (FR-015) (failing monitor does not alter workflow step statuses)

### Implementation
- [ ] T055 [P] [OBS] Add event constants & logging wrappers in `src/logging/logger.go`
- [ ] T056 [OBS] Insert logging calls in workflow execution path `src/exec/workflow.go`
- [ ] T057 [OBS] Insert logging in monitor scheduler `src/exec/monitors.go`
- [ ] T058 [OBS] Insert logging in UI interactions (focus changes, overlay toggle) `src/ui/model/root.go`

Checkpoint: Logs provide traceability for major actions.

---
## Phase 9: Polish & Cross-Cutting Concerns

- [ ] T059 [P] Refine layout styling & gradient rendering performance
- [ ] T060 [P] Add advisory warning when output lines exceed threshold (e.g., 100k) in buffer
- [ ] T061 [P] Final interaction audit (parity/help/focus resize stress) consolidating previous parity/help/focus tasks (references T045, T046, T047, T077)
- [ ] T064 Performance profiling (measure input latency under load) & adjustments (post T089) 
// T086 moved earlier as T086A–T086G (Foundational phase) to satisfy Composition-First gate.
- [ ] T087 Usability heuristic / SC-005 evaluation doc (record method or convert to assumption with justification)
- [ ] T065 Security review: whitelist enforcement & no shell injection possibility
- [ ] T066 Documentation updates in `quickstart.md` (add workflow + overlay instructions)
- [ ] T067 Code cleanup & comments pass

Checkpoint: Ready for release of MVP feature set.

---
## Dependencies & Execution Order

Phase Dependencies:
- Setup → Foundational → (User Story phases 3–6 can only start after Foundational)
- Accessibility, Logging phases (7–8) depend on core user stories 1–4; can start after US4 or partially after needed components exist.
- Polish last.

User Story Independence:
- US1 independent once foundational ready.
- US2 builds on US1 execution structures.
- US3 independent of US2 logic but shares foundational + execution integration.
- US4 relies on color logic from foundational + usage patterns (can begin after US1 color usage validated).

Parallel Opportunities:
- Marked [P] tasks within different files.
- Tests for a story can run in parallel after skeletons exist.

---
## Parallel Execution Examples

User Story 1 (after Foundational):
- Parallel: T016, T017, T018
- Then: T019, T020, T021 (parallel), followed by T022, T023, T024, T025

User Story 3 (after Foundational):
- Parallel: T034, T035 tests; T036, T037 impl
- Then: T038, T039, T040 integration

---
## Implementation Strategy

MVP Path:
1. Phase 1 + 2
2. Phase 3 (US1) → deliver MVP (single command execution)
3. Optionally stage release, then proceed with US2–US4 → Accessibility → Logging → Polish

Incremental Delivery:
- Deliver value after each user story checkpoint; ensure independence & tests pass.

---
## Task Counts
Total Tasks: 90
- Setup: 5
- Foundational: 10 (+8 remediation tests = 18 including added harness tasks)
- US1: 10
- US2: 6
- US3: 7 (+3 remediation tests = 10)
- US4: 4 (+1 exhaustion test = 5)
- Accessibility: 8 (+2 remediation tests = 10)
- Observability: 6 (+7 remediation tests = 13)
- Polish: 11 (+1 consolidated audit +2 documentation tasks = 14)
Additional Remediation / Early Tests: 11 earlier (T068–T079) + 9 (T080–T087) + 3 (T088–T090) = 23 total remediation-related tasks

Parallelizable Tasks (marked [P]): 39

---
## Independent Test Criteria Summary
- US1: Single command run with colored output & status
- US2: Fail-fast multi-step workflow with accurate progress
- US3: Monitors reflect state changes within interval
- US4: Stable deterministic colors across runs
- Accessibility: Overlay speed, focus cycle, scroll latency
- Observability: Event emission coverage, progress latency, combined load latency
- Early Foundations: Viewport retention, layout degrade, hanging detection, advisory threshold, error pattern, composition-first assembly, parity snapshot

---
