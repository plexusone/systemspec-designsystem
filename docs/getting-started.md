# Getting Started

This guide walks you through creating your first design system spec and generating code from it.

## Prerequisites

- Go 1.22 or later
- A React/TypeScript project (optional, for validation)

## Installation

Install the `dss` CLI tool:

```bash
go install github.com/plexusone/systemspec-designsystem/cmd/dss@latest
```

Verify installation:

```bash
dss --version
```

## Create a Design System Specification

### 1. Create Directory Structure

```bash
mkdir -p my-design-system/foundations my-design-system/components
```

### 2. Create Meta File

Create `my-design-system/meta.json`:

```json
{
  "name": "My Design System",
  "version": "1.0.0",
  "description": "A custom design system for my application"
}
```

### 3. Define Colors

Create `my-design-system/foundations/colors.json`:

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
      "usage": "Primary brand color for CTAs and key actions"
    },
    {
      "id": "secondary",
      "value": "hsl(210 40% 96.1%)",
      "semantic": "secondary",
      "usage": "Secondary buttons and backgrounds"
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

### 4. Define Typography

Create `my-design-system/foundations/typography.json`:

```json
{
  "fontFamilies": [
    {
      "id": "sans",
      "name": "Sans Serif",
      "stack": "ui-sans-serif, system-ui, sans-serif",
      "usage": "Primary font for all UI text"
    }
  ],
  "fontSizes": [
    { "id": "xs", "value": "0.75rem" },
    { "id": "sm", "value": "0.875rem" },
    { "id": "base", "value": "1rem" },
    { "id": "lg", "value": "1.125rem" },
    { "id": "xl", "value": "1.25rem" },
    { "id": "2xl", "value": "1.5rem" }
  ],
  "fontWeights": [
    { "id": "normal", "value": 400 },
    { "id": "medium", "value": 500 },
    { "id": "semibold", "value": 600 },
    { "id": "bold", "value": 700 }
  ]
}
```

### 5. Define a Component

Create `my-design-system/components/button.json`:

```json
{
  "id": "Button",
  "name": "Button",
  "description": "Interactive button for triggering actions",
  "variants": [
    {
      "id": "default",
      "name": "Default",
      "description": "Primary action button",
      "isDefault": true
    },
    {
      "id": "secondary",
      "name": "Secondary",
      "description": "Less prominent actions"
    },
    {
      "id": "destructive",
      "name": "Destructive",
      "description": "Dangerous actions like delete"
    }
  ],
  "props": [
    {
      "name": "variant",
      "type": "enum",
      "description": "Visual style variant"
    },
    {
      "name": "disabled",
      "type": "boolean",
      "description": "Whether the button is disabled"
    }
  ],
  "llm": {
    "intent": "Trigger user-initiated actions",
    "allowedContexts": ["forms", "dialogs", "toolbars"],
    "forbiddenContexts": ["inline-text"],
    "antiPatterns": [
      "Multiple primary buttons in the same view",
      "Using buttons for navigation"
    ],
    "exampleUsage": [
      "<Button>Save</Button>",
      "<Button variant=\"secondary\">Cancel</Button>"
    ]
  }
}
```

## Generate Code

### View Design System Info

```bash
cd my-design-system
dss info
```

Output:

```
Design System: My Design System
Version: 1.0.0
Description: A custom design system for my application

Foundations:
  Colors: 5 tokens
  Font Families: 1
  Font Sizes: 6

Components: 1
  - Button (3 variants)

✓ Spec validation passed
```

### Generate CSS

```bash
dss generate --css ./output/styles.css
```

This generates Tailwind v4 compatible CSS:

```css
@import "tailwindcss";

@theme {
  --color-background: hsl(0 0% 100%);
  --color-foreground: hsl(222.2 84% 4.9%);
  --color-primary: hsl(222.2 47.4% 11.2%);
  /* ... */
}
```

### Generate TypeScript Types

```bash
dss generate --types ./output/types.ts
```

This generates typed interfaces:

```typescript
export type ButtonVariant = 'default' | 'secondary' | 'destructive';

export interface ButtonProps {
  variant?: ButtonVariant;
  disabled?: boolean;
  children: React.ReactNode;
}
```

### Generate LLM Context

```bash
dss generate --llm ./DESIGN_CONTEXT.md
```

This generates a markdown file optimized for AI code generation, including:

- Design principles
- Token documentation
- Component specs with anti-patterns
- Accessibility requirements

### Generate NPM Package

Generate a publishable NPM package with framework-specific outputs:

```bash
dss generate --package ./dist --targets all
```

This creates a complete package with:

- `package.json` with proper exports
- CSS custom properties
- Tailwind v4 preset
- ShadCN theme variables
- MkDocs Material theme
- SCSS variables
- JSON and W3C token formats

**Specific Targets:**

```bash
# Only CSS and Tailwind
dss generate --package ./dist --targets css,tailwind

# With custom scope and name
dss generate --package ./dist --scope @myorg --name tokens --targets css,tailwind,shadcn
```

The generated package can be published to NPM and consumed in your projects.

## Validate Implementation

If you have React components, validate them against your spec:

```bash
dss validate ../src/components
```

The validator checks for:

- Hardcoded colors (should use CSS variables)
- Invalid variant values
- Missing accessibility attributes
- Anti-pattern violations

## Next Steps

- Read the [CLI Reference](cli.md) for all available commands
- Explore the [Specification](specification/index.md) for all available fields
- See [Examples](examples.md) for real-world design systems
- Set up the [MCP Server](mcp-server.md) for AI assistant integration
- Learn about [embedding specs](sdk.md#from-embedded-filesystem) for single-binary distribution
