# systemspec-designsystem

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/systemspec-designsystem/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/systemspec-designsystem/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/systemspec-designsystem/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/systemspec-designsystem/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/systemspec-designsystem/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/systemspec-designsystem/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/systemspec-designsystem
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/systemspec-designsystem
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/systemspec-designsystem
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fsystemspec-designsystem
 [loc-svg]: https://tokei.rs/b1/github/plexusone/systemspec-designsystem
 [repo-url]: https://github.com/plexusone/systemspec-designsystem
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/systemspec-designsystem/blob/main/LICENSE

A declarative, machine-readable specification for defining design systems as code, built for AI-native development.

> **Note:** This project is part of the PlexusOne `systemspec-{domain}` family of machine-readable system definitions (alongside `systemspec-apistyle`, `systemspec-crypto`, and `systemspec-deploy`). It was formerly named `design-system-spec` and was renamed in September 2026; earlier releases were published under the old module path `github.com/plexusone/design-system-spec`.

## Primary Goals

DSS exists to solve two problems in AI-assisted development:

### 1. AI Agents Build Compliant UI

When an AI agent (Claude, Copilot, etc.) generates UI code, it should automatically follow the design system—using correct colors, spacing, component variants, and patterns.

**How DSS helps:** The `dss generate --llm` command produces a structured context document containing:
- Design principles and constraints
- Token values with semantic meanings
- Component specs with `allowedContexts` and `antiPatterns`
- Accessibility requirements

Include this in your AI context, and generated code follows the design system.

### 2. AI Agents Review Code for Compliance

After code is written (by humans or AI), an agent should be able to review it and report design system violations.

**How DSS helps:** The `dss validate` command checks implementations against the spec:
- Hardcoded colors → should use CSS variables
- Invalid variant values → must match component spec
- Missing accessibility attributes → required by spec
- Anti-pattern violations → flagged in component spec

JSON output (`--json`) enables programmatic integration with CI or agent workflows.

## Assessment: How Well Does DSS Achieve These Goals?

| Goal | Effectiveness | Notes |
|------|---------------|-------|
| **AI builds compliant UI** | **High** | LLM context with intent, examples, and anti-patterns significantly improves AI output quality |
| **AI reviews for compliance** | **Medium-High** | Catches token violations and accessibility issues; cannot verify visual correctness or behavioral logic |

### What DSS Catches

- Hardcoded colors, spacing values
- Invalid component variants
- Missing `alt`, `aria-label` attributes
- Anti-patterns (multiple primary buttons, nested cards)
- Non-standard token values

### What DSS Cannot Catch

- Visual correctness (does it *look* right?)
- Interaction behavior (does hover/focus work?)
- Semantic appropriateness (is this the *right* component for the use case?)
- Layout and responsive issues

For visual validation, combine DSS with visual regression testing (Chromatic, Percy).

## Installation

```bash
# CLI
go install github.com/plexusone/systemspec-designsystem/cmd/dss@latest

# Go SDK
go get github.com/plexusone/systemspec-designsystem
```

## Quick Start

### 1. Create a Spec

```
my-design-system/
├── meta.json
├── foundations/
│   └── colors.json
└── components/
    └── button.json
```

**meta.json:**
```json
{
  "name": "My Design System",
  "version": "1.0.0"
}
```

**components/button.json:**
```json
{
  "id": "Button",
  "name": "Button",
  "variants": [
    { "id": "default", "isDefault": true },
    { "id": "secondary" },
    { "id": "destructive" }
  ],
  "llm": {
    "intent": "Trigger user actions",
    "allowedContexts": ["forms", "dialogs", "toolbars"],
    "antiPatterns": ["Multiple primary buttons per view"]
  }
}
```

### 2. Generate LLM Context

```bash
dss generate --llm ./DESIGN_CONTEXT.md
```

Add `DESIGN_CONTEXT.md` to your AI assistant's context (Claude Projects, Copilot instructions, etc.).

### 3. Validate Implementations

```bash
# Human-readable
dss validate ./src/components

# JSON for CI/agents
dss validate --json ./src/components
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `dss info` | Display design system metadata |
| `dss generate` | Generate CSS, TypeScript types, LLM prompt |
| `dss validate <dir>` | Validate component implementations |
| `dss lint-spec` | Lint spec for completeness and agent-readiness |
| `dss bind` | Generate theme bindings CSS/TypeScript/SCSS |
| `dss contract show` | Display component theming contract |
| `dss contract validate` | Validate all theming contracts |

## MCP Server

The `dss-mcp` command exposes your design system as an MCP (Model Context Protocol) server, enabling AI assistants like Claude to directly query and validate against your spec.

### Installation

```bash
go install github.com/plexusone/systemspec-designsystem/cmd/dss-mcp@latest
```

### Claude Desktop Integration

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json`):

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

### Available Tools (37 total)

