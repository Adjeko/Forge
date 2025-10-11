<!--
Sync Impact Report
Version change: 1.2.0 → 1.3.0 (MINOR)
Modified principles:
	Added new Principle VI: Accessibility & Dual Input (Hotkeys + Mouse)
Added sections: None (principle appended)
Removed sections: None
Templates requiring updates:
	✅ .specify/templates/plan-template.md (added accessibility gate)
	✅ .specify/templates/spec-template.md (added accessibility mapping guidance)
	✅ .specify/templates/tasks-template.md (tasks must include hotkey + BubbleZone binding steps)
	⚠ README.md (document universal hotkey conventions & interaction parity)
Follow-up TODOs:
	TODO(README): Add table of global hotkeys & interaction zones
-->

# Forge Constitution

## Core Principles

### I. Terminal-First UI
Forge MUST deliver all user-facing functionality via a TUI, prioritizing clarity, accessibility, and responsiveness in the terminal. All components MUST be interactive and visually consistent.
Rationale: Ensures the app remains focused on its core value as a TUI and leverages the strengths of terminal-based workflows.

### II. Charm Ecosystem Compliance & Composition-First
All visual and interactive components MUST use Charm libraries (Bubbletea, Bubble, Lipgloss, Huh, BubbleZone) for rendering, styling, and interactivity. Before implementing any "new" UI component, you MUST attempt to compose it from existing Charm primitives and reusable Bubbles (composition-first rule). Only if the desired behavior cannot be achieved through composition MAY a bespoke implementation be created.

Mandatory order of approach:
1. Reuse an existing Forge component as-is (prefer).
2. Compose from existing Charm primitives (e.g., Bubbletea model + Modal + Label + Button + Lipgloss styles).
3. Extend/wrap an existing Bubble non-invasively.
4. Create a new custom component (last resort) WITH written justification in the PR describing: gap, attempted compositions, why insufficient. If broadly useful, open/track an upstream issue in the Charm ecosystem.

Example: No ready-made confirmation dialog? Compose one from Modal + Label + Button(s) instead of writing a new dialog renderer. (Beispiel in Deutsch: Wenn es keinen Bestätigungs-Dialog gibt, baue ihn aus Modal, Label und Button zusammen, bevor du eine eigene Implementierung schreibst.)

Prohibited without justification: Direct low-level terminal drawing logic bypassing Charm abstractions.

Rationale: Guarantees a unified look/feel, minimizes maintenance surface, accelerates development, and encourages upstream contribution rather than local divergence.

### III. Test-Driven Development (NON-NEGOTIABLE)
All new features and bug fixes MUST be developed using TDD. Tests MUST be written before implementation, fail initially, and pass only after code is complete. Red-Green-Refactor cycle is strictly enforced.
Rationale: Ensures reliability, maintainability, and user trust in every release.

### IV. Integration & Contract Testing
Integration tests MUST cover all interactions between Bubbletea models, BubbleZone click events, and external Go packages. Contract changes require new/updated tests.
Rationale: Prevents regressions and ensures seamless user experience across all interactive flows.

### V. Observability & Simplicity
Forge MUST provide structured logging for all major actions and errors. Simplicity is enforced: avoid unnecessary features, keep UI minimal, and follow YAGNI (You Aren't Gonna Need It) principles.
Rationale: Improves debuggability and keeps the app maintainable and user-friendly.

### VI. Accessibility & Dual Input (Hotkeys + Mouse)
Every user-facing action MUST be operable via keyboard hotkeys AND via mouse (BubbleZone clickable areas) with functional parity. No feature may exist that is only accessible by mouse or only by keyboard. All interactive regions MUST register BubbleZone hit areas; all actions MUST have mnemonic or consistent hotkey bindings (e.g., first letter, or global scheme). Hotkeys MUST be discoverable: either listed in a contextual help panel, footer hint, or inline legend. Navigation between focusable elements MUST be possible with standard keys (Tab/Shift+Tab or arrow keys) and MUST update visual focus state.

Normative rules:
1. Parity: For each command/action, doc MUST show: Action | Hotkey(s) | BubbleZone ID.
2. Discoverability: Pressing a universal help key (e.g., `?`) MUST toggle an overlay showing available hotkeys for current context within 150ms.
3. Consistency: Global hotkeys (`q` quit, `?` help, `Esc` cancel/close, `Enter` confirm) MUST NOT be repurposed locally without justification.
4. Focus Management: Focus changes MUST be visually indicated (style or cursor) and never lost on window resize.
5. No Hidden Zones: BubbleZone areas MUST match visible interactive UI; invisible click targets are prohibited.

Rationale: Ensures Forge is usable efficiently by power users (keyboard) and discoverable for exploratory users (mouse), improving accessibility, lowering learning curve, and enabling mixed-modality workflows.

## Technology Stack & Constraints
Forge is written in Go. All TUI components MUST use Charm libraries (Bubbletea, Bubble, Lipgloss, Huh, BubbleZone). No other UI libraries are permitted unless explicitly justified and documented. All dependencies MUST be open source and actively maintained.

## Development Workflow & Quality Gates
All code changes MUST pass automated tests and code review. Every PR MUST demonstrate compliance with TDD and integration testing principles. Releases are versioned using semantic versioning (MAJOR.MINOR.PATCH). Breaking changes require a migration plan and explicit user communication.

## Governance

This constitution supersedes all other development practices. Amendments require documentation and a migration plan for affected features. As the sole developer, you are responsible for verifying compliance with all principles before making changes. Complexity must be justified in writing. Runtime guidance is maintained in the README.md and updated with every constitution change.

<!--
Sync Impact Report
Version change: 1.0.0 → 1.1.0
Modified principles: All principles rewritten for Forge context
Added sections: Technology Stack & Constraints, Development Workflow & Quality Gates
Removed sections: None
Templates requiring updates:
✅ plan-template.md (Constitution Check gates updated)
✅ spec-template.md (requirements and testing alignment)
✅ tasks-template.md (task types and test-first discipline)
⚠ README.md (ensure runtime guidance matches principles)
Follow-up TODOs:
TODO(RATIFICATION_DATE): Set original ratification date if known
-->
**Version**: 1.3.0 | **Ratified**: 2025-10-10 | **Last Amended**: 2025-10-10