# MCP Server

The `dss-mcp` command exposes design system operations as an MCP (Model Context Protocol) server, enabling AI assistants like Claude to query and validate against your design system specification.

## Overview

The MCP server provides 37 tools organized into nine categories:

| Category | Tools | Purpose |
|----------|-------|---------|
| Spec Reading | 7 | Query components, tokens, patterns, and metadata |
| Guidance | 4 | Generate prompts, get variants, props, and anti-patterns |
| Validation | 4 | Validate files against the design system |
| Fix | 6 | Auto-fix design system violations |
| Lint | 3 | Check spec completeness and agent-readiness |
| Validators | 4 | External validator discovery and delegation |
| Compliance | 5 | Compliance reporting and release gates |
| Accessibility | 4 | A11y requirements, anti-patterns, contrast suggestions |

## Installation

```bash
# Build from source
go install github.com/plexusone/systemspec-designsystem/cmd/dss-mcp@latest

# Or build locally
go build -o dss-mcp ./cmd/dss-mcp
```

## Usage

### Basic Usage

```bash
# Start MCP server with your design system spec
dss-mcp --spec ./design-system/

# Or with a single spec file
dss-mcp --spec ./design-system.yaml
```

### With Browser Validation (w3pilot)

```bash
# Enable w3pilot browser tools for visual validation
dss-mcp --spec ./design-system/ --browser
```

This adds 169 additional browser automation tools from w3pilot for:

- Screenshot comparison
- Computed style verification
- Accessibility tree inspection
- Visual regression testing

## Claude Desktop Integration

Add to your Claude Desktop configuration:

### macOS

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "design-system": {
      "command": "dss-mcp",
      "args": ["--spec", "/path/to/your/design-system"]
    }
  }
}
```

### With Browser Tools

```json
{
  "mcpServers": {
    "design-system": {
      "command": "dss-mcp",
      "args": ["--spec", "/path/to/your/design-system", "--browser"]
    }
  }
}
```

## Available Tools

### Spec Reading Tools

#### `get_component`

Get a component definition including variants, props, states, events, and accessibility requirements.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | Yes | Component ID (e.g., "button", "input") |

**Example:**

```json
{
  "id": "button"
}
```

#### `list_components`

List all components with ID, name, category, and description.

**Parameters:** None

#### `get_token`

Get a design token value.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `type` | string | Yes | Token type: color, spacing, typography, elevation, borderRadius, breakpoint |
| `name` | string | Yes | Token name/ID (e.g., "primary-500") |

#### `list_tokens`

List all tokens of a given type.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `type` | string | Yes | Token type |

#### `get_pattern`

Get a pattern definition including components, layout, and variations.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `id` | string | Yes | Pattern ID (e.g., "form-validation") |

#### `list_patterns`

List all patterns with ID, name, category, and description.

**Parameters:** None

#### `get_meta`

Get design system metadata including name, version, and description.

**Parameters:** None

### Guidance Tools

#### `generate_prompt`

Generate an LLM context prompt for the design system.

**Parameters:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `format` | string | No | "markdown" | Output format: markdown, xml |
| `include_foundations` | boolean | No | true | Include design tokens |
| `include_components` | boolean | No | true | Include component definitions |
| `include_patterns` | boolean | No | true | Include pattern recommendations |
| `include_accessibility` | boolean | No | true | Include accessibility requirements |
| `include_anti_patterns` | boolean | No | true | Include anti-patterns to avoid |

#### `get_variants`

Get all valid variants for a component.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component_id` | string | Yes | Component ID |

#### `get_props`

Get prop definitions for a component including types, defaults, and constraints.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component_id` | string | Yes | Component ID |

#### `get_anti_patterns`

Get anti-patterns to avoid for a component.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component_id` | string | Yes | Component ID |

### Validation Tools

#### `validate_file`

Validate a file against the design system spec.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to validate |
| `rules` | array | No | Specific rules to check (default: all) |
| `include_context` | boolean | No | Include code snippets in violations |

**Available Rules:**

- `no-hardcoded-colors` - Check for hardcoded hex/rgb/hsl colors
- `use-spacing-scale` - Check for non-standard spacing values
- `img-alt-required` - Check for images without alt attributes
- `button-accessible-name` - Check for icon buttons without aria-label
- `valid-variant` - Check for unknown variant values
- `single-primary-button` - Check for multiple primary buttons
- `no-nested-cards` - Check for nested card components

