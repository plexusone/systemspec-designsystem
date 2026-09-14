# Foundations

Foundations define the design tokens that form the visual foundation of your design system: colors, typography, spacing, and more.

## File Location

```
my-design-system/
└── foundations/
    ├── colors.json
    ├── typography.json
    ├── spacing.json
    └── border-radius.json
```

## Colors

### Schema

```json
{
  "colors": [
    {
      "id": "primary",
      "value": "hsl(222.2 47.4% 11.2%)",
      "semantic": "primary",
      "usage": "Primary brand color for CTAs"
    }
  ]
}
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique token identifier |
| `value` | string | Yes | Color value (hex, hsl, rgb) |
| `semantic` | string | No | Semantic category |
| `usage` | string | No | Usage description |

### Example

```json
{
  "colors": [
    {
      "id": "background",
      "value": "hsl(0 0% 100%)",
      "semantic": "background",
      "usage": "Primary page background"
    },
    {
      "id": "foreground",
      "value": "hsl(222.2 84% 4.9%)",
      "semantic": "foreground",
      "usage": "Primary text color"
    },
    {
      "id": "primary",
      "value": "hsl(222.2 47.4% 11.2%)",
      "semantic": "primary",
      "usage": "CTAs and key interactive elements"
    },
    {
      "id": "destructive",
      "value": "hsl(0 84.2% 60.2%)",
      "semantic": "destructive",
      "usage": "Error states and destructive actions"
    }
  ]
}
```

## Typography

### Schema

```json
{
  "fontFamilies": [...],
  "fontSizes": [...],
  "fontWeights": [...],
  "lineHeights": [...],
  "typeScale": [...]
}
```

### Font Families

```json
{
  "fontFamilies": [
    {
      "id": "sans",
      "name": "Sans Serif",
      "stack": "ui-sans-serif, system-ui, sans-serif",
      "usage": "Primary UI text"
    },
    {
      "id": "mono",
      "name": "Monospace",
      "stack": "ui-monospace, SFMono-Regular, monospace",
      "usage": "Code and technical content"
    }
  ]
}
```

### Font Sizes

```json
{
  "fontSizes": [
    { "id": "xs", "value": "0.75rem" },
    { "id": "sm", "value": "0.875rem" },
    { "id": "base", "value": "1rem" },
    { "id": "lg", "value": "1.125rem" },
    { "id": "xl", "value": "1.25rem" },
    { "id": "2xl", "value": "1.5rem" },
    { "id": "3xl", "value": "1.875rem" },
    { "id": "4xl", "value": "2.25rem" }
  ]
}
```

### Font Weights

```json
{
  "fontWeights": [
    { "id": "normal", "value": 400 },
    { "id": "medium", "value": 500 },
    { "id": "semibold", "value": 600 },
    { "id": "bold", "value": 700 }
  ]
}
```

## Spacing

### Schema

```json
{
  "baseUnit": "0.25rem",
  "scale": [...]
}
```

### Example

```json
{
  "baseUnit": "0.25rem",
  "scale": [
    { "id": "0", "value": "0" },
    { "id": "1", "value": "0.25rem" },
    { "id": "2", "value": "0.5rem" },
    { "id": "3", "value": "0.75rem" },
    { "id": "4", "value": "1rem" },
    { "id": "6", "value": "1.5rem" },
    { "id": "8", "value": "2rem" },
    { "id": "12", "value": "3rem" },
    { "id": "16", "value": "4rem" },
    { "id": "24", "value": "6rem" }
  ]
}
```

## Border Radius

```json
{
  "borderRadius": [
    { "id": "none", "value": "0" },
    { "id": "sm", "value": "0.125rem" },
    { "id": "md", "value": "0.375rem" },
    { "id": "lg", "value": "0.5rem" },
    { "id": "xl", "value": "0.75rem" },
    { "id": "full", "value": "9999px" }
  ]
}
```

## Generated Output

When you run `dss generate --css`, foundations are converted to CSS custom properties:

```css
@theme {
  --color-background: hsl(0 0% 100%);
  --color-foreground: hsl(222.2 84% 4.9%);
  --color-primary: hsl(222.2 47.4% 11.2%);

  --font-sans: ui-sans-serif, system-ui, sans-serif;
  --font-mono: ui-monospace, SFMono-Regular, monospace;

  --text-xs: 0.75rem;
  --text-sm: 0.875rem;
  --text-base: 1rem;

  --spacing-0: 0;
  --spacing-1: 0.25rem;
  --spacing-2: 0.5rem;

  --radius-sm: 0.125rem;
  --radius-md: 0.375rem;
  --radius-lg: 0.5rem;
}
```

## Modes

A design system may declare discrete presentation modes at the document level:

```json
{ "modes": ["light", "dark", "high-contrast"] }
```

Color tokens provide per-mode values through a generalized `modes` map (the
`lightModeValue`/`darkModeValue` fields remain as sugar for the `light` and
`dark` entries; explicit map entries win):

```json
{
  "id": "surface",
  "value": "#0A0E1A",
  "modes": { "light": "#FFFFFF", "dark": "#0A0E1A" }
}
```

The `mode-completeness` lint rule (`dss lint-spec`) verifies that every
mode-aware token covers every declared mode and that no token uses an
undeclared mode. CSS generation emits one `[data-mode="<mode>"]` override
block per mode; switch modes at runtime by setting that attribute.

## Densities

Densities define spacing variants (e.g. comfortable vs. compact):

```json
{
  "densities": [
    { "id": "comfortable", "scale": 1.0 },
    { "id": "compact", "scale": 0.75, "spacingOverrides": { "4": "0.65rem" } }
  ]
}
```

`scale` is a spacing multiplier; `spacingOverrides` replaces specific spacing
tokens where a linear multiplier is not enough. Generated CSS publishes a
`--density` custom property (`1` at `:root`, the variant's scale inside each
`[data-density="<id>"]` block) plus any spacing overrides — components
consume it via `calc(<base> * var(--density, 1))`.
