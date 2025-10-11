# Usability Heuristics & Assumptions (Forge TUI)

Date: 2025-10-11

## Purpose
Capture the lightweight heuristic evaluation guiding polish decisions without full formal UX study.

## Heuristics Applied
1. Visibility of System Status
   - Progress bar updates immediately after workflow completion (latency test T078).
   - Monitors panel LEDs change on each poll event; event logs trace transitions.
2. Match Between System & Real World
   - Whitelist commands shown plainly; error states marked with `exit=-1` and red color override.
3. User Control & Freedom
   - Help overlay toggle (`?`) is reversible and fast (<150ms).
   - Focus traversal (Tab) cycles predictably; Escape not yet implemented (future enhancement).
4. Consistency & Standards
   - Hotkey + BubbleZone parity enforced (T077, parity audit at Init).
   - Advisory line prefixed with `ADVISORY:` across buffer.
5. Error Prevention
   - Non-whitelisted commands rejected early (`ValidateCommand` returns structured error).
   - Zero-step workflow guard (T084) prevents undefined execution.
6. Recognition Over Recall
   - Help overlay lists actions dynamically from registry; no need to memorize keys.
7. Flexibility & Efficiency of Use
   - Env flag `FORGE_WORKFLOW_STEPS` enables multi-step workflows for demo while keeping single-step fast by default.
8. Minimalist Design
   - Cached header/footer reduce render overhead; output viewport caching tail slice avoids repeated joins under load.
9. Help & Documentation
   - `quickstart.md` includes workflow and help overlay usage.

## Assumptions
- Terminal width stable; resize-specific styling deferred.
- Single command repeated in workflow acceptable for MVP demonstration.
- Output color collisions are acceptable with documented palette reuse when exhausted.
- Performance thresholds (scroll ≤200ms, help ≤150ms, combined load <200ms) are sufficient for perceived responsiveness.

## Deferred Items
- Mouse wheel scroll parity.
- Escape/Cancel hotkey for in-flight workflow.
- Rich monitor detail panel activation.
- Comprehensive internationalization.

## Acceptance Summary
Polish tasks finalized with tests ensuring latency and interaction stability; no blocking usability risks identified for MVP scope.