#### `validate_directory`

Validate all files in a directory.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Directory path to validate |
| `extensions` | array | No | File extensions to check (default: .tsx, .jsx, .ts, .js, .css) |
| `rules` | array | No | Specific rules to check |

#### `check_colors`

Check a file for hardcoded colors that should use design tokens.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to check |

#### `check_spacing`

Check a file for hardcoded spacing values that should use the spacing scale.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to check |

### Fix Tools

#### `fix_file`

Fix design system violations in a file. Replaces hardcoded colors with CSS variables, fixes spacing values, and adds missing accessibility attributes.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `rules` | array | No | Specific rules to fix (default: all) |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `suggest_fixes`

Suggest fixes for design system violations without applying them. Returns the original content, suggested fixes, and what the fixed content would look like.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to analyze |
| `rules` | array | No | Specific rules to check (default: all) |

#### `fix_colors`

Fix hardcoded colors in a file by replacing them with CSS variables from the design system.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `fix_spacing`

Fix hardcoded spacing values in a file by replacing them with spacing tokens.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `fix_accessibility`

Fix accessibility issues in a file. Adds missing alt attributes to images and aria-label attributes to icon-only buttons.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the file to fix |
| `dry_run` | boolean | No | If true, return fixes without applying them |

#### `fix_directory`

Fix design system violations in all files in a directory.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `path` | string | Yes | Path to the directory to fix |
| `rules` | array | No | Specific rules to fix (default: all) |
| `dry_run` | boolean | No | If true, return fixes without applying them |

### Lint Tools

#### `lint_spec`

Lint the design system spec for completeness and best practices. Returns a score (0-100), coverage metrics, and detailed issues.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `rules` | array | No | Specific rules to check (default: all) |
| `min_score` | number | No | Minimum acceptable score (0-100) |

**Available Rules:**

- `meta-required` - Design system must have name and version
- `component-has-variants` - Components should define variants
- `component-has-props` - Components should define props
- `component-has-llm-context` - Components should have LLM context
- `llm-has-intent` - LLM context must have intent field
- `llm-has-anti-patterns` - LLM context should document anti-patterns
- `llm-has-allowed-contexts` - LLM context should specify allowed contexts
- `tokens-have-descriptions` - Tokens should have descriptions
- `token-references-valid` - Token references must resolve
- `no-orphan-tokens` - Tokens should be referenced by components
- `component-uses-valid` - Component uses references must be valid
- `accessibility-defined` - Design system should define accessibility requirements
- `theming-contract-valid` - Theming contracts must be properly configured
- `validators-configured` - External validators should be configured
- `validator-tool-required` - Validators must specify a tool name
- `validator-type-valid` - Validator type must be valid

#### `list_lint_rules`

List all available lint rules with descriptions.

**Parameters:** None

#### `check_agent_readiness`

