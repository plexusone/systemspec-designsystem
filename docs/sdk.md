# Go SDK

The Go SDK provides programmatic access to DSS for loading, validating, and generating code from design system specifications.

## Installation

```bash
go get github.com/plexusone/systemspec-designsystem/sdk/go
```

## Quick Start

```go
package main

import (
    "fmt"
    dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func main() {
    // Load from directory
    ds, err := dss.LoadDesignSystem("./my-design-system/")
    if err != nil {
        panic(err)
    }

    // Validate
    if err := ds.Validate(); err != nil {
        panic(err)
    }

    fmt.Printf("Loaded: %s v%s\n", ds.Meta.Name, ds.Meta.Version)
}
```

## Loading Design Systems

### From Directory

```go
ds, err := dss.LoadDesignSystem("./my-design-system/")
```

The loader scans for:

- `meta.json` (required)
- `principles.json`, `accessibility.json`, etc.
- `foundations/`, `components/`, `patterns/` subdirectories

### From Single File

```go
ds, err := dss.LoadDesignSystem("./design-system.json")
```

### From Embedded Filesystem

Load design systems from Go's `embed.FS` or any `fs.FS` interface. This enables bundling specs into binaries for distribution.

```go
import (
    "embed"
    "io/fs"
    dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

//go:embed spec/*
var specFS embed.FS

func main() {
    // Get the subdirectory within the embedded FS
    sub, err := fs.Sub(specFS, "spec")
    if err != nil {
        panic(err)
    }

    // Load from embedded filesystem
    ds, err := dss.LoadDesignSystemFromFS(sub)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Loaded: %s v%s\n", ds.Meta.Name, ds.Meta.Version)
}
```

This is useful for:

- **MCP servers** with embedded specs (single binary distribution)
- **CLI tools** with bundled design systems
- **Testing** with in-memory filesystems
- **Lambda/serverless** deployments

## Service Layer

The SDK provides a `Service` type that wraps a `DesignSystem` and exposes operations suitable for CLI tools, MCP servers, and SDK consumers.

### Creating a Service

```go
ds, _ := dss.LoadDesignSystem("./my-design-system/")
service := dss.NewService(ds)
```

### Service Methods

```go
// Component operations
comp, err := service.GetComponent(ctx, "button")
components := service.ListComponents(ctx)

// Token operations
token, err := service.GetToken(ctx, "color", "primary")
tokens := service.ListTokens(ctx, "color")

// Pattern operations
pattern, err := service.GetPattern(ctx, "form-validation")
patterns := service.ListPatterns(ctx)

// Metadata
meta := service.GetMeta(ctx)

// Validation
result, err := service.ValidateFile(ctx, "./src/Button.tsx", nil)
result, err := service.ValidateDirectory(ctx, "./src/components", nil)

// Auto-fixing
fixResult, err := service.FixFile(ctx, "./src/Button.tsx", nil)
fixResult, err := service.SuggestFixes(ctx, "./src/Button.tsx", nil) // dry run
fixResults, err := service.FixDirectory(ctx, "./src/components", nil)

// Generation
prompt, err := service.GenerateLLMPrompt(ctx, nil)
```

### Validation Results

```go
type ValidationResult struct {
    Files      int
    Violations []Violation
    Summary    ValidationSummary
}

type Violation struct {
    File     string
    Line     int
    Column   int
    Rule     string
    Message  string
    Severity string // "error", "warning", "info"
    Context  string
}
```

### Fix Results

```go
type FixResult struct {
    Path          string
    Fixes         []Fix
    ModifiedLines int
    NewContent    string    // Only populated if DryRun is false
    Summary       FixSummary
}

type Fix struct {
    Line        int
    Column      int
    Rule        string
    Original    string
    Replacement string
    TokenUsed   string  // Token ID used for replacement
}

type FixSummary struct {
    ColorsFixed        int
    SpacingFixed       int
    AccessibilityFixed int
}
```

### Auto-Fixing Violations

The Service provides methods to automatically fix design system violations:

```go
// Fix violations in a single file (writes changes)
result, err := service.FixFile(ctx, "./src/Button.tsx", &dss.FixOptions{
    Rules:  []string{"no-hardcoded-colors", "use-spacing-scale"},
    DryRun: false,
})

// Preview fixes without applying (dry run)
result, err := service.SuggestFixes(ctx, "./src/Button.tsx", nil)
for _, fix := range result.Fixes {
    fmt.Printf("Line %d: %s → %s (token: %s)\n",
        fix.Line, fix.Original, fix.Replacement, fix.TokenUsed)
}

// Fix all files in a directory
results, err := service.FixDirectory(ctx, "./src/components", nil)
for _, r := range results {
    fmt.Printf("%s: %d fixes applied\n", r.Path, len(r.Fixes))
}
```

**Fix Options:**

```go
type FixOptions struct {
    // Rules to apply (default: all)
    Rules []string

    // DryRun previews fixes without writing
    DryRun bool
}
```

**Available Fix Rules:**