| Category | Tools | Purpose |
|----------|-------|---------|
| Spec Reading | `get_component`, `list_components`, `get_token`, `list_tokens`, `get_pattern`, `list_patterns`, `get_meta` | Query the design system |
| Guidance | `generate_prompt`, `get_variants`, `get_props`, `get_anti_patterns` | Implementation guidance |
| Validation | `validate_file`, `validate_directory`, `check_colors`, `check_spacing` | Code compliance checking |
| Fix | `fix_file`, `suggest_fixes`, `fix_colors`, `fix_spacing`, `fix_accessibility`, `fix_directory` | Auto-fix violations |
| Lint | `lint_spec`, `list_lint_rules`, `check_agent_readiness` | Spec completeness checking |
| Validators | `list_validators`, `get_validator`, `get_validation_requirements`, `get_validator_invocation` | External validator delegation |
| Compliance | `generate_compliance_report`, `check_release_gate`, `run_fix_loop`, `fix_and_verify`, `get_compliance_certificate` | Release workflow |
| Accessibility | `get_accessibility_requirements`, `get_a11y_anti_patterns`, `suggest_contrast_token`, `get_component_fix_context` | Accessibility integration for agent-a11y |

### Example Workflows

**Component Implementation:**
```
User: "Implement a button component following the design system"

Claude: [calls get_component with id="button"]
        [calls get_variants with component_id="button"]
        [implements component]
        [calls validate_file to check compliance]
```

**Release Workflow:**
```
User: "Prepare the codebase for release"

Claude: [calls lint_spec to check spec completeness]
        [calls run_fix_loop to auto-fix violations]
        [calls check_release_gate for go/no-go decision]
        "Release approved with score 92/100. Certificate: a3f2c8b1"
```

See [MCP Server Documentation](docs/mcp-server.md) for full details.

### Generate Options

```bash
dss generate --llm ./DESIGN_CONTEXT.md  # LLM context (primary use case)
dss generate --css ./src/index.css      # Tailwind v4 @theme block
dss generate --types ./src/lib/types.ts # TypeScript interfaces
dss generate --package ./dist           # NPM package with framework presets
```

### Generating NPM Packages

Generate publishable NPM packages with framework-specific outputs:

```bash
# Generate package with all targets
dss generate --package ./dist --targets all

# Generate with specific targets
dss generate --package ./dist --targets css,tailwind,shadcn

# Custom scope and name
dss generate --package ./dist --scope @myorg --name tokens
```

**Available Targets:**

| Target | Output | Use Case |
|--------|--------|----------|
| `css` | CSS custom properties | Vanilla CSS/HTML projects |
| `tailwind` | Tailwind v4 preset and theme | Tailwind CSS v4+ projects |
| `shadcn` | ShadCN/UI theme variables | ShadCN/UI component library |
| `mkdocs-material` | MkDocs Material theme | Documentation sites |
| `scss` | SCSS variables | Sass/SCSS projects |
| `json` | Raw JSON tokens | Build tool pipelines |
| `w3c` | W3C Design Tokens format | Standards-compliant token exchange |

**Using Generated Packages:**

Tailwind v4:

```javascript
import tokens from '@myorg/design-tokens/tailwind';

export default {
  presets: [tokens],
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
};
```

ShadCN:

```css
@import '@myorg/design-tokens/shadcn';
```

MkDocs Material:

```yaml
extra_css:
  - https://unpkg.com/@myorg/design-tokens/mkdocs/theme.css
```

### Theme Bindings

```bash
# Generate CSS bindings from themeBindings configuration
dss bind --output ./theme.css

# Generate TypeScript constants
dss bind --format typescript --output ./theme.ts

# Use semantic auto-mapping strategy
dss bind --strategy semantic --output ./theme.css
```

### Contract Commands

```bash
# Show a component's theming contract
dss contract show button

# Validate all theming contracts
dss contract validate
```

### Validate Output

```
✓ Passed:
  Button: validated against spec

⚠ Warnings (2):
  [no-hardcoded-colors] ./src/components/Card.tsx:45
    Hardcoded color '#333' - use CSS variable from design system
  [button-accessible-name] ./src/components/IconButton.tsx:12
    Icon-only button should have aria-label

Summary: 3 checks (1 passed, 2 warnings, 0 errors)
```

## LLM Context Structure

The generated `DESIGN_CONTEXT.md` includes:

```markdown
# My Design System

## Design Principles
- Clarity Over Complexity: ...

## Design Tokens
| Token | Value | Usage |
|-------|-------|-------|
| `primary` | `hsl(222 47% 11%)` | Primary actions and CTAs |

## Components

### Button
**Intent:** Trigger user actions
**Use in:** forms, dialogs, toolbars
**Don't use in:** inline-text, navigation

**Anti-patterns:**
- Multiple primary buttons per view
- Using button for navigation (use Link)

**Examples:**
<Button>Save</Button>
<Button variant="destructive">Delete</Button>
```

This structure is optimized for LLM comprehension and consistent code generation.

## Theming Contracts

DSS supports formal theming contracts between component libraries and consuming applications:

**Component Side (themingContract):**
```json
{
  "themingContract": {
    "prefix": "--btn",
    "tokens": [
      {
        "id": "background",
        "cssProperty": "--btn-background",
        "semantic": "primary",
        "defaultLight": "#0066CC",
        "defaultDark": "#3399FF"
      }
    ]
  }
}
```