Check if the design system spec is ready for AI agent code generation. Verifies LLM context completeness and validator configuration.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component_id` | string | No | Check readiness for a specific component |

### Validator Tools

Tools for discovering and delegating to external validators (e.g., agent-a11y for WCAG, spectral for API style).

#### `list_validators`

List all configured external validators.

**Parameters:** None

**Returns:**

```json
{
  "configured": true,
  "validators": [
    {
      "domain": "accessibility",
      "tool": "agent-a11y",
      "type": "mcp",
      "required": true
    }
  ]
}
```

#### `get_validator`

Get configuration for a specific validator by domain.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `domain` | string | Yes | Validator domain: `accessibility`, `api`, or custom |

#### `get_validation_requirements`

Get all validation requirements defined in the spec with their corresponding validators.

**Parameters:** None

#### `get_validator_invocation`

Get invocation details for a validator (MCP command, CLI command, etc.).

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `domain` | string | Yes | Validator domain |

### Compliance Tools

Tools for the complete release workflow: validate → fix → verify → release.

#### `generate_compliance_report`

Generate a comprehensive compliance report for release.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `directory` | string | Yes | Directory to validate |
| `min_score` | number | No | Minimum acceptable score (default: 80) |
| `blocking_categories` | array | No | Categories that block release (default: colors, accessibility) |
| `include_issues` | boolean | No | Include detailed issue list (default: true) |

**Returns:**

```json
{
  "status": "pass",
  "score": 92,
  "categories": [
    {"category": "colors", "status": "pass", "score": 100},
    {"category": "accessibility", "status": "warn", "score": 85}
  ],
  "summary": {
    "totalChecks": 15,
    "passed": 12,
    "errors": 0,
    "warnings": 3
  }
}
```

#### `check_release_gate`

Check if the codebase passes the release gate. Returns approved/rejected with reasons.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `directory` | string | Yes | Directory to validate |
| `min_score` | number | No | Minimum acceptable score (default: 80) |
| `require_zero_errors` | boolean | No | Require zero error-level issues (default: true) |
| `allow_warnings` | boolean | No | Allow release with warnings (default: true) |

**Returns:**

```json
{
  "approved": true,
  "reason": "All checks passed",
  "score": 92,
  "certificate": {
    "id": "a3f2c8b1",
    "timestamp": "2024-01-15T10:30:00Z",
    "hash": "sha256:..."
  }
}
```

#### `run_fix_loop`

Run an automated fix-validate loop until convergence or max iterations.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `directory` | string | Yes | Directory to fix |
| `max_iterations` | number | No | Maximum fix-validate cycles (default: 3) |
| `dry_run` | boolean | No | Preview fixes without applying (default: false) |
| `stop_on_regression` | boolean | No | Stop if fixes introduce new violations (default: true) |
| `rules` | array | No | Specific rules to fix |

**Returns:**

```json
{
  "converged": true,
  "iterations": 2,
  "initialViolations": 8,
  "finalViolations": 0,
  "fixedCount": 8,
  "status": "pass"
}
```

#### `fix_and_verify`

Fix a single file and verify the fix resolved the violations.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `file` | string | Yes | File to fix |
| `rules` | array | No | Specific rules to fix |
| `dry_run` | boolean | No | Preview fixes without applying |

**Returns:**

```json
{
  "violationsBefore": 3,
  "fixesApplied": 3,
  "violationsAfter": 0,
  "verified": true,
  "reason": "all violations fixed"
}
```

#### `get_compliance_certificate`

Generate a compliance certificate after validation. Provides proof of compliance for release artifacts.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `directory` | string | Yes | Directory that was validated |

### Accessibility Integration Tools

Tools for querying accessibility requirements, anti-patterns, and token contrast. Designed for integration with agent-a11y.

#### `get_accessibility_requirements`

Get accessibility requirements for a component including required props, keyboard interactions, focus management, and WCAG criteria.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component` | string | Yes | Component ID (e.g., "button", "input", "modal") |

**Returns:**

```json
{
  "componentId": "button",
  "component": "Button",
  "requiredProps": [
    { "name": "aria-label", "when": "icon-only", "example": "aria-label=\"Close\"" }
  ],
  "keyboard": {
    "interactions": {
      "Enter": "activates button",
      "Space": "activates button"
    }
  },
  "focus": {
    "visible": true,
    "order": "natural"
  },
  "colorContrast": {
    "text": "4.5:1",
    "largeText": "3.0:1",
    "uiComponents": "3.0:1"
  },
  "wcagCriteria": ["2.1.1", "2.4.7", "4.1.2"],
  "role": "button"
}
```

#### `get_a11y_anti_patterns`