| Rule | Description |
|------|-------------|
| `no-hardcoded-colors` | Replace hex/rgb colors with closest design token |
| `use-spacing-scale` | Replace pixel values with spacing scale tokens |
| `img-alt-required` | Add alt text to images (derived from filename) |
| `button-accessible-name` | Add aria-label to icon-only buttons |

**Color Matching:**

The fixer uses intelligent color matching:

1. Exact match - if hex value matches a token exactly
2. Closest match - finds the perceptually closest token using color distance

```go
// Example: #3b82f6 matches "primary-500" exactly
// Example: #3a7fef maps to closest token "primary-500"
```

**Spacing Matching:**

Spacing values are snapped to the nearest token in the spacing scale:

```go
// Example: "12px" → "spacing-3" (if scale is 4px base)
// Example: "17px" → "spacing-4" (snaps to 16px)
```

## Code Generation

### Generate CSS

```go
opts := dss.DefaultCSSOptions()
opts.Format = "tailwind4"  // or "css-vars", "scss"
opts.IncludeComments = true

css, err := ds.GenerateCSS(opts)
if err != nil {
    panic(err)
}
fmt.Println(css)
```

### CSS Formats

| Format | Output |
|--------|--------|
| `tailwind4` | Tailwind v4 `@theme` block |
| `css-vars` | Standard `:root` custom properties |
| `scss` | SCSS variables |

### Generate TypeScript Types

```go
opts := dss.DefaultReactOptions()

types, err := ds.GenerateReactTypes(opts)
if err != nil {
    panic(err)
}
fmt.Println(types)
```

Output example:

```typescript
export type ButtonVariant = 'default' | 'secondary' | 'destructive';

export interface ButtonProps {
  variant?: ButtonVariant;
  disabled?: boolean;
  children: React.ReactNode;
}
```

### Generate LLM Prompt

```go
opts := dss.DefaultLLMPromptOptions()
opts.Format = "markdown"  // or "xml"

prompt, err := ds.GenerateLLMPrompt(opts)
if err != nil {
    panic(err)
}
fmt.Println(prompt)
```

## Type Reference

### DesignSystem

```go
type DesignSystem struct {
    Meta          Meta
    Principles    []Principle
    Foundations   Foundations
    Components    []Component
    Patterns      []Pattern
    Templates     []Template
    Content       *Content
    Accessibility *Accessibility
    Governance    *Governance
    ThemeBindings []ThemeBindings  // Application token mappings
}
```

### Component

```go
type Component struct {
    ID              string
    Name            string
    Description     string
    Category        string
    Variants        []Variant
    States          []State
    Props           []Prop
    Events          []ComponentEvent   // Custom events
    Slots           []Slot
    Uses            []string           // Component dependencies
    Constraints     *Constraints
    ThemingContract *ThemingContract   // External theming API
    Accessibility   *ComponentA11y
    LLM             *LLMContext
}
```

### ThemingContract

```go
type ThemingContract struct {
    Prefix      string       // CSS custom property prefix (e.g., "--btn")
    Description string
    Tokens      []ThemeToken
}

type ThemeToken struct {
    ID           string  // Token identifier
    CSSProperty  string  // Full CSS property name
    Semantic     string  // Semantic category for auto-mapping
    Description  string
    DefaultLight string  // Light mode default
    DefaultDark  string  // Dark mode default
}
```

### ThemeBindings

```go
type ThemeBindings struct {
    Component string          // Component ID to theme
    SpecURL   string          // Optional URL to fetch component spec
    ThemeMode string          // "light", "dark", or empty for both
    Strategy  string          // "explicit", "semantic", "inherit"
    Mappings  []TokenMapping
}

type TokenMapping struct {
    From      string  // Application token ID
    To        string  // Component token ID
    Transform string  // Optional CSS transform function
}
```

### LLMContext

```go
type LLMContext struct {
    Intent            string
    AllowedContexts   []string
    ForbiddenContexts []string
    ExampleUsage      []string
    AntiPatterns      []string
    SemanticMeaning   string
    RelatedElements   []string
    PriorityScore     int
}
```

## Validation

```go
if err := ds.Validate(); err != nil {
    fmt.Printf("Validation failed: %v\n", err)
}
```

Validation checks:

- Required fields present (name, version)
- Component IDs are unique
- Variant references are valid
- Color values are parseable

## Theming Contracts

### Validate Contracts

```go
// Validate a single component's contract
report := dss.ValidateContract(&component)
if !report.Passed {
    for _, err := range report.Errors {
        fmt.Printf("Error: %s: %s\n", err.Field, err.Message)
    }
}

// Validate all contracts in a design system
reports := dss.ValidateAllContracts(ds)
fmt.Printf("Errors: %d, Warnings: %d\n",
    dss.TotalErrors(reports),
    dss.TotalWarnings(reports))
```

### Generate Theme Bindings

```go
opts := dss.BindingOptions{
    Format:          dss.FormatCSS,  // or FormatTypeScript, FormatSCSS
    DefaultStrategy: "semantic",      // or "explicit", "inherit"
}

bindings, err := dss.GenerateBindings(ds, opts)
if err != nil {
    panic(err)
}

// Write to file
f, _ := os.Create("theme.css")
dss.WriteBindings(f, bindings)
```

