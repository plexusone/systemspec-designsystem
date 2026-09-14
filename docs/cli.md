# CLI Reference

The `dss` command-line tool generates code artifacts and validates implementations against design system specifications.

## Installation

```bash
go install github.com/plexusone/systemspec-designsystem/cmd/dss@latest
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--dir` | `-d` | Directory containing the design system spec (default: `.`) |
| `--help` | `-h` | Help for any command |
| `--version` | `-v` | Version information |

## Commands

### dss info

Display information about a design system specification.

```bash
dss info [flags]
```

**Examples:**

```bash
# Show info for current directory
dss info

# Show info for a specific directory
dss info -d ./my-design-system
```

**Output includes:**

- Design system name and version
- Description
- Principles count
- Foundation token counts (colors, typography, spacing)
- Component list with variant counts
- Accessibility target level
- Validation status

---

### dss generate

Generate code artifacts from a design system specification.

```bash
dss generate [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--css` | Output path for CSS file |
| `--types` | Output path for TypeScript types |
| `--llm` | Output path for LLM prompt/context |
| `--css-format` | CSS format: `tailwind4` (default), `css-vars`, `scss` |

**Examples:**

```bash
# Preview all outputs to stdout
dss generate

# Generate CSS only
dss generate --css ./src/index.css

# Generate all artifacts
dss generate \
  --css ./src/index.css \
  --types ./src/lib/design-system-types.ts \
  --llm ./DESIGN_CONTEXT.md

# Use different spec directory
dss generate -d ./design-system --css ./web/src/styles.css

# Generate SCSS variables instead of Tailwind
dss generate --css-format scss --css ./src/variables.scss

# Generate standard CSS custom properties
dss generate --css-format css-vars --css ./src/vars.css

# Generate NPM package with all targets
dss generate --package ./dist --targets all

# Generate NPM package with specific targets
dss generate --package ./dist --targets css,tailwind,shadcn
```

