# CLAUDE.md - systemspec-designsystem

Project-specific instructions for Claude Code.

## Core Principle: Go-First Schema Design

**Go structs are the source of truth. JSON Schema is generated, never hand-written.**

This ensures schemas are always Go-friendly and pass `schemalint` validation.

### Workflow for Schema Changes

```bash
# 1. Edit Go structs (source of truth)
vim sdk/go/components.go

# 2. Rebuild to verify compilation
go build ./...

# 3. Regenerate JSON schemas
go run tools/generate/main.go

# 4. Validate with schemalint
schemalint lint schema/design-system.schema.json

# 5. Run tests
go test ./...

# 6. Commit both Go and generated schema together
git add sdk/go/ schema/
```

### Why Go-First?

| Benefit | Explanation |
|---------|-------------|
| Type safety | Go compiler catches errors before schema generation |
| No union ambiguity | Go structs map to concrete JSON objects |
| Consistent casing | Go struct tags control JSON property names |
| schemalint compliance | Generated schemas pass validation automatically |
| Single source of truth | No drift between Go types and JSON Schema |

### What NOT to Do

- **Never hand-edit** files in `schema/` directory
- **Never add** `anyOf`/`oneOf` without discriminators (Go doesn't produce these)
- **Never use** `interface{}` unless absolutely necessary (use concrete types)
- **Never use** `additionalProperties: true` (set to `false` in reflector)

## Project Architecture

```
systemspec-designsystem/
├── sdk/go/                    # Go SDK (source of truth)
│   ├── designsystem.go        # Root DesignSystem type
│   ├── components.go          # Component, Prop, Event types
│   ├── patterns.go            # Pattern types
│   ├── foundations.go         # Token types (colors, spacing, etc.)
│   ├── gen_css.go             # CSS generator
│   ├── gen_react.go           # TypeScript generator
│   ├── gen_llm.go             # LLM prompt generator
│   └── jsonschema.go          # Schema generation functions
│
├── schema/                    # Generated JSON Schemas (DO NOT EDIT)
│   ├── design-system.schema.json
│   ├── components/
│   └── foundations/
│
├── cmd/dss/                   # CLI tool
│   └── cmd/
│       ├── generate.go        # dss generate command
│       ├── validate.go        # dss validate command
│       └── info.go            # dss info command
│
├── tools/generate/            # Schema generation tool
│   └── main.go
│
└── examples/                  # Example design system specs
```

## DSL Design Philosophy

The schema describes a **semantic graph of UI systems**, not rendering details.

### Include (Semantic Contracts)

- Components with props, events, slots, states
- Composition relationships (`uses`, `allowedComponents`)
- Token references (`tokensUsed`)
- Accessibility requirements
- Validation constraints

### Exclude (Implementation Details)

- CSS rules or styling
- Framework-specific code (Lit, React lifecycle)
- Pixel-level layout details
- Runtime behavior logic

## Adding New Schema Features

When adding new fields or types:

1. **Design the Go struct first** in `sdk/go/`
2. **Use concrete types** - avoid `interface{}` where possible
3. **Add JSON tags** with `omitempty` for optional fields
4. **Add jsonschema tags** for validation (format, pattern, etc.)
5. **Regenerate and validate** before committing

Example:

```go
// Good - concrete type with proper tags
type ComponentEvent struct {
    Name        string             `json:"name"`
    Description string             `json:"description,omitempty"`
    Bubbles     bool               `json:"bubbles,omitempty"`
    Detail      []EventDetailField `json:"detail,omitempty"`
}

// Avoid - interface{} loses type information
type BadEvent struct {
    Data interface{} `json:"data"` // Becomes "additionalProperties: true"
}
```

## CLI Commands

```bash
# Show design system info
dss info

# Generate outputs (CSS, TypeScript, LLM prompts)
dss generate --css ./output.css --types ./types.ts --llm ./CONTEXT.md

# Validate component implementations
dss validate

# Check spec completeness
dss lint

# Report coverage metrics
dss coverage
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./sdk/go/

# Run specific test
go test -v ./sdk/go/ -run TestLoadDesignSystem
```

## Linting

```bash
# Go linting
golangci-lint run

# Schema linting (Go-friendliness)
schemalint lint schema/design-system.schema.json

# Strict mode (no composition keywords)
schemalint lint --profile scale schema/design-system.schema.json
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/invopop/jsonschema` | Generate JSON Schema from Go types |
| `github.com/spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | YAML support for spec files |

## Related Projects

- `grokify/schemalint` - JSON Schema linting for Go-friendliness
- `grokify/markdown-editor` - Markdown editor using this spec pattern
