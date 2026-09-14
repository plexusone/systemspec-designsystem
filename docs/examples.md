# Examples

This page shows real-world examples of design system specifications.

## Minimal Example

The smallest possible design system spec:

```
minimal/
└── meta.json
```

**meta.json:**

```json
{
  "name": "Minimal Design System",
  "version": "1.0.0"
}
```

## Starter Example

A basic design system with colors, typography, and a button:

```
starter/
├── meta.json
├── foundations/
│   ├── colors.json
│   └── typography.json
└── components/
    └── button.json
```

### meta.json

```json
{
  "name": "Starter Design System",
  "version": "1.0.0",
  "description": "A starter design system template"
}
```

### foundations/colors.json

```json
{
  "colors": [
    {
      "id": "background",
      "value": "hsl(0 0% 100%)",
      "semantic": "background",
      "usage": "Page backgrounds"
    },
    {
      "id": "foreground",
      "value": "hsl(240 10% 3.9%)",
      "semantic": "foreground",
      "usage": "Primary text"
    },
    {
      "id": "primary",
      "value": "hsl(240 5.9% 10%)",
      "semantic": "primary",
      "usage": "Primary actions"
    },
    {
      "id": "secondary",
      "value": "hsl(240 4.8% 95.9%)",
      "semantic": "secondary",
      "usage": "Secondary actions"
    },
    {
      "id": "muted",
      "value": "hsl(240 4.8% 95.9%)",
      "semantic": "muted",
      "usage": "Muted backgrounds"
    },
    {
      "id": "destructive",
      "value": "hsl(0 84.2% 60.2%)",
      "semantic": "destructive",
      "usage": "Destructive actions"
    },
    {
      "id": "border",
      "value": "hsl(240 5.9% 90%)",
      "semantic": "border",
      "usage": "Borders and dividers"
    }
  ]
}
```

### foundations/typography.json

```json
{
  "fontFamilies": [
    {
      "id": "sans",
      "name": "Sans Serif",
      "stack": "ui-sans-serif, system-ui, -apple-system, sans-serif",
      "usage": "Primary UI text"
    },
    {
      "id": "mono",
      "name": "Monospace",
      "stack": "ui-monospace, SFMono-Regular, monospace",
      "usage": "Code blocks"
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

### components/button.json

```json
{
  "id": "Button",
  "name": "Button",
  "description": "Interactive button for actions",
  "variants": [
    {
      "id": "default",
      "name": "Default",
      "description": "Primary action",
      "isDefault": true
    },
    {
      "id": "secondary",
      "name": "Secondary",
      "description": "Secondary action"
    },
    {
      "id": "outline",
      "name": "Outline",
      "description": "Outlined button"
    },
    {
      "id": "ghost",
      "name": "Ghost",
      "description": "Minimal button"
    },
    {
      "id": "destructive",
      "name": "Destructive",
      "description": "Dangerous action"
    }
  ],
  "props": [
    {
      "name": "variant",
      "type": "enum",
      "description": "Visual variant"
    },
    {
      "name": "size",
      "type": "enum",
      "description": "Button size"
    },
    {
      "name": "disabled",
      "type": "boolean",
      "description": "Disabled state"
    }
  ],
  "llm": {
    "intent": "Trigger user actions",
    "allowedContexts": ["forms", "dialogs", "toolbars"],
    "antiPatterns": [
      "Multiple primary buttons per view",
      "Using for navigation"
    ],
    "exampleUsage": [
      "<Button>Save</Button>",
      "<Button variant=\"secondary\">Cancel</Button>"
    ]
  }
}
```

## Using the Examples

```bash
# Clone the repo
git clone https://github.com/plexusone/systemspec-designsystem.git
cd systemspec-designsystem/examples/starter

# View info
dss info

# Generate CSS
dss generate --css ./output.css

# Validate (if you have components)
dss validate ./src/components
```

## Real-World Example: ProofMinds

The [ProofMinds](https://github.com/proofminds/proofminds-web) project uses DSS for its design system:

```
proofminds-web/design-system/
├── meta.json
├── principles.json
├── accessibility.json
├── foundations/
│   ├── colors.json      # 19 color tokens
│   ├── typography.json  # Font families, sizes, weights
│   ├── spacing.json     # 13-step spacing scale
│   └── border-radius.json
└── components/
    ├── button.json      # 6 variants
    ├── card.json        # 2 variants
    ├── badge.json       # 5 variants
    └── input.json
```

Key features:

- Professional credentialing platform colors (navy-based)
- WCAG AA accessibility target
- LLM context for AI-assisted development
- Compliance validation against React components