## Diagram Generation

### Mermaid Diagrams

```go
opts := dss.DefaultMermaidOptions()
opts.Direction = "LR"  // Left to Right
opts.ShowSlots = true

// Component relationship diagram
mermaid, err := ds.GenerateMermaid(opts)

// Single component focus diagram
diagram, err := ds.GenerateMermaidComponentDiagram("button", opts)

// Token usage diagram
tokenDiagram, err := ds.GenerateMermaidTokenDiagram(opts)
```

### D2 Diagrams

```go
opts := dss.DefaultD2Options()
opts.Sketch = true  // Hand-drawn look

// Full architecture diagram
d2, err := ds.GenerateD2(opts)

// Single component diagram
diagram, err := ds.GenerateD2ComponentDiagram("button", opts)
```

## W3C Design Tokens Export

```go
opts := dss.DefaultW3CTokensOptions()
opts.IncludeDescriptions = true
opts.IncludeExtensions = true

tokens, err := ds.GenerateW3CTokens(opts)
// Output: W3C Design Tokens Community Group format JSON
```

## Documentation Generation

```go
opts := dss.DefaultDocsOptions()
opts.OutputDir = "docs/spec"
opts.IncludeMermaid = true
opts.IncludeD2 = true
opts.MkDocsCompatible = true

output, err := ds.GenerateDocs(opts)
// output.Files contains map of path -> content
for path, content := range output.Files {
    // Write each file
}
```

## NPM Package Generation

Generate publishable NPM packages with framework-specific presets and configurations.

### GeneratePackage Method

```go
opts := dss.PackageGeneratorOptions{
    OutputDir:   "./dist",
    Scope:       "@myorg",
    PackageName: "design-tokens",
    Targets: []dss.PackageTarget{
        dss.TargetCSS,
        dss.TargetTailwind,
        dss.TargetShadCN,
    },
    IncludeReadme: true,
}

err := ds.GeneratePackage(opts)
if err != nil {
    panic(err)
}
```

### PackageGeneratorOptions

```go
type PackageGeneratorOptions struct {
    // OutputDir is the directory to write the package to
    OutputDir string

    // Targets specifies which framework outputs to generate
    Targets []PackageTarget

    // Scope is the NPM scope (e.g., "@plexusone")
    // If empty, derived from meta
    Scope string

    // PackageName is the package name (default: "design-tokens")
    PackageName string

    // DryRun previews output without writing files
    DryRun bool

    // IncludeReadme generates a README.md file
    IncludeReadme bool
}
```

### PackageTarget Constants

```go
const (
    TargetCSS            PackageTarget = "css"
    TargetTailwind       PackageTarget = "tailwind"
    TargetShadCN         PackageTarget = "shadcn"
    TargetMkDocsMaterial PackageTarget = "mkdocs-material"
    TargetSCSS           PackageTarget = "scss"
    TargetJSON           PackageTarget = "json"
    TargetW3C            PackageTarget = "w3c"
)
```

### Example: Generate All Targets

```go
ds, _ := dss.LoadDesignSystem("./my-design-system/")

opts := dss.DefaultPackageOptions()
opts.OutputDir = "./dist"
opts.Targets = []dss.PackageTarget{
    dss.TargetCSS,
    dss.TargetTailwind,
    dss.TargetShadCN,
    dss.TargetMkDocsMaterial,
    dss.TargetSCSS,
    dss.TargetJSON,
    dss.TargetW3C,
}

err := ds.GeneratePackage(opts)
```

### Example: Preview Without Writing

```go
opts := dss.DefaultPackageOptions()
opts.OutputDir = "./dist"
opts.DryRun = true

err := ds.GeneratePackage(opts)
// No files written, validates package structure
```

### Generated Package Structure

```
dist/
├── package.json           # NPM manifest with exports
├── index.js               # CommonJS entry point
├── index.mjs              # ES module entry point
├── index.d.ts             # TypeScript declarations
├── README.md              # Usage documentation
├── css/
│   └── tokens.css
├── tailwind/
│   ├── preset.js
│   └── theme.json
├── shadcn/
│   └── theme.css
├── mkdocs/
│   └── theme.css
├── scss/
│   └── tokens.scss
├── json/
│   └── tokens.json
└── w3c/
    └── tokens.json
```

## JSON Schema Generation

```go
schema := dss.GenerateSchema()
data, _ := json.MarshalIndent(schema, "", "  ")
fmt.Println(string(data))
```

## Building a Custom Tool

```go
package main

import (
    "flag"
    "fmt"
    "os"

    dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func main() {
    dir := flag.String("dir", ".", "Design system directory")
    flag.Parse()

    ds, err := dss.LoadDesignSystem(*dir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Custom processing
    for _, c := range ds.Components {
        fmt.Printf("Component: %s (%d variants)\n", c.Name, len(c.Variants))
    }
}
```
