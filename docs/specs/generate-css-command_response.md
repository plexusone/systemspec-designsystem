# Review: `dss generate css` command spec

**Reviewer:** Claude Code (systemspec-designsystem session)  
**Date:** 2026-07-22  
**Status:** Reviewed — ready to implement with clarifications below

## Overall Assessment

Strong spec. Clear motivation, concrete examples, well-scoped. The schema already defines the `output.css` contract, so this is filling a real gap — teams are manually transcribing tokens today, which causes drift.

No blockers. The items below are clarifications that will matter during implementation.

## Clarifications Needed

### 1. Mapping resolution: merge vs. replace

The spec says (line 94): if a profile defines `css.mappings`, "use those instead of the default `output.css.mappings`."

This means **replace**, not merge. If the default has 20 mappings and a profile defines 3, the profile gets only 3. That's a sharp edge — every profile must duplicate the full mapping table or accept a subset.

**Recommendation:** Document this explicitly. Consider supporting merge semantics (profile mappings override matching keys, default mappings fill the rest) as the default, with an opt-in `"mappingsMode": "replace"` if full replacement is needed.

### 2. Mapping format

The spec references `output.css.mappings` but never shows what a mapping object looks like in the JSON. Implementers need to know the structure.

**Recommendation:** Add an example, e.g.:

```json
{
  "output": {
    "css": {
      "prefix": "--brand",
      "selector": ":root",
      "mappings": {
        "colors.dark": "dark",
        "colors.navy": "navy",
        "typography.fontFamily.sans": "font-sans"
      }
    }
  }
}
```

### 3. Tailwind v4 token coverage

The `--format css` example shows colors, fonts, radius, and spacing. The `--format tailwind` example only shows colors and fonts. Are non-color tokens intentionally excluded from Tailwind output, or is the example abbreviated?

For Tailwind v4, variable naming determines which utilities get generated (`--color-*` → `bg-*`/`text-*`, `--font-*` → `font-*`, `--spacing-*` → `p-*`/`m-*`). The mapping must produce correct Tailwind namespace prefixes or the tokens exist as variables but don't surface as utility classes.

**Recommendation:** Show the full Tailwind output including spacing and radius tokens, and document which Tailwind namespaces are supported.

### 4. Error handling

The acceptance criteria say "exit 1 on invalid input" but don't specify behavior for:

- Profile name not found in the JSON
- Profile exists but has no `css` config
- `output.css` section missing entirely
- Token path in a mapping doesn't resolve to a value

**Recommendation:** Define fallback behavior for each case (error with message vs. use defaults vs. skip).

### 5. `--input` default vs. real usage

The `--input` flag defaults to `design-system.json` (CWD), but every integration example passes an explicit relative path like `../../standards/brand-design-system/design-system.json`. The default is fine, but worth noting whether the CLI should also support a `DSS_INPUT` environment variable or config file for CI use.

## Minor Observations

- **Tailwind import order** — the code example is correct (`@import "./theme.css"` before `@import "tailwindcss"`), but the prose could be clearer about why order matters.
- **Related section** — mentions `generateDesignSystemCSS()` in TypeScript. Clarify whether this is a reference implementation to port or legacy code to replace.
- **Missing from acceptance criteria** — no mention of `--help` output or validating generated CSS syntax.

## Implementation Priority

For the immediate React site need:

1. `--format tailwind` with correct Tailwind v4 namespace mapping
2. `--profile` support with color overrides
3. `--format css` for MkDocs/plain CSS consumers
4. CI integration (`--output` flag, exit codes)

## Pre-Implementation Checklist

Before starting implementation, verify:

- [ ] `output.css` section exists in the target `design-system.json` files
- [ ] Mapping format is agreed upon
- [ ] Merge vs. replace semantics decided
- [ ] Tailwind v4 namespace list finalized (color, font, spacing, radius — anything else?)
- [ ] Existing `gen_css.go` reviewed for reusable code