Get accessibility anti-patterns to avoid for a component or rule. Returns bad examples, good examples, and WCAG criteria.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component` | string | No | Component ID to get anti-patterns for |
| `rule_id` | string | No | Rule ID like "color-contrast", "missing-label", "keyboard", "focus" |

**Returns:**

```json
{
  "componentId": "button",
  "antiPatterns": [
    {
      "id": "missing-button-name",
      "description": "Buttons without accessible names",
      "badExample": "<button><svg>...</svg></button>",
      "goodExample": "<button aria-label=\"Close dialog\"><svg>...</svg></button>",
      "wcagCriteria": ["4.1.2"],
      "components": ["button", "icon-button"]
    }
  ]
}
```

#### `suggest_contrast_token`

Suggest color tokens that meet contrast requirements against a background color.

**Parameters:**

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `background` | string | Yes | - | Background color as hex (e.g., "#ffffff") |
| `min_contrast` | number | No | 4.5 | Minimum contrast ratio (4.5 for AA normal text) |

**Returns:**

```json
[
  {
    "token": "gray-900",
    "value": "#1a1a2e",
    "contrastRatio": 12.6,
    "meetsAA": true,
    "meetsAAA": true,
    "meetsAALarge": true
  },
  {
    "token": "primary-700",
    "value": "#1d4ed8",
    "contrastRatio": 4.8,
    "meetsAA": true,
    "meetsAAA": false,
    "meetsAALarge": true
  }
]
```

#### `get_component_fix_context`

Get full context needed to fix accessibility issues in a component including file patterns, props to add, styles to check, and available tokens.

**Parameters:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `component` | string | Yes | Component ID (e.g., "button", "input") |
| `issue_type` | string | Yes | Issue type: "color-contrast", "missing-label", "keyboard", "focus" |

**Returns:**

```json
{
  "componentId": "button",
  "issueType": "color-contrast",
  "filePattern": "src/components/Button.{tsx,vue,svelte}",
  "propsToAdd": [],
  "stylesToCheck": [
    { "property": "color", "token": "color.text.*" },
    { "property": "background-color", "token": "color.bg.*" }
  ],
  "tokensAvailable": {
    "colors": ["primary", "secondary", "gray-900", "white"],
    "spacing": ["1", "2", "4", "8"]
  },
  "relatedComponents": ["icon", "tooltip"],
  "antiPatterns": ["Using color alone to indicate errors"]
}
```

## Example Workflows

### Component Implementation

1. Ask Claude to implement a button component
2. Claude uses `get_component` to get the button spec
3. Claude uses `get_variants` to understand available variants
4. Claude implements the component
5. Claude uses `validate_file` to check compliance

### Design Token Lookup

```
User: "What color should I use for primary actions?"

Claude: [calls get_token with type="color", name="primary-500"]
"The primary action color is #3B82F6 (primary-500)"
```

### Code Review

```
User: "Review my Button.tsx for design system compliance"

Claude: [calls validate_file with path="./src/components/Button.tsx"]
"Found 2 warnings:
- Line 15: Hardcoded color '#ff0000' - use CSS variable
- Line 23: Missing aria-label on icon button"
```

### Auto-Fix Violations

```
User: "Fix the design system violations in my components folder"

Claude: [calls fix_directory with path="./src/components", dry_run=true]
"Found 5 fixable issues across 3 files:
- Button.tsx: 2 color fixes, 1 accessibility fix
- Card.tsx: 1 spacing fix
- Input.tsx: 1 color fix

Would you like me to apply these fixes?"

User: "Yes, apply them"

Claude: [calls fix_directory with path="./src/components"]
"Applied 5 fixes:
- Replaced #ff0000 with var(--color-error)
- Added aria-label to icon button
- Replaced 15px with var(--spacing-4)
..."
```

### Complete Release Workflow

```
User: "Prepare the codebase for release"

Claude: [calls lint_spec]
"Spec completeness: 92/100
- All components have LLM context ✓
- Accessibility validator configured ✓"

Claude: [calls run_fix_loop with directory="./src"]
"Fix-validate loop completed:
- Initial violations: 12
- Fixes applied: 10
- Final violations: 2 (unfixable - need human review)
- Converged in 2 iterations"

Claude: [calls check_release_gate with directory="./src"]
"Release gate: APPROVED
- Score: 88/100
- Blocking issues: 0
- Warnings: 2 (non-blocking)
- Certificate: a3f2c8b1"

User: "What are the warnings?"

Claude: [reviews remaining issues]
"2 semantic issues requiring human review:
1. Card.tsx:45 - Consider if this color choice matches brand guidelines
2. Form.tsx:78 - Multiple buttons in close proximity

These are non-blocking. Ready to release."
```

### Check Agent Readiness

```
User: "Is our design system ready for AI code generation?"

Claude: [calls check_agent_readiness]
"Agent Readiness: 85%

✓ 12/14 components have LLM context
✓ Intent defined for all components with LLM context
✓ Anti-patterns documented
✓ Accessibility validator configured

⚠ Missing LLM context:
- Tooltip: Add intent and allowedContexts
- Divider: Add intent field