**Application Side (themeBindings):**
```json
{
  "themeBindings": [
    {
      "component": "button",
      "mappings": [
        { "from": "brand-primary", "to": "background" }
      ]
    }
  ]
}
```

**Generate CSS:**
```bash
dss bind --output ./theme.css
```

See [Theming Specification](docs/specification/theming.md) for details.

## Figma Integration

> **Note:** Figma integration is for teams transitioning from traditional design workflows. For AI-native development, DSS specs are authored directly—Figma is not required.

For teams that still use Figma:

- **Tokens Studio** can sync Figma variables ↔ JSON tokens
- Future: `dss import figma` / `dss export figma` commands

If your workflow is AI-first (specs authored in JSON, UI generated by AI), skip Figma entirely.

## Canonical Layers

DSS defines 9 layers (most projects only need 3):

| Layer | Required | Purpose |
|-------|----------|---------|
| **Meta** | Yes | Name, version |
| **Foundations** | Yes | Tokens (colors, typography, spacing) |
| **Components** | Yes | UI elements with LLM context |
| Principles | No | Design philosophy |
| Patterns | No | Multi-component solutions |
| Templates | No | Page layouts |
| Content | No | Voice & tone |
| Accessibility | No | WCAG requirements |
| Governance | No | Versioning policies |

## Go SDK

```go
import dss "github.com/plexusone/systemspec-designsystem/sdk/go"

ds, _ := dss.LoadDesignSystem("./my-design-system/")
ds.Validate()

// Generate for AI context
prompt, _ := ds.GenerateLLMPrompt(dss.DefaultLLMPromptOptions())

// Generate for build
css, _ := ds.GenerateCSS(dss.DefaultCSSOptions())
types, _ := ds.GenerateReactTypes(dss.DefaultReactOptions())
```

### Loading from Embedded Filesystems

Bundle design system specs into binaries using Go's `embed` package:

```go
import (
    "embed"
    "io/fs"
    dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

//go:embed spec/*
var specFS embed.FS

func main() {
    sub, _ := fs.Sub(specFS, "spec")
    ds, _ := dss.LoadDesignSystemFromFS(sub)
    // Use ds...
}
```

This enables single-binary distribution for MCP servers and CLI tools.

## Project Structure

```
systemspec-designsystem/
├── cmd/
│   ├── dss/              # CLI tool
│   │   └── cmd/
│   │       ├── generate.go    # Generate CSS, TypeScript, LLM
│   │       ├── validate.go    # Validate implementations
│   │       ├── lint_spec.go   # Spec completeness linting
│   │       ├── bind.go        # Theme bindings generation
│   │       ├── contract.go    # Theming contract commands
│   │       └── info.go
│   └── dss-mcp/          # MCP server
│       └── main.go       # MCP server entry point
├── sdk/go/               # Go SDK
│   ├── loader.go         # Load design systems
│   ├── service.go        # Service layer for CLI/MCP
│   ├── validate_file.go  # File validation logic
│   ├── fix_file.go       # Auto-fix violations
│   ├── lint_spec.go      # Spec completeness linting
│   ├── compliance.go     # Compliance reporting
│   ├── fix_loop.go       # Fix-validate loop
│   ├── validators.go     # External validator types
│   ├── theming.go        # Theming contract types
│   ├── gen_css.go        # CSS generator
│   ├── gen_react.go      # TypeScript generator
│   ├── gen_llm.go        # LLM context generator
│   └── ...               # Other generators
├── skills/designsystem/  # MCP skill definition (37 tools)
│   ├── skill.go          # Skill interface
│   ├── tools_spec.go     # Spec reading tools
│   ├── tools_guidance.go # Guidance tools
│   ├── tools_validate.go # Validation tools
│   ├── tools_fix.go      # Fix tools
│   ├── tools_lint.go     # Lint tools
│   ├── tools_validators.go # Validator delegation tools
│   ├── tools_compliance.go # Compliance & release tools
│   └── tools_a11y.go     # Accessibility integration tools
├── schema/               # JSON Schemas (generated)
├── ui/                   # Web component viewer (Lit)
└── docs/                 # MkDocs documentation
```

## Roadmap

- [x] Core SDK (9 canonical layers)
- [x] CLI (`generate`, `validate`, `info`)
- [x] Code generators (CSS, TypeScript, LLM)
- [x] Documentation (MkDocs)
- [x] Theming contracts and bindings
- [x] Diagram generators (Mermaid, D2)
- [x] W3C Design Tokens export
- [x] Web component viewer (ui/)
- [x] MCP server for AI assistants (37 tools)
- [x] `dss lint-spec` for spec completeness
- [x] External validator delegation (WCAG, API style)
- [x] Compliance reporting and release gates
- [x] Fix-validate loop with convergence detection
- [x] Visual regression testing with w3pilot integration
- [x] Evaluation system with rubric-based scoring
- [x] HTML documentation generator with Material Web demos
- [ ] `dss init` scaffolding
- [ ] CI/CD GitHub Actions integration
- [ ] Color contrast validation
- [ ] Figma tokens import/export (for transitioning teams)

## License

MIT
