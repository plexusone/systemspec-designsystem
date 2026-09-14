# Design Systems Catalog

Design systems that can be implemented as DSS specifications. Community contributions welcome.

## Tier 1: Enterprise (High Priority)

| Design System | Organization | Status | Repo |
|---------------|--------------|--------|------|
| **Material Design 3** | Google | In Progress | `dss-material/` |
| **Carbon** | IBM | Planned | `dss-carbon/` |
| **Fluent 2** | Microsoft | Planned | `dss-fluent/` |
| **Spectrum** | Adobe | Planned | `dss-spectrum/` |

## Tier 2: Popular Ecosystems

| Design System | Organization | Status | Notes |
|---------------|--------------|--------|-------|
| **Atlassian Design** | Atlassian | Planned | B2B/SaaS standard |
| **Lightning** | Salesforce | Planned | Enterprise CRM |
| **Polaris** | Shopify | Planned | E-commerce |
| **Ant Design** | Ant Group | Planned | Popular in Asia |

## Tier 3: Modern/Community

| Design System | Organization | Status | Notes |
|---------------|--------------|--------|-------|
| **Radix/Shadcn** | Community | Planned | Headless primitives, Tailwind-native |
| **Chakra UI** | Community | Planned | Accessibility-first |
| **Primer** | GitHub | Planned | Developer-focused |
| **Mantine** | Community | Planned | React ecosystem |

## Platform Design Systems

| Platform | Guidelines | Notes |
|----------|------------|-------|
| **Apple HIG** | Human Interface Guidelines | iOS, macOS, watchOS |
| **Android** | Material (native) | Android platform |
| **Windows** | Fluent | Windows apps |
| **Web** | WCAG, WAI-ARIA | Accessibility standards |

## Implementation Priority

1. **Material Design 3** - Most widely adopted, archived official repos
2. **Carbon** - Enterprise-focused, strong accessibility
3. **Fluent 2** - Cross-platform (Web, Windows, mobile)
4. **Spectrum** - Enterprise, accessibility-first

## Strategic Context

Official Material implementations are in maintenance mode or archived:

- `material-components/material-web` - maintenance mode, pending new maintainers
- `material-components/material-components-android` - maintenance mode, migrating to Compose
- `material-components-flutter/` - archived Nov 30, 2023
- `material-components-ios/` - archived Dec 10, 2025

DSS provides a platform-agnostic, machine-readable source of truth that can:

1. Drive AI-assisted implementation across all platforms
2. Enable community-maintained Flutter/iOS specs (since official repos archived)
3. Provide consistent evaluation and quality metrics
4. Extract specs from archived repos, maintain going forward

## Contributing

To create a DSS implementation for a design system:

1. Fork this pattern: `plexusone/dss-material`
2. Extract specs from official documentation
3. Run `dss eval` to measure completeness
4. Submit PR to add to this catalog

### Creating a New DSS Spec Repository

```bash
# Clone the template
git clone https://github.com/plexusone/dss-material dss-<name>
cd dss-<name>

# Update module path
go mod edit -module github.com/<org>/dss-<name>

# Create spec structure
mkdir -p specs/v1/{foundations,components}

# Add meta.json
cat > specs/v1/meta.json << 'EOF'
{
  "name": "<Design System Name>",
  "version": "1.0.0",
  "description": "DSS specification for <Design System>"
}
EOF

# Evaluate completeness
dss eval --spec ./specs/v1 --json

# Generate documentation
dss render --spec ./specs/v1 --output ./docs
```

## Status Legend

- Complete - Fully specified with evaluation score > 90%
- In Progress - Active development
- Planned - On roadmap, not started
- Needs Maintainer - Looking for community contribution

## Platform Conformance Model

Design systems for apps often need to conform with platform design systems:

```
                    App Design System
                  (e.g., Material Design)
                           |
                           | must conform to
                           v
                 Platform Design System
    +-------------+-------------+-------------+--------+
    | Apple HIG   | Android     | Salesforce  | Web    |
    | (UIKit/     | (native     | Lightning   | (WCAG, |
    |  SwiftUI)   |  Material)  | (AppExch)   |  W3C)  |
    +-------------+-------------+-------------+--------+
```

### Conformance Examples

| App DS | Platform | Conformance Requirements |
|--------|----------|--------------------------|
| Material iOS | Apple iOS | HIG navigation patterns, UIKit/SwiftUI conventions |
| Material Android | Android | Native Material (direct), platform UI conventions |
| Material Web | Salesforce | Lightning Design System container constraints |
| Material Web | Enterprise | WCAG AA, corporate brand guidelines |

### DSS Platform Conformance Fields

```yaml
platformConformance:
  - platform: ios
    guidelines: "Apple Human Interface Guidelines"
    adaptations:
      - "Use native navigation bar on iOS"
      - "Respect safe areas"
      - "Support Dynamic Type"
  - platform: salesforce
    guidelines: "Lightning Design System"
    adaptations:
      - "Use SLDS spacing tokens"
      - "Respect Salesforce container widths"
```

## Related Resources

- [systemspec-designsystem Documentation](./index.md)
- [Getting Started Guide](./getting-started.md)
- [SDK Reference](./sdk.md)
- [CLI Reference](./cli.md)