Recommendation: Add LLM context to remaining components for best AI generation results."
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     AI Agents (Claude, etc.)                     │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ MCP Protocol (stdio)
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    dss-mcp MCP Server                            │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                  omniskill Runtime                          ││
│  │  ┌──────────────────────────┐  ┌──────────────────────────┐││
│  │  │ designsystem skill       │  │ w3pilot skill            │││
│  │  │ (37 native Go tools)     │  │ (169 browser tools)      │││
│  │  │ - Spec reading (7)       │  │                          │││
│  │  │ - Guidance (4)           │  │                          │││
│  │  │ - Validation (4)         │  │                          │││
│  │  │ - Fix (6)                │  │                          │││
│  │  │ - Lint (3)               │  │                          │││
│  │  │ - Validators (4)         │  │                          │││
│  │  │ - Compliance (5)         │  │                          │││
│  │  │ - Accessibility (4)      │  │                          │││
│  │  └──────────┬───────────────┘  └────────────┬─────────────┘││
│  └─────────────┼───────────────────────────────┼───────────────┘│
└────────────────┼───────────────────────────────┼────────────────┘
                 │                               │
                 ▼                               ▼
┌────────────────────────────────────┐  ┌───────────────────────┐
│      sdk/go/ (Service Layer)       │  │  w3pilot MCP Server   │
│  - Component operations            │  │  (subprocess)         │
│  - Token operations                │  └───────────────────────┘
│  - Validation & Fix                │
│  - Compliance & Release            │
│  - Spec linting                    │
│  - Accessibility integration       │
└────────────────────────────────────┘
```

## Troubleshooting

### Server Won't Start

```bash
# Check if spec path is correct
ls -la /path/to/your/design-system

# Run with verbose output
dss-mcp --spec ./design-system/ 2>&1 | head -20
```

### Tool Errors

```bash
# Test with MCP inspector
npx @anthropic/mcp-inspector dss-mcp --spec ./design-system/
```

### Browser Tools Not Available

Ensure w3pilot is installed and in your PATH:

```bash
which w3pilot
w3pilot --version
```

## Embedding Design System Specifications

For distribution, you can create a custom MCP server binary that embeds your design system spec using Go's `embed` package. This eliminates the need for external spec files.

### Example: Embedded MCP Server

```go
package main

import (
    "context"
    "embed"
    "fmt"
    "io/fs"
    "os"

    "github.com/modelcontextprotocol/go-sdk/mcp"
    dss "github.com/plexusone/systemspec-designsystem/sdk/go"
    "github.com/plexusone/systemspec-designsystem/skills/designsystem"
    "github.com/plexusone/systemspec-designsystem/internal/omniskill/mcp/server"
)

//go:embed spec/*
var specFS embed.FS

func main() {
    if err := run(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func run() error {
    ctx := context.Background()

    // Load from embedded filesystem
    sub, err := fs.Sub(specFS, "spec")
    if err != nil {
        return fmt.Errorf("failed to get spec subdirectory: %w", err)
    }

    ds, err := dss.LoadDesignSystemFromFS(sub)
    if err != nil {
        return fmt.Errorf("failed to load design system: %w", err)
    }

    // Create service and skill
    service := dss.NewService(ds)
    dsSkill := designsystem.New(service)
    if err := dsSkill.Init(ctx); err != nil {
        return fmt.Errorf("failed to initialize skill: %w", err)
    }
    defer dsSkill.Close()

    // Create and run MCP server
    rt := server.New(&mcp.Implementation{
        Name:    "my-design-system",
        Version: "1.0.0",
    }, nil)
    rt.RegisterSkill(dsSkill)

    return rt.ServeStdio(ctx)
}
```

### Project Structure

```
my-design-system-mcp/
├── go.mod
├── main.go
└── spec/
    ├── meta.json
    ├── foundations/
    │   └── colors.json
    └── components/
        └── button.json
```

### Building and Distributing

```bash
# Build single binary with embedded spec
go build -o my-design-system-mcp .

# Distribute just the binary - no external files needed
./my-design-system-mcp
```

### Claude Desktop Configuration

No `--spec` flag needed since the spec is embedded:

```json
{
  "mcpServers": {
    "my-design-system": {
      "command": "/path/to/my-design-system-mcp"
    }
  }
}
```

## Related

- [CLI Reference](cli.md) - Command-line interface
- [SDK Reference](sdk.md) - Go SDK for programmatic access
- [Getting Started](getting-started.md) - Quick start guide
