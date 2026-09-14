# SystemSpec: Design System - Implementation Tasks

## Project Goals

DSS exists to solve two problems in AI-native development:

1. **AI agents build compliant UI** - Generated code follows the design system
2. **AI agents review code for compliance** - Automated validation catches violations

## Phase 1: Core SDK ✅

**Status**: Complete

Go SDK with types for all 9 canonical layers.

- [x] `sdk/go/` nested module
- [x] Meta, Principles, Foundations types
- [x] Components with LLM context
- [x] Patterns, Templates, Content types
- [x] Accessibility, Governance types
- [x] JSON/YAML loader
- [x] JSON Schema generation

## Phase 2: CLI Tool ✅

**Status**: Complete

The `dss` CLI using Cobra.

- [x] `dss info` - Display design system metadata
- [x] `dss generate` - Generate CSS, TypeScript, LLM prompt
- [x] `dss validate` - Validate implementations against spec

### Code Generators

- [x] `gen_css.go` - Tailwind v4, CSS vars, SCSS, MkDocs Material
- [x] `gen_react.go` - TypeScript interfaces
- [x] `gen_llm.go` - LLM-optimized context (Markdown)
- [x] `gen_package.go` - NPM package generation with framework presets

### NPM Package Generation

Generate publishable NPM packages with framework-specific outputs:

- [x] `--package` flag for `dss generate`
- [x] Targets: css, tailwind, shadcn, mkdocs-material, scss, json, w3c
- [x] Auto-generated package.json, README, TypeScript types
- [x] Tailwind v4 preset generation
- [x] ShadCN theme CSS generation

## Phase 3: Documentation ✅

**Status**: Complete

- [x] MkDocs site (`docs/`)
- [x] CLI reference
- [x] Specification reference
- [x] Getting started guide
- [x] README with goals and assessment

## Phase 4: Enhanced LLM Context

**Status**: Not Started

Improve AI code generation quality.

- [ ] XML output format for Claude (better structured parsing)
- [ ] Component relationship graph (what goes with what)
- [ ] Example compositions (not just single components)
- [ ] Negative examples (what bad code looks like)
- [ ] Token semantic groupings (surface colors vs. feedback colors)

## Phase 5: Enhanced Validation

**Status**: Not Started

Improve AI code review capability.

- [ ] Color contrast validation (WCAG AA/AAA ratios)
- [ ] Spacing consistency checks
- [ ] Component composition validation (valid parent-child)
- [ ] Custom validation rules (user-defined patterns)
- [ ] Severity levels (error, warning, info)
- [ ] Auto-fix suggestions in output

## Phase 6: CLI Enhancements

**Status**: Not Started

- [ ] `dss init` - Scaffold new design system spec
- [ ] `dss lint` - Lint spec for completeness
- [ ] `dss diff` - Compare spec versions (breaking changes)

## Phase 7: Integrations

**Status**: Not Started

### CI/CD (Priority)

- [ ] GitHub Action for `dss validate`
- [ ] JSON Schema validation in editors
- [ ] Pre-commit hook configuration

### Figma (For Transitioning Teams)

> Note: Figma integration is for teams migrating from traditional workflows. AI-native teams author specs directly.

- [ ] `dss import figma` - Import from Tokens Studio
- [ ] `dss export figma` - Export to Figma format
- [ ] Document migration workflow

## Goal Assessment

### Goal 1: AI Builds Compliant UI

| Aspect | Status | Effectiveness |
|--------|--------|---------------|
| Token context | ✅ Done | High - values + semantic meanings |
| Component intent | ✅ Done | High - `intent`, `allowedContexts` |
| Anti-patterns | ✅ Done | High - explicit "don't do this" |
| Examples | ✅ Done | Medium - single component examples |
| Compositions | ❌ Todo | Would improve multi-component generation |
| Negative examples | ❌ Todo | Would reduce common mistakes |

**Current effectiveness: High** - Significantly improves AI output quality

### Goal 2: AI Reviews for Compliance

| Check Type | Status | Effectiveness |
|------------|--------|---------------|
| Hardcoded colors | ✅ Done | High |
| Invalid variants | ✅ Done | High |
| Accessibility attrs | ✅ Done | High |
| Anti-patterns | ✅ Done | Medium - pattern matching |
| Color contrast | ❌ Todo | Would catch WCAG issues |
| Composition rules | ❌ Todo | Would catch structural issues |
| Visual correctness | N/A | Out of scope (use Chromatic) |

**Current effectiveness: Medium-High** - Catches token/accessibility issues, not visual

## Known Limitations

DSS intentionally does not:

1. **Implement components** - DSS is metadata, not code
2. **Validate visuals** - Use Chromatic/Percy for visual regression
3. **Replace Figma for designers** - Designers use design tools, not JSON
4. **Run at runtime** - DSS is build-time/CI only

These are handled by complementary tools.

## Future Ideas (Not Committed)

- Vue/Svelte type generators
- Dark mode CSS generation
- Static documentation site generator
- VS Code extension (schema validation, autocomplete)
- AI-powered spec generation from existing CSS
- Storybook story generation from spec
