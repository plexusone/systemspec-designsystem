# Theming

DSS provides a formal mechanism for component libraries to publish themeable tokens and for applications to bind their design tokens to those components.

## Overview

The theming system consists of two complementary schemas:

| Schema | Used By | Purpose |
|--------|---------|---------|
| `themingContract` | Component libraries | Declares which CSS properties can be themed externally |
| `themeBindings` | Applications | Maps application tokens to component theming APIs |

This separation enables:

- Components remain agnostic of consuming applications
- Applications can theme any component that publishes a contract
- Tools can validate bindings and generate CSS automatically

## Theming Contract (Component Side)

Components declare their external theming API via `themingContract`:

```json
{
  "themingContract": {
    "prefix": "--mde-theme",
    "description": "CSS custom properties for external theming",
    "tokens": [
      {
        "id": "bg-primary",
        "cssProperty": "--mde-theme-bg-primary",
        "semantic": "background",
        "description": "Main background color",
        "defaultLight": "#ffffff",
        "defaultDark": "#0f172a"
      }
    ]
  }
}
```

### Contract Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `prefix` | string | Yes | CSS custom property prefix (e.g., `--mde-theme`) |
| `description` | string | No | Human-readable description of the theming API |
| `tokens` | array | Yes | List of themeable tokens |

### Token Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Token identifier (e.g., `bg-primary`) |
| `cssProperty` | string | Yes | Full CSS custom property name |
| `semantic` | string | Yes | Semantic category for automatic mapping |
| `description` | string | No | What this token controls |
| `defaultLight` | string | Yes | Default value for light mode |
| `defaultDark` | string | Yes | Default value for dark mode |

### Semantic Categories

The `semantic` field enables automatic token matching between component and application:

| Semantic | Meaning | Typical Tokens |
|----------|---------|----------------|
| `background` | Surface colors | `bg-primary`, `bg-secondary`, `bg-tertiary` |
| `foreground` | Text colors | `text-primary`, `text-secondary`, `text-muted` |
| `primary` | Brand/accent colors | `accent`, `link`, `focus` |
| `border` | Edge/divider colors | `border`, `divider`, `outline` |
| `success` | Positive states | `success`, `valid`, `complete` |
| `warning` | Caution states | `warning`, `caution` |
| `danger` | Negative states | `error`, `destructive`, `invalid` |
| `spacing` | Layout spacing | `spacing-sm`, `spacing-md`, `spacing-lg` |
| `radius` | Border radius | `radius-sm`, `radius-md`, `radius-lg` |
| `typography` | Font properties | `font-sans`, `font-mono`, `font-size-*` |
| `elevation` | Shadows/depth | `shadow-sm`, `shadow-md`, `shadow-lg` |

### Implementation Pattern

Components should use the `var(--theme-*, default)` pattern internally:

```css
:host {
  /* External override takes precedence, falls back to internal default */
  --internal-bg: var(--mde-theme-bg-primary, #ffffff);
}

:host([theme='dark']) {
  --internal-bg: var(--mde-theme-bg-primary, #0f172a);
}

.content {
  background: var(--internal-bg);
}
```

This pattern ensures:

- Components work standalone with sensible defaults
- External theming overrides when properties are set
- Both light and dark modes are supported

## Theme Bindings (Application Side)

Applications declare how to map their tokens to component theming APIs:

```json
{
  "themeBindings": [
    {
      "component": "@grokify/markdown-editor",
      "specUrl": "https://github.com/grokify/markdown-editor/blob/main/design-system.json",
      "themeMode": "dark",
      "mappings": [
        { "from": "plexus-dark", "to": "bg-primary" },
        { "from": "plexus-slate", "to": "bg-secondary" },
        { "from": "plexus-cyan", "to": "accent" },
        { "from": "foreground", "to": "text-primary" }
      ]
    }
  ]
}
```

### Binding Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `component` | string | Yes | Component package name or identifier |
| `specUrl` | string | No | URL to component's design-system.json |
| `themeMode` | string | No | Default theme mode (`light`, `dark`, `auto`) |
| `mappings` | array | Yes | Token mapping definitions |
| `strategy` | string | No | Mapping strategy (`explicit`, `semantic`, `hybrid`) |

### Mapping Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `from` | string | Yes | Application token ID (from `foundations.colors`) |
| `to` | string | Yes | Component token ID (from `themingContract.tokens`) |

### Mapping Strategies

| Strategy | Description | Use Case |
|----------|-------------|----------|
| `explicit` | Manual `from`/`to` mapping for each token | Full control, custom branding |
| `semantic` | Auto-match by `semantic` field | Quick integration, standard colors |
| `hybrid` | Semantic defaults + explicit overrides | Balance of automation and control |

## Generated Output

### CSS Generation

```bash
dss generate --css --bindings
```

Produces:

