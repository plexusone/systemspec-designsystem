# Minimal Design System Example

A minimal example demonstrating the systemspec-designsystem format.

## Structure

```
minimal-system/
├── meta.json                    # System metadata
├── foundations/
│   ├── colors.json              # Color tokens
│   ├── typography.json          # Typography tokens
│   └── spacing.json             # Spacing scale
└── components/
    ├── button.json              # Button component
    └── input.json               # Input component
```

## Loading with Go SDK

```go
package main

import (
    "fmt"
    dss "github.com/plexusone/systemspec-designsystem/sdk/go"
)

func main() {
    ds, err := dss.LoadDesignSystem("examples/minimal-system")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Loaded: %s v%s\n", ds.Meta.Name, ds.Meta.Version)
    fmt.Printf("Colors: %d\n", len(ds.Foundations.Colors))
    fmt.Printf("Components: %d\n", len(ds.Components))
}
```

## Validating Against Schema

```bash
# Using ajv-cli or similar JSON Schema validator
ajv validate -s schema/design-system.schema.json -d examples/minimal-system/meta.json
```

## Contents

### Foundations

- **Colors**: Primary, secondary, success, danger, and neutral palette
- **Typography**: Sans and mono font families, size scale, type styles
- **Spacing**: 4px base unit with standard scale (0-64px)

### Components

- **Button**: 5 variants, 6 states, full LLM context
- **Input**: Text input with validation states
