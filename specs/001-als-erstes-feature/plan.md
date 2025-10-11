# Implementation Plan: Forge TUI Task Runner Core Interface & Workflow Execution

**Branch**: `001-als-erstes-feature` | **Date**: 2025-10-10 | **Spec**: ./spec.md
**Input**: Feature specification from `/specs/001-als-erstes-feature/spec.md`

## Summary

Deliver the initial Forge fullscreen TUI: run whitelisted primitive commands singly or as sequential fail-fast workflows, render colored step output with progress footer, and display independent background status monitors (LEDs). Enforce accessibility parity (hotkeys + BubbleZone). Establish buffers and performance scaffolding for large streaming output.

## Technical Context

**Language/Version**: Go 1.23 (assumed latest stable)  
**Primary Dependencies**: bubbletea, lipgloss, bubblezone, huh  
**Storage**: In-memory only (no persistence MVP)  
**Testing**: `go test` + golden snapshot tests for rendering + timing assertions  
**Target Platform**: Cross-platform terminal (Windows, macOS, Linux)  
**Project Type**: single  
**Performance Goals**: <200ms input latency while streaming 50k lines; monitor state reflect change ≤6s; progress update ≤250ms after step end  
**Constraints**: Min width 100 columns for dual panel; hang detection 60s; unlimited buffer (risk accepted)  
**Scale/Scope**: Single user; ≤10 primitive commands initial whitelist; ≤5 concurrent status monitors  

Unknowns / Research Items (Initial) → Resolution Status:
-- R1 Buffer rendering: Virtualization vs incremental diff → DECIDED: Incremental diff + viewport retention (supports FR-010a)
-- R2 Zones & focus: Integration strategy → DECIDED: Central FocusRegistry + BubbleZone IDs (supports FR-022–024)
-- R3 Logging choice: slog vs zerolog → DECIDED: Go slog (structured + stdlib) with event taxonomy
-- R4 Memory risk: Advisory warning at threshold → DECIDED: Soft warning at 100k lines (Assumption & future window strategy optional)

Remaining Open Research: NONE (All critical pre-implementation unknowns resolved)

## Constitution Check (Pre-Design)

- Terminal-first: PASS (TUI only)
- Charm compliance: PASS (Charm libs exclusively)
- Composition-first: PASS (compose before custom)
- Accessibility & dual input: PARTIAL (help overlay FR missing → add FR-022)
- TDD: PLAN (tests enumerated per story; write before code)
- Observability: PARTIAL (logger choice & event taxonomy pending research)
- Simplicity: PASS (cancellation deferred intentionally)

Actions Required: Add FR-022 (help overlay). Define logging events & library.

## Project Structure

```
src/
├── ui/
│   ├── model/
│   ├── components/
│   ├── help/
│   ├── styles/
│   ├── zones/
│   └── accessibility/
├── exec/
│   ├── primitives.go
│   ├── workflow.go
│   └── monitors.go
├── output/
│   ├── buffer.go
│   └── color.go
├── logging/
│   └── logger.go
├── main.go
└── internal/

tests/
├── integration/
├── contract/
├── ui/
└── exec/
```

**Structure Decision**: Domain-based layering for clarity & composition-first design.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| (none) | | |

## Post-Design Constitution Recheck (Remediation Pass)

Gate Status:
* Terminal-first: PASS
* Charm compliance: PASS
* Composition-first: PASS (interfaces & entities defined pre-implementation)
* Accessibility & dual input: PASS (Cancel removed/deferred; parity table aligned; FR-010 split does not affect accessibility)
* TDD: IMPROVEMENT NEEDED → Added upcoming tasks for: hanging detection test (FR-017), layout degrade (FR-013), monitor accuracy (SC-007/008), progress latency (SC-003/004), combined latency (SC-009), error pattern (FR-021), advisory threshold (Assumption), viewport retention (FR-010a)
* Observability: PASS (event taxonomy & slog chosen; tests to assert log emission patterns will be queued)
* Simplicity: PASS (Cancel deferred; no premature complexity)

Action: Insert new foundational test tasks before any production code to enforce TDD & composition-first documentation.

Spec Deltas Incorporated: FR-010 split (FR-010a/FR-010b), FR-013 threshold concretized (<100 width), FR-017 algorithm precise, FR-021 error pattern defined, advisory threshold assumption, cancel gap documented.

Result: No blocking gate violations remain once new tasks committed.

## Foundational Test Coverage Additions (Pre-Code)

Early tests to be authored PRIOR to implementation:
1. Viewport Retention (simulate >5000 lines and assert last 5000 accessible) – FR-010a
2. Input Latency Harness skeleton (placeholder, will measure later) – FR-010b
3. Layout Degrade Rendering (<100 cols) – FR-013
4. Hanging Detection Simulation (no output for 60s virtual clock) – FR-017
5. Non-Whitelist Command Rejection Pattern – FR-021
6. Monitor Poll Accuracy & Timing (jitter within ≤500ms, reflect ≤1.2× interval) – SC-007/008
7. Progress Bar Update Latency (≤250ms) – SC-003/004
8. Combined Load Latency (<200ms input under stream + monitors) – SC-009
9. Advisory Threshold Warning at 100k lines – Assumption
10. Composition-First Smoke Test (assemble UI from existing bubbletea primitives only) – Constitution
11. Accessibility Parity Early Snapshot (hotkey vs zone activation produce identical state transition) – Principle VI
12. Startup Readiness Time (process start → command runnable ≤5s) – SC-001
13. High-Volume Stream Latency (synthetic 50k lines generator maintains ≤200 ms input latency) – SC-002 / FR-010b
14. Monitor/Workflow State Separation (monitor failure does not mutate workflow step statuses) – FR-015

### TDD Gating & Composition Documentation

Implementation MAY NOT proceed beyond skeletal stubs until alle oben gelisteten Tests (1–11) als initial failing Tests committed wurden (TDD Gate – Constitution Principle III). Jede neue UI-Komponente MUSS eine dokumentierte "Composition Attempt" Notiz besitzen (Task wird ergänzt), welche beschreibt: (a) Welche vorhandenen Charm-Komponenten ausprobiert wurden, (b) Warum sie unzureichend waren, (c) Begründung für Custom Code (falls notwendig). Diese Notizen werden als Markdown-Snippets unter `docs/composition/` abgelegt und durch einen Meta-Test validiert.
