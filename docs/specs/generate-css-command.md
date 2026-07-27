# Enhancement: `dss generate css` command

**Status:** Proposed  
**Priority:** High  
**Requested by:** a customer style guide team  
**Date:** 2026-07-22  

## Summary

Add a `dss generate css` command that reads any conforming `design-system.json` and emits CSS output suitable for direct consumption by Tailwind v4, MkDocs Material, or plain CSS sites.

## Motivation

Today, teams manually copy token values from `design-system.json` into their CSS files. This creates drift — when the design system updates, downstream sites don't automatically pick up the changes.

The `design-system.json` schema already defines an `output.css` section with `prefix`, `selector`, and `mappings`, and the `profiles` section defines per-surface overrides. But there's no CLI command that actually uses this information to generate output.

### Current workflow (manual, error-prone)

```
design-system.json → human reads values → manually writes CSS → site consumes CSS
```

### Desired workflow (automated)

```
design-system.json → `dss generate css` → CSS file → site consumes CSS
```

## Proposed Solution

### Command

```bash
dss generate css [flags]
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--input`, `-i` | Yes | `design-system.json` | Path to design system JSON file |
| `--format`, `-f` | No | `css` | Output format: `css`, `tailwind` |
| `--profile`, `-p` | No | `website` | Profile to use (applies colorOverrides) |
| `--output`, `-o` | No | stdout | Output file path (writes to stdout if omitted) |

### Output Formats

#### `--format css` (default)

Emits standard CSS custom properties:

```css
:root {
  --brand-dark: #0a0a0a;
  --brand-navy: #0f0f18;
  --brand-green: #00ff00;
  --brand-blue: #3b82f6;
  --brand-text: #f1f5f9;
  --brand-font-sans: 'Inter', system-ui, -apple-system, sans-serif;
  --brand-radius-lg: 0.5rem;
  --brand-space-4: 1rem;
  /* ... */
}
```

The selector and prefix come from `output.css.selector` and `output.css.prefix` in the JSON, or from the profile's `css` section if specified.

#### `--format tailwind`

Emits a Tailwind v4 `@theme` block:

```css
@theme {
  --color-brand-dark: #0a0a0a;
  --color-brand-navy: #0f0f18;
  --color-brand-green: #00ff00;
  --color-brand-blue: #3b82f6;
  --color-brand-text: #f1f5f9;
  --font-sans: 'Inter', system-ui, -apple-system, sans-serif;
  /* ... */
}
```

This file gets imported in the site's `index.css` before `@import "tailwindcss"` and makes all tokens available as Tailwind utilities (`bg-brand-dark`, `text-brand-green`, etc.).

### Profile Support

When `--profile mkdocs` is specified:

1. Start with the default token values
2. Apply `colorOverrides` from the profile
3. Use the profile's `css.selector` and `css.prefix` (if defined)
4. If profile defines `css.mappings`, use those instead of the default `output.css.mappings`

Example with a customer's MkDocs profile:

```bash
dss generate css -i design-system.json --profile mkdocs
```

Output:
```css
[data-md-color-scheme="slate"] {
  --md-default-bg-color: #0f0f18;
  --md-code-bg-color: #1a1a24;
  --md-accent-fg-color: #00ff00;
  --md-default-fg-color: #f1f5f9;
}
```

## Integration Examples

### React/Vite site (Tailwind v4)

```bash
# In package.json scripts or CI
dss generate css -i ../../standards/brand-design-system/design-system.json --format tailwind -o src/theme.css
```

```css
/* src/index.css */
@import "./theme.css";
@import "tailwindcss";
```

### MkDocs site

```bash
# In CI before mkdocs build
dss generate css -i ../../standards/brand-design-system/design-system.json --profile mkdocs -o docs/stylesheets/brand.css
```

```yaml
# mkdocs.yml
extra_css:
  - stylesheets/brand.css
```

### CI Pipeline (keeps sites in sync)

```yaml
generate-theme:
  stage: build
  script:
    - dss generate css -i design-system.json --format tailwind -o src/theme.css
    - dss generate css -i design-system.json --profile mkdocs -o docs/stylesheets/brand.css
```

## Schema Requirements

This feature uses existing schema fields:

- `output.css.prefix` — variable prefix (e.g., `--sav`)
- `output.css.selector` — CSS selector (e.g., `:root`)
- `output.css.mappings` — how token paths map to variable names
- `profiles[name].css.prefix` — per-profile prefix override
- `profiles[name].css.selector` — per-profile selector override
- `profiles[name].css.mappings` — per-profile custom mappings
- `profiles[name].colorOverrides` — token value overrides

No schema changes needed — this is purely a CLI implementation.

## Acceptance Criteria

- [ ] `dss generate css` reads any valid `design-system.json`
- [ ] `--format css` produces valid CSS with `:root` selector by default
- [ ] `--format tailwind` produces valid Tailwind v4 `@theme` block
- [ ] `--profile` applies color overrides and uses profile-specific CSS config
- [ ] Output to stdout by default (pipeable), `--output` writes to file
- [ ] Exit 0 on success, exit 1 on invalid input
- [ ] Works with plexusone-design-system and acme/brand-design-system JSON files

## Related

- `plexusone-design-system` — has `generateDesignSystemCSS()` in TypeScript (nav package), but it's not a CLI command
- `acme/brand-design-system` — first consumer that needs this for both React and MkDocs sites
- Schema: `output.css` section already defines the contract