```css
/* markdown-editor theme bindings */
markdown-editor {
  --mde-theme-bg-primary: var(--plexus-dark);
  --mde-theme-bg-secondary: var(--plexus-slate);
  --mde-theme-accent: var(--plexus-cyan);
  --mde-theme-text-primary: var(--plexus-foreground);
}

/* Light theme override */
[data-theme="light"] markdown-editor {
  --mde-theme-bg-primary: #ffffff;
  --mde-theme-bg-secondary: #f8fafc;
  --mde-theme-accent: #3b82f6;
  --mde-theme-text-primary: #1e293b;
}
```

### TypeScript Generation

```bash
dss generate --typescript --bindings
```

Produces:

```typescript
export const markdownEditorTheme = {
  dark: {
    '--mde-theme-bg-primary': 'var(--plexus-dark)',
    '--mde-theme-bg-secondary': 'var(--plexus-slate)',
    '--mde-theme-accent': 'var(--plexus-cyan)',
    '--mde-theme-text-primary': 'var(--plexus-foreground)',
  },
  light: {
    '--mde-theme-bg-primary': '#ffffff',
    '--mde-theme-bg-secondary': '#f8fafc',
    '--mde-theme-accent': '#3b82f6',
    '--mde-theme-text-primary': '#1e293b',
  },
} as const;
```

## Validation

### Contract Validation (Component CI)

```bash
dss contract validate
```

Checks:

- All tokens have `defaultLight` and `defaultDark`
- CSS properties follow prefix convention
- Semantic values are from allowed set
- No duplicate token IDs

### Binding Validation (Application CI)

```bash
dss bindings validate --warn-unbound
```

Checks:

- Referenced components have accessible specs
- All `to` tokens exist in component contract
- All `from` tokens exist in application foundations
- Warns about unmapped contract tokens

## Automatic Semantic Mapping

When application and component tokens share semantics, mappings can be generated:

```bash
dss bind @grokify/markdown-editor --strategy=semantic
```

The tool:

1. Fetches component's `themingContract`
2. Matches app tokens by `semantic` field
3. Generates `themeBindings` section
4. Warns about unmatched tokens

```
✓ Matched: bg-primary ← plexus-dark (semantic: background)
✓ Matched: accent ← plexus-cyan (semantic: primary)
⚠ Unmatched: bg-tertiary (no app token with semantic: background)
  Suggestion: Add mapping or use plexus-slate
```

## Complete Example

### Component: markdown-editor/design-system.json

```json
{
  "meta": {
    "name": "Markdown Editor",
    "version": "1.0.0"
  },
  "themingContract": {
    "prefix": "--mde-theme",
    "tokens": [
      {
        "id": "bg-primary",
        "cssProperty": "--mde-theme-bg-primary",
        "semantic": "background",
        "description": "Main background color",
        "defaultLight": "#ffffff",
        "defaultDark": "#0f172a"
      },
      {
        "id": "accent",
        "cssProperty": "--mde-theme-accent",
        "semantic": "primary",
        "description": "Primary accent for interactive elements",
        "defaultLight": "#3b82f6",
        "defaultDark": "#60a5fa"
      }
    ]
  }
}
```

### Application: plexusone/design-system.json

```json
{
  "meta": {
    "name": "PlexusOne",
    "version": "1.0.0"
  },
  "foundations": {
    "colors": [
      { "id": "plexus-dark", "value": "#0a0e1a", "semantic": "background" },
      { "id": "plexus-cyan", "value": "#06b6d4", "semantic": "primary" }
    ]
  },
  "themeBindings": [
    {
      "component": "@grokify/markdown-editor",
      "mappings": [
        { "from": "plexus-dark", "to": "bg-primary" },
        { "from": "plexus-cyan", "to": "accent" }
      ]
    }
  ]
}
```

### Generated CSS

```css
markdown-editor {
  --mde-theme-bg-primary: var(--plexus-dark);
  --mde-theme-accent: var(--plexus-cyan);
}
```

## Web Components and Shadow DOM

This theming system is designed specifically for Web Components:

- **Shadow DOM encapsulation** prevents normal CSS from styling internal elements
- **CSS custom properties** are the only values that pierce the shadow boundary
- **Theming contracts** formalize which properties components expose

Without a contract, consumers must read source code to discover themeable properties. With a contract, the theming API is explicit, documented, and machine-readable.

## See Also

- [Component Theming Case Study](../case-studies/component-theming.md) - Real-world example with markdown-editor and PlexusOne
- [Foundations](foundations.md) - Defining design tokens
- [Components](components.md) - Component specification including `llm` context

## Modes and density in theming contracts

`ThemeToken` supports generalized per-mode defaults alongside the
`defaultLight`/`defaultDark` sugar fields:

```json
{
  "id": "background",
  "cssProperty": "--btn-background",
  "defaults": { "light": "#ffffff", "dark": "#0A0E1A", "high-contrast": "#000000" },
  "densitySensitive": true
}
```

`themeBindings.themeMode` accepts any mode the document declares; binding
resolution (explicit, semantic, and inherit strategies) resolves the active
mode's value for referenced color tokens and defaults, preserving the
dark-first fallback when no mode is specified. `densitySensitive` marks
tokens whose values should scale with the active density.
