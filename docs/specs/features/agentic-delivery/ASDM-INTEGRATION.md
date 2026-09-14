# SystemSpec: Design System: Enabling ASDM Level 6+ Delivery

This document describes how systemspec-designsystem enables organizations to progress from Level 5 (Agentic Engineering) to Level 6 (Autonomous Coding & Review) and beyond in the [Autonomous Software Delivery Model (ASDM)](https://productbuildershq.github.io/frameworks/software-delivery-autonomy/).

## ASDM Context

The Software Delivery Autonomy Model defines 7 levels of increasing autonomy:

| Level | Name | Human Role | Key Constraint |
|-------|------|------------|----------------|
| 5 | Agentic Engineering | Batch reviewer | Human review bottleneck |
| 6 | Autonomous Coding & Review | Specification owner | Scenario-based validation required |
| 7 | Autonomous Operations | Governor | Production telemetry validation |

**The Level 5→6 transition is the hardest**: it requires replacing human code review with automated scenario-based validation. This is where systemspec-designsystem provides critical infrastructure.

## The Level 5 Bottleneck

At Level 5, AI agents generate code autonomously for extended periods (hours), but every change still requires human review before merge. The AWS Project Mantle team described this as "queue prompts before sleep, check results in morning."

**The constraint**: Human review throughput limits deployment velocity. More agents → more PRs → review backlog → delayed delivery.

**The solution**: Replace human code review with specification-satisfying validation. If a change satisfies the specification scenarios, it can merge without human review.

## How Design-System-Spec Enables Level 6

Design-system-spec provides the three pillars required for Level 6 UI development:

### 1. Machine-Readable Specifications

The design system specification (`design-system.yaml`) defines:

```yaml
components:
  button:
    props:
      variant:
        type: enum
        values: [primary, secondary, ghost]
      size:
        type: enum
        values: [sm, md, lg]
    states: [default, hover, focus, disabled]
    accessibility:
      role: button
      ariaRequired: [aria-label]
```

This specification is the **source of truth** that agents implement against. Humans own the specification; agents own the implementation.

### 2. Scenario-Based Validation (Visual Regression Tests)

Visual regression tests are **scenarios** in ASDM terminology:

```yaml
# visual-tests.yaml
tests:
  - id: button-primary-default
    component: button
    variant: primary
    url: "/storybook/button--primary"
    viewports: [desktop, mobile]
    threshold: 0.001  # 0.1% pixel difference allowed
```

Each test captures a user-visible scenario. The test suite becomes a **holdout set** that validates whether generated code satisfies the specification.

### 3. Probabilistic Satisfaction Metrics

Visual diff thresholds provide **probabilistic satisfaction**:

```
Scenario: button-primary-default @ desktop
Baseline: v2.1.0
Actual diff: 0.0003 (0.03%)
Threshold: 0.001 (0.1%)
Status: SATISFIED
```

A component satisfies the specification when **all visual scenarios pass within threshold**. This replaces binary test pass/fail with graded confidence.

## The Level 6 Workflow with Design-System-Spec

```
┌─────────────────────────────────────────────────────────────┐
│                    SPECIFICATION LAYER                       │
│  Humans own: design-system.yaml, visual-tests.yaml          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    GENERATION LAYER                          │
│  Agents implement components using spec + LLM context       │
│  dss generate --llm → implementation guidance               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    VALIDATION LAYER                          │
│  Visual regression tests validate against baselines          │
│  dss visual test → scenario satisfaction report             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    MERGE DECISION                            │
│  All scenarios satisfied? → Auto-merge (no human review)    │
│  Scenarios failed? → Reject or escalate                     │
└─────────────────────────────────────────────────────────────┘
```

**Key insight**: The design system specification + visual test suite becomes the **Digital Twin** of your UI. Agents generate code, run it against the visual scenarios, and iterate until satisfaction metrics pass.

## Implementing Level 6 with Design-System-Spec

### Phase 1: Establish Baseline Infrastructure

```bash
# 1. Define design system specification
dss init --name "Acme Design System"

# 2. Create visual test suite
# visual-tests.yaml covers all component variants/states/viewports

# 3. Generate initial baselines from reference implementation
dss visual baseline generate v1.0.0

# 4. Establish satisfaction thresholds
# Default: 0.1% pixel diff tolerance
```

### Phase 2: Agent Implementation Loop

```bash
# Agent workflow (autonomous, no human in loop)
while not satisfied:
    # 1. Generate implementation from spec
    implementation = agent.implement(spec, llm_context)

    # 2. Render in test environment
    deploy_to_storybook(implementation)

    # 3. Run visual validation
    report = dss visual test --baseline v1.0.0 --json

    # 4. Check satisfaction
    if report.all_passed:
        satisfied = true
    else:
        # Feed failures back to agent for iteration
        agent.improve(implementation, report.failures)
```

### Phase 3: Autonomous Merge

```yaml
# CI/CD pipeline (no human approval gate)
on: pull_request

jobs:
  validate:
    steps:
      - run: dss visual test --baseline ${{ env.BASELINE_VERSION }}

      - name: Auto-merge if satisfied
        if: success()
        run: gh pr merge --auto --squash

      - name: Escalate if failed
        if: failure()
        run: |
          gh pr comment --body "Visual regression detected. Escalating to specification owner."
          gh pr edit --add-label "needs-spec-review"
```

## Comparison: Level 5 vs Level 6 Review

| Aspect | Level 5 (Human Review) | Level 6 (Scenario Validation) |
|--------|------------------------|------------------------------|
| Review trigger | Every PR | Only on scenario failure |
| Review focus | Code quality, style | Specification completeness |
| Reviewer role | Approve/reject code | Define/update specifications |
| Throughput limit | Human review capacity | CI pipeline capacity |
| Feedback loop | Hours (human availability) | Minutes (automated) |

## Extending to Level 7: Autonomous Operations

At Level 7, visual validation extends to **production monitoring**:

```yaml
# Production visual monitoring
monitors:
  - name: homepage-hero
    url: "https://app.example.com/"
    selector: "[data-testid='hero']"
    baseline: production-v2.1.0
    threshold: 0.005
    schedule: "*/15 * * * *"  # Every 15 minutes

    on_regression:
      - alert: pagerduty
      - action: rollback_deployment
```

**Production telemetry as validation**: If live UI diverges from baseline beyond threshold, autonomous systems can rollback, generate fixes, validate, and redeploy—without human intervention.

## ASDM Practice Mapping

| ASDM Practice | Design-System-Spec Implementation |
|---------------|-----------------------------------|
| Specification-first | `design-system.yaml` as source of truth |
| Scenario-based validation | `visual-tests.yaml` test suites |
| Probabilistic satisfaction | Pixel diff thresholds (0.001 = 0.1%) |
| Digital Twin testing | Storybook + w3pilot browser automation |
| Zero human code review | Auto-merge on scenario satisfaction |
| Full provenance | Baseline versions, manifest checksums |
| Governance-by-policy | Threshold configs, escalation rules |

## Investment Requirements

### Level 5 → Level 6 Prerequisites

| Requirement | Design-System-Spec Solution |
|-------------|----------------------------|
| Comprehensive scenario library | Visual test suite covering all variants/states |
| Baseline infrastructure | Versioned baselines with manifests |
| Satisfaction metrics | Configurable diff thresholds |
| Validation pipeline | `dss visual test` in CI |
| Escalation policies | Label-based routing on failure |

### Token Economics

Level 6 requires significant compute for iterative generation:

```
Agent iteration loop:
  - LLM tokens for generation: ~$0.10/component
  - Browser automation: ~$0.02/screenshot
  - Comparison compute: ~$0.01/diff

Per component validation: ~$0.15
Full design system (100 components × 3 variants × 3 viewports): ~$135/run
Daily iteration budget: ~$500-1,000
```

This aligns with ASDM's "$1,000+/day per engineer in tokens" estimate for Level 6.

## When to Adopt

**Prerequisites for Level 6 UI delivery**:

1. Design system specification is complete and machine-readable
2. Visual test suite covers >90% of component variants/states
3. Baseline infrastructure is operational
4. CI pipeline supports autonomous merge
5. Escalation policies are defined for specification owners

**When NOT to adopt**:

- Design system is still evolving rapidly (specification churn)
- Visual testing infrastructure is immature
- Organization requires human attestation (regulatory)
- Token budget is constrained

## Summary

Design-system-spec provides the **scenario-based validation infrastructure** required for Level 6 autonomous UI delivery:

1. **Specifications** define what components should be
2. **Visual tests** capture scenarios that must be satisfied
3. **Diff thresholds** provide probabilistic satisfaction metrics
4. **Baselines** establish the reference implementation
5. **Automated validation** replaces human code review

This enables the transition from "human reviews every PR" (Level 5) to "scenarios validate every change" (Level 6), removing the human review bottleneck and enabling true autonomous software delivery for UI components.

## Companion Projects

For complete Level 6 UI validation, combine systemspec-designsystem with:

| Project | Validation Type | Specification |
|---------|----------------|---------------|
| **systemspec-designsystem** | Visual regression | Design system spec |
| **[agent-a11y](https://github.com/plexusone/agent-a11y)** | Accessibility | WCAG 2.x |

Both projects integrate with multi-agent-spec for GO/WARN/NO-GO decision orchestration.

## References

- [Software Delivery Autonomy Levels](https://productbuildershq.github.io/frameworks/software-delivery-autonomy/) - ASDM framework
- [agent-a11y ASDM Integration](https://github.com/plexusone/agent-a11y/docs/architecture/asdm-integration.md) - Accessibility validation
- [Visual Regression Testing TRD](../visual-regression-testing/TRD.md) - Technical implementation
- [Visual Regression Testing PLAN](../visual-regression-testing/PLAN.md) - Implementation roadmap
