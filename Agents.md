# Project Rules & Tooling Guidelines

This document outlines the architectural decisions, tools, and conventions used in the **Forge** project.

## 🛠 Technology Stack

The project is built using **Go** and the **Charm** ecosystem for Terminal User Interfaces (TUI).

| Component | Library | Purpose |
|-----------|---------|---------|
| **Core Framework** | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | The Elm Architecture (Model-View-Update) for TUI apps. |
| **Styling** | [Lip Gloss](https://github.com/charmbracelet/lipgloss) | CSS-like styling for terminal layouts. |
| **Components** | [Bubbles](https://github.com/charmbracelet/bubbles) | Reusable TUI components (Help, Inputs, Spinners, etc.). |
| **Color Handling** | [Go-Colorful](https://github.com/lucasb-eyer/go-colorful) | Advanced color manipulation (gradients, blending). |

### 🧩 Component Strategy

*   **Charm First**: For every feature, check if it can be solved using existing Charm libraries (Bubbles, Lip Gloss, etc.).
*   **Composition over Creation**: Do not implement custom components from scratch if they can be built by composing existing Charm components.

---

## 🏗 Architecture: The Elm Model (MVU)

We follow the **Model-View-Update** pattern strictly:

1.  **Model**: A struct (`type model struct`) holding the application state (window size, data, sub-models).
2.  **Init**: A function returning the initial command (`tea.Cmd`) to run (e.g., I/O, timers).
3.  **Update**: A function `(msg tea.Msg) (tea.Model, tea.Cmd)` that handles messages and returns a new model and command.
4.  **View**: A function returning a `string` that represents the UI based on the current model.

### Rules for MVU

*   **Immutability**: Treat the model as immutable where possible. Return a modified copy.
*   **Side Effects**: Never perform side effects (I/O, API calls) directly in `Update` or `View`. Use `tea.Cmd` for all side effects.
*   **Composition**: Break down complex UIs into smaller sub-models (e.g., `Header`, `Footer`) with their own `Update` and `View` methods if they have internal state.

---

## 🎨 Styling & Layout (Lip Gloss)

*   **Separation of Concerns**: Define styles (colors, padding, borders) in separate variables or functions, not inline within the `View` logic if they are complex.
*   **Responsiveness**: Always use the `width` and `height` from `tea.WindowSizeMsg` to calculate layouts dynamically.
*   **Gradients**: Use `go-colorful` for complex color interpolations (like the header gradient) and convert them to Lip Gloss colors.

### Component Guidelines

*   **Header (`header.go`)**: Handles the branding and top-level status. Uses custom ASCII art and gradients.
*   **Footer (`Footer.go`)**: Displays key bindings and help info using `bubbles/help`.
*   **Content Area**: The central space should dynamically fill the remaining height between Header and Footer.

---

## ⌨️ Input & Commands

*   **Key Bindings**: Use `bubbles/key` to define keymaps. This allows for auto-generating the Help menu in the footer.
*   **Commands (`commands.go`)**: Encapsulate external system calls (like Git operations) in `tea.Cmd` functions.
    *   Return a `Msg` type on completion (e.g., `GitStatusMsg`) to update the model.

## 📝 File Structure

*   `main.go`: Entry point, main model definition, and global update loop.
*   `header.go`: Header component logic.
*   `Footer.go`: Footer and help component logic.
*   `commands.go`: Definitions for `tea.Cmd` actions and their result messages.
*   `Agents.md`: This documentation file.

---

## 🚀 Development Workflow

1.  **Run**: `go run .`
2.  **Build**: `go build .`
3.  **Format**: `go fmt ./...` (Always run before committing)