**NPM Package Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--package` | `-p` | Output directory for NPM package |
| `--scope` | `-s` | NPM scope (e.g., `@myorg`) |
| `--name` | `-n` | Package name (default: `design-tokens`) |
| `--targets` | `-t` | Comma-separated targets (default: `css,tailwind`) |
| `--dry-run` | | Preview without writing files |

**Available Targets:**

| Target | Description |
|--------|-------------|
| `css` | CSS custom properties |
| `tailwind` | Tailwind CSS v4 preset |
| `shadcn` | ShadCN/UI theme variables |
| `mkdocs-material` | MkDocs Material theme |
| `scss` | SCSS variables |
| `json` | Raw JSON tokens |
| `w3c` | W3C Design Tokens format |
| `all` | All of the above |

**Generated Package Exports:**

The generated `package.json` includes these exports for consumers:

| Export Path | File | Usage |
|-------------|------|-------|
| `.` | `index.mjs` / `index.js` | `import { colors } from '@myorg/design-tokens'` |
| `./css` | `css/tokens.css` | `@import '@myorg/design-tokens/css'` |
| `./tailwind` | `tailwind/preset.js` | `import preset from '@myorg/design-tokens/tailwind'` |
| `./shadcn` | `shadcn/theme.css` | `@import '@myorg/design-tokens/shadcn'` |
| `./mkdocs` | `mkdocs/extra.css` | `extra_css: ['.../@myorg/design-tokens/mkdocs']` |
| `./scss` | `scss/_variables.scss` | `@import '@myorg/design-tokens/scss'` |

See [NPM Package Spec](specs/ROADMAP.md) for full `package.json` structure.

**CSS Formats:**

| Format | Description |
|--------|-------------|
| `tailwind4` | Tailwind CSS v4 `@theme` block with CSS custom properties |
| `css-vars` | Standard CSS `:root` with custom properties |
| `scss` | SCSS variables (`$color-primary`, etc.) |

---

### dss validate

Validate that component implementations comply with the design system specification.

```bash
dss validate <components-dir> [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `components-dir` | Directory containing React/TypeScript component files |

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output results as JSON (useful for CI) |

**Examples:**

```bash
# Validate components
dss validate ./src/components

# Validate with different spec directory
dss validate -d ./design-system ./web/src/components

# JSON output for CI integration
dss validate --json ./src/components
```

**Validation Checks:**

| Rule | Severity | Description |
|------|----------|-------------|
| `no-hardcoded-colors` | Warning | Hex/RGB/HSL colors should use CSS variables |
| `use-spacing-scale` | Warning | Pixel values should use spacing scale |
| `img-alt-required` | Error | Images must have alt attributes |
| `button-accessible-name` | Warning | Icon-only buttons need aria-label |
| `valid-variant` | Error | Variant values must match component spec |
| `single-primary-button` | Warning | Only one primary button per view |
| `no-nested-cards` | Warning | Cards should not be nested |

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | Validation passed (may have warnings) |
| 1 | Validation failed with errors |

---

### dss lint-spec

Lint the design system specification for completeness and best practices.

```bash
dss lint-spec [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--rules` | | Specific rules to check (comma-separated) |
| `--min-score` | | Minimum acceptable score (0-100) |
| `--json` | | Output as JSON |
| `--verbose` | `-v` | Show all issues including info level |

**Subcommands:**

```bash
# List available lint rules
dss lint-spec rules
```

**Examples:**

```bash
# Lint current directory
dss lint-spec

# Lint with minimum score requirement
dss lint-spec --min-score 80

# Lint specific rules only
dss lint-spec --rules component-has-llm-context,llm-has-intent

# JSON output for CI
dss lint-spec --json

# Show all issues including info level
dss lint-spec --verbose

# Show available rules
dss lint-spec rules
```

**Available Rules:**

| Rule | Severity | Description |
|------|----------|-------------|
| `meta-required` | Error | Design system must have name and version |
| `component-has-variants` | Warning | Components should define variants |
| `component-has-props` | Info | Components should define props |
| `component-has-llm-context` | Warning | Components should have LLM context |
| `llm-has-intent` | Error | LLM context must have intent field |
| `llm-has-anti-patterns` | Warning | LLM context should document anti-patterns |
| `llm-has-allowed-contexts` | Info | LLM context should specify allowed contexts |
| `tokens-have-descriptions` | Info | Tokens should have descriptions |
| `token-references-valid` | Error | Token references must resolve to valid tokens |
| `no-orphan-tokens` | Info | Tokens should be referenced by components |
| `component-uses-valid` | Error | Component uses references must be valid |
| `accessibility-defined` | Warning | Design system should define accessibility |
| `theming-contract-valid` | Error/Warning | Theming contracts must be properly configured |
| `validators-configured` | Warning | External validators should be configured |
| `validator-tool-required` | Error | Validators must specify a tool name |
| `validator-type-valid` | Error | Validator type must be valid |

**Output:**

```
Spec Completeness Score: 85/100

Coverage:
  Components with LLM context: 80%
  Components with variants:    90%
  Components with props:       100%
  Tokens with descriptions:    75%
  Tokens referenced:           95%

Issues: 0 errors, 2 warnings, 3 info

Warnings:
  components[2].llm [card]: Component 'card' missing LLM context for AI code generation
    → Add llm field with intent, allowedContexts, and antiPatterns

✓ Spec is agent-ready
```

**Exit Codes:**

| Code | Meaning |
|------|---------|
| 0 | Lint passed (score >= min-score, no errors) |
| 1 | Lint failed (score < min-score or has errors) |

---

### dss eval

Evaluate design system spec completeness and quality using rubric-based scoring.

```bash
dss eval [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--dir` | `-d` | Directory containing the design system spec |
| `--json` | | Output as JSON |
| `--min-score` | | Minimum acceptable score (1-5) |

**Examples:**

```bash
# Evaluate current directory
dss eval

# Output as JSON
dss eval --json > eval.json

# Require minimum score
dss eval --min-score 4
```

**Evaluation Categories (weighted):**

| Category | Weight | Description |
|----------|--------|-------------|
| Completeness | 25% | Required fields, no gaps |
| Agent-Readiness | 30% | LLM context, anti-patterns, examples |
| Accessibility | 25% | WCAG, keyboard, screen reader support |
| Documentation | 20% | Descriptions, usage guidance |

**Output:**

```
Design System Evaluation: My Design System v1.0.0

Overall Score: 4/5 (Good)

Categories:
  Completeness:    4/5 (25%)
  Agent-Readiness: 5/5 (30%)
  Accessibility:   3/5 (25%)
  Documentation:   4/5 (20%)

Coverage:
  Components: 12
  Foundations: 45 tokens
  Patterns: 3
```

---

### dss render

Generate HTML documentation for the design system with live component demos.

```bash
dss render [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--dir` | `-d` | Directory containing the design system spec |
| `--output` | `-o` | Output directory for HTML files |
| `--title` | | Custom site title |
| `--eval` | | Path to evaluation JSON file |
| `--mkdocs` | | Generate MkDocs-compatible markdown |

**Examples:**

```bash
# Generate HTML documentation
dss render --output ./docs

# With custom title
dss render --output ./docs --title "Material Design 3"

# Include evaluation data
dss render --output ./docs --eval ./evals/v3.json

# Generate MkDocs-compatible output
dss render --output ./docs/specs --mkdocs
```

**Generated Files:**

| File | Description |
|------|-------------|
| `index.html` | Dashboard landing page |
| `components.html` | Component gallery |
| `component-{id}.html` | Individual component pages with live demos |
| `tokens.html` | Token palette visualization |
| `eval.html` | Evaluation dashboard |

**Features:**

- Material Web CDN integration for live component demos
- Dark/light theme toggle
- Variant selector and disabled state controls
- PlexusOne unified navigation integration
- Self-contained CSS (no external dependencies)

---

### dss visual

Visual regression testing commands using w3pilot for browser automation.

```bash
dss visual <subcommand> [flags]
```

**Subcommands:**

#### dss visual test

Run visual regression tests against baselines.

```bash
dss visual test [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--tests` | | Filter tests by ID (comma-separated) |
| `--viewports` | | Filter by viewport (e.g., "1920x1080") |
| `--threshold` | | Diff threshold percentage (default: 0.1) |
| `--parallel` | `-p` | Number of parallel workers (default: 4) |
| `--json` | | Output results as JSON |

**Examples:**

```bash
# Run all tests
dss visual test

# Run specific tests
dss visual test --tests button-primary,card-elevated

# Specific viewport
dss visual test --viewports 375x812

# JSON output for CI
dss visual test --json
```

#### dss visual baseline generate

Generate new baselines for a version.

```bash
dss visual baseline generate <version> [flags]
```

**Examples:**

```bash
# Generate baselines for v1.0.0
dss visual baseline generate v1.0.0
```

#### dss visual baseline update

Update specific test baselines.

```bash
dss visual baseline update <version> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--tests` | Tests to update (comma-separated) |

**Examples:**

```bash
# Update specific tests
dss visual baseline update v1.0.0 --tests button-primary
```

#### dss visual baseline list

List available baseline versions.

```bash
dss visual baseline list
```

#### dss visual baseline prune

Remove old baseline versions.

```bash
dss visual baseline prune <version>
```

**Test Definition Format:**

Create `visual-tests.yaml` in your project:

```yaml
version: "1.0"
name: "Component Visual Tests"
defaults:
  threshold: 0.1
  viewports:
    - { width: 1920, height: 1080 }
    - { width: 375, height: 812 }

tests:
  - id: button-primary
    url: http://localhost:3000/components/button
    selector: ".button-primary"

  - id: card-elevated
    url: http://localhost:3000/components/card
    selector: ".card-elevated"
    stabilization:
      waitForSelector: ".card-content"
      delay: 100
```

---

### dss bind

Generate theme bindings from design system token mappings.

```bash
dss bind [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output file (default: stdout) |
| `--format` | `-f` | Output format: `css`, `typescript`, `scss` (default: `css`) |
| `--strategy` | | Mapping strategy: `explicit`, `semantic`, `inherit` (default: `explicit`) |

**Examples:**

```bash
# Generate CSS bindings to stdout
dss bind

# Generate CSS to file
dss bind --output ./theme.css

# Generate TypeScript constants
dss bind --format typescript --output ./theme.ts

# Use semantic auto-mapping
dss bind --strategy semantic --output ./theme.css

# Generate SCSS variables
dss bind --format scss --output ./theme.scss
```

**Mapping Strategies:**

| Strategy | Description |
|----------|-------------|
| `explicit` | Only use defined `from`/`to` mappings, skip unmapped tokens |
| `semantic` | Auto-map by semantic field, fall back to component defaults |
| `inherit` | Use component defaults for all unmapped tokens |

---

### dss contract

Display and validate component theming contracts.

```bash
dss contract <subcommand> [flags]
```

**Subcommands:**

#### dss contract show

Display a component's theming contract.

```bash
dss contract show <component-id> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON instead of table |

**Examples:**

```bash
# Show button's theming contract
dss contract show button

# JSON output
dss contract show button --json
```

**Output:**

```
Theming Contract: Button
Prefix: --btn

ID          CSS PROPERTY      SEMANTIC  LIGHT DEFAULT  DARK DEFAULT
--          ------------      --------  -------------  ------------
background  --btn-background  primary   #0066CC        #3399FF
text        --btn-text        text      #FFFFFF        #0A0E1A
```

#### dss contract validate

Validate all theming contracts in the design system.

```bash
dss contract validate [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON instead of table |

**Examples:**

```bash
# Validate all contracts
dss contract validate

# JSON output for CI
dss contract validate --json
```

**Validation Checks:**

| Check | Severity | Description |
|-------|----------|-------------|
| Prefix required | Error | Contract must have a prefix starting with `--` |
| Token ID required | Error | Each token must have a unique ID |
| CSS property required | Error | Each token must have a cssProperty |
| Duplicate token ID | Error | Token IDs must be unique within a contract |
| CSS property prefix | Warning | cssProperty should start with the contract prefix |
| Invalid semantic | Warning | Semantic value should be from the allowed set |
| Missing defaults | Warning | Tokens should have defaultLight and defaultDark |

**Output:**

```
✓ button: OK
⚠ card: 2 warnings
    WARN  [bg] themingContract.tokens[0].defaultLight: defaultLight not provided
✗ input: 1 error, 0 warnings
    ERROR themingContract.prefix: prefix is required

Validation passed: 3 contracts validated, 2 warnings
```

---

## Usage in Makefiles

Example `Makefile` for a project using DSS:

```makefile
DSS ?= dss
WEB_SRC = ./web/src

.PHONY: generate validate package

generate:
	@$(DSS) generate \
		--css $(WEB_SRC)/index.css \
		--types $(WEB_SRC)/lib/design-system-types.ts \
		--llm ./DESIGN_CONTEXT.md

validate:
	@$(DSS) validate $(WEB_SRC)/components

# CI validation with JSON output
validate-ci:
	@$(DSS) validate --json $(WEB_SRC)/components

# Check if generated files are current
check:
	@$(DSS) generate --css /tmp/check.css
	@diff -q $(WEB_SRC)/index.css /tmp/check.css

# Generate NPM package
package:
	@$(DSS) generate --package ./dist --targets all
```

---

## Using Generated NPM Packages

### In a Tailwind Project

```javascript
// tailwind.config.js
import preset from '@myorg/design-tokens/tailwind'

export default {
  presets: [preset],
  content: ['./src/**/*.{js,ts,jsx,tsx}'],
}
```

### In a ShadCN Project

```css
/* Import theme variables */
@import '@myorg/design-tokens/shadcn/theme.css';
```

### In MkDocs

```yaml
# mkdocs.yml
extra_css:
  - node_modules/@myorg/design-tokens/mkdocs/extra.css
```

### Direct CSS Import

```css
@import '@myorg/design-tokens/css/tokens.css';

.my-button {
  background: var(--color-primary);
  color: var(--color-text);
}
```

### JavaScript/TypeScript

```javascript
import { colors, spacing } from '@myorg/design-tokens'

console.log(colors.primary) // #...
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Design System Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Install dss
        run: go install github.com/plexusone/systemspec-designsystem/cmd/dss@latest

      - name: Validate components
        run: dss validate -d ./design-system ./src/components

      - name: Check generated files are current
        run: |
          dss generate -d ./design-system --css /tmp/index.css
          diff ./src/index.css /tmp/index.css
```

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Validate design system compliance
dss validate ./src/components

if [ $? -ne 0 ]; then
  echo "Design system validation failed!"
  exit 1
fi
```

---

## MCP Server (dss-mcp)

The `dss-mcp` command runs an MCP (Model Context Protocol) server that exposes design system operations to AI assistants like Claude.

### Installation

```bash
go install github.com/plexusone/systemspec-designsystem/cmd/dss-mcp@latest
```

### Usage

```bash
# Start MCP server with your design system spec
dss-mcp --spec ./design-system/

# Enable browser validation tools (requires w3pilot)
dss-mcp --spec ./design-system/ --browser
```

### Flags

| Flag | Description |
|------|-------------|
| `--spec`, `-s` | Path to design system spec directory (required) |
| `--browser` | Enable w3pilot browser tools for visual validation |

### Claude Desktop Configuration

Add to `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "design-system": {
      "command": "dss-mcp",
      "args": ["--spec", "/path/to/design-system"]
    }
  }
}
```

### Available Tools

The MCP server exposes 37 tools organized into eight categories:

| Category | Tools | Purpose |
|----------|-------|---------|
| Spec Reading | 7 | Query components, tokens, patterns, metadata |
| Guidance | 4 | Generate prompts, get variants, props, anti-patterns |
| Validation | 4 | Validate files/directories, check colors/spacing |
| Fix | 6 | Auto-fix colors, spacing, accessibility violations |
| Lint | 3 | Spec completeness and agent-readiness checking |
| Validators | 4 | External validator discovery and delegation |
| Compliance | 5 | Release gates, compliance reports, certificates |
| Accessibility | 4 | A11y requirements, contrast suggestions, fix context |

See [MCP Server](mcp-server.md) for full documentation.
