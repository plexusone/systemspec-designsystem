# systemspec-designsystem

A declarative, machine-readable specification for defining complete design systems as code.

## What is DSS?

systemspec-designsystem (DSS) provides a canonical framework for expressing design systems (like Material Design, Carbon, Fluent) in a structured, version-controlled, LLM-optimized format.

## Key Features

- **Declarative definitions** - Define tokens, components, patterns as data
- **Multi-format support** - JSON/YAML for tokens, Markdown for documentation
- **LLM optimization** - Explicit intent, contexts, and constraints for AI code generation
- **Code generation** - Generate CSS, TypeScript types, and LLM prompts from spec
- **MCP Server** - Expose design systems to AI assistants via Model Context Protocol
- **Theming contracts** - Formal theming API between component libraries and applications
- **Diagram generation** - Mermaid and D2 diagrams for architecture visualization
- **Compliance validation** - Validate implementations against the spec
- **Visual regression testing** - Screenshot-based testing with w3pilot integration
- **Evaluation system** - Rubric-based spec evaluation and coverage metrics
- **HTML documentation** - Generate static docs with Material Web live demos
- **Go-first approach** - Go structs as source of truth, generating JSON Schema
- **Embedded specs** - Bundle specs into binaries with Go's embed package

## Quick Install

```bash
# Install CLI
go install github.com/plexusone/systemspec-designsystem/cmd/dss@latest

# Install Go SDK
go get github.com/plexusone/systemspec-designsystem/sdk/go
```

## Quick Example

```bash
# Show design system info
dss info

# Generate CSS, TypeScript types, and LLM context
dss generate --css ./src/index.css --types ./src/lib/types.ts --llm ./DESIGN_CONTEXT.md

# Generate theme bindings for external components
dss bind --output ./theme.css

# Validate theming contracts
dss contract validate

# Validate component implementations
dss validate ./src/components

# Evaluate spec completeness
dss eval --json > eval.json

# Generate HTML documentation with live demos
dss render --output ./docs --title "My Design System"

# Run visual regression tests
dss visual test

# Start MCP server for AI assistant integration
dss-mcp --spec ./design-system
```

## Canonical Layers

DSS defines **10 canonical layers** for complete design system specification:

| Layer | Purpose |
|-------|---------|
| **Meta** | System metadata (name, version, maintainers) |
| **Principles** | Design philosophy and guidelines |
| **Foundations** | Design tokens (colors, typography, spacing) |
| **Components** | UI elements with variants and states |
| **Patterns** | Multi-component solutions |
| **Templates** | Page layouts |
| **Content** | Voice & tone guidelines |
| **Accessibility** | WCAG compliance requirements |
| **Governance** | Versioning and deprecation policies |
| **Theming** | Token mappings to external components |

## Next Steps

- [Getting Started](getting-started.md) - Create your first design system spec
- [CLI Reference](cli.md) - Full CLI command documentation
- [MCP Server](mcp-server.md) - Integrate with AI assistants like Claude
- [Go SDK](sdk.md) - Programmatic access and embedded specs
- [Specification](specification/index.md) - Detailed spec reference
