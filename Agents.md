# Project Rules

## Styling
- **Framework**: TailwindCSS (v4)
- **Constraint**: ALL styling must be done using TailwindCSS utility classes. Do not use <style> blocks or external CSS files unless absolutely necessary for complex animations or legacy overrides.
- **Design System**:
  - Fonts: Inter (Google Fonts)
  - Colors: Use arbitrary values or theme extensions for premium gradients (e.g. `bg-[radial-gradient(...)]`).

## Verification
- **Automated First**: Always attempt to verify changes automatically using the available tools (browser integration, tests) before requesting manual verification from the user. Only ask for manual verification if automated methods are insufficient or impossible.
