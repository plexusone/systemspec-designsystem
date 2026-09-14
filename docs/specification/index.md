# Specification Overview

DSS defines **9 canonical layers** for complete design system specification. Each layer serves a specific purpose and can be defined in separate JSON/YAML files or combined.

## Layer Summary

| Layer | File(s) | Purpose |
|-------|---------|---------|
| [Meta](meta.md) | `meta.json` | System name, version, maintainers |
| [Principles](principles.md) | `principles.json` | Design philosophy and guidelines |
| [Foundations](foundations.md) | `foundations/*.json` | Design tokens (colors, typography, spacing, density) |
| [Components](components.md) | `components/*.json` | UI elements with variants, states, props |
| [Patterns](patterns.md) | `patterns/*.json` | Multi-component solutions |
| Templates | `templates/*.json` | Page-level layouts |
| Content | `content.json` | Voice & tone guidelines |
| [Accessibility](accessibility.md) | `accessibility.json` | WCAG compliance requirements |
| Governance | `governance.json` | Versioning and deprecation policies |
| [Theming](theming.md) | `themeBindings.json` | Token mappings to external components |

## Directory Structure

A complete design system spec typically follows this structure:

```
my-design-system/
├── meta.json                 # Required: name, version
├── modes.json                # Declared discrete modes (e.g. ["light", "dark"])
├── principles.json           # Design philosophy
├── accessibility.json        # WCAG requirements
├── governance.json           # Policies
├── content.json              # Voice & tone
├── themeBindings.json        # Token mappings to external components
├── foundations/
│   ├── colors.json           # Color tokens (may include per-mode values)
│   ├── typography.json       # Font definitions
│   ├── spacing.json          # Spacing scale
│   ├── densities.json        # Density variants (scale + spacing overrides)
│   └── border-radius.json    # Border radius values
├── components/
│   ├── button.json           # Button component spec (includes themingContract)
│   ├── card.json             # Card component spec
│   └── input.json            # Input component spec
├── patterns/
│   └── form.json             # Form pattern spec
└── templates/
    └── dashboard.json        # Dashboard template
```

## Minimal Example

The minimum required spec is just `meta.json`:

```json
{
  "name": "My Design System",
  "version": "1.0.0"
}
```

## Loading Behavior

The DSS loader:

1. Looks for `meta.json` in the specified directory
2. Loads `principles.json`, `accessibility.json`, etc. if present
3. Scans `foundations/`, `components/`, `patterns/`, `templates/` subdirectories
4. Merges all files into a single `DesignSystem` object
5. Validates the combined spec

## File Formats

DSS supports both JSON and YAML:

```json
// colors.json
{
  "colors": [
    { "id": "primary", "value": "#0066CC" }
  ]
}
```

```yaml
# colors.yaml
colors:
  - id: primary
    value: "#0066CC"
```

## LLM Context

Any component, pattern, or template can include an `llm` field for AI code generation optimization:

```json
{
  "llm": {
    "intent": "What this element is for",
    "allowedContexts": ["where-to-use-it"],
    "forbiddenContexts": ["where-not-to-use-it"],
    "antiPatterns": ["common mistakes to avoid"],
    "exampleUsage": ["<Component>Example</Component>"]
  }
}
```

See [Components](components.md) for full LLM context documentation.
