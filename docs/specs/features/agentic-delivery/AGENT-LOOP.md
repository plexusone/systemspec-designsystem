# Agentic Delivery Loop with Design-System-Spec

This document describes the autonomous agent workflow for implementing UI components using systemspec-designsystem as the validation layer.

## Overview

The agentic delivery loop replaces human code review with specification-based validation:

```
┌──────────────────────────────────────────────────────────────────┐
│                         AGENT LOOP                                │
│                                                                   │
│   ┌─────────┐    ┌──────────┐    ┌──────────┐    ┌─────────┐    │
│   │  Spec   │───▶│ Generate │───▶│ Validate │───▶│ Decide  │    │
│   │  Read   │    │   Code   │    │  Visual  │    │         │    │
│   └─────────┘    └──────────┘    └──────────┘    └────┬────┘    │
│                                                        │         │
│                       ◀────────────────────────────────┘         │
│                              (iterate if failed)                  │
└──────────────────────────────────────────────────────────────────┘
```

## Loop Phases

### Phase 1: Specification Read

The agent retrieves the component specification and LLM context:

```bash
# Get component specification
dss info --component button --json

# Generate LLM-optimized implementation guidance
dss generate --llm --component button
```

Output provides:

- Props with types, defaults, constraints
- Valid states and transitions
- Accessibility requirements
- Anti-patterns to avoid
- Token references (colors, spacing, typography)

### Phase 2: Code Generation

The agent implements the component using the specification:

```typescript
// Agent prompt structure
const prompt = `
Implement a React component for ${componentId} following this specification:

${llmContext}

Requirements:
- Use design tokens from the specification
- Implement all variants: ${variants.join(', ')}
- Handle all states: ${states.join(', ')}
- Meet accessibility requirements: ${a11y}

Output only the component code, no explanations.
`;

const implementation = await llm.generate(prompt);
```

### Phase 3: Visual Validation

The agent deploys to a test environment and runs visual regression:

```bash
# Deploy to Storybook or similar
npm run build-storybook
npx serve storybook-static -p 6006 &

# Run visual tests against baseline
dss visual test \
  --tests button \
  --baseline v2.1.0 \
  --json > validation-result.json
```

### Phase 4: Decision

Based on validation results:

```javascript
const result = JSON.parse(fs.readFileSync('validation-result.json'));

if (result.summary.failed === 0 && result.summary.errors === 0) {
  // All scenarios satisfied - proceed to merge
  await git.commit('feat(button): implement button component');
  await git.push();
  await github.createPR({ autoMerge: true });
} else {
  // Scenarios failed - iterate
  const failures = result.results.filter(r => r.status === 'failed');
  await agent.improve(implementation, failures);
}
```

## Complete Agent Implementation

```python
#!/usr/bin/env python3
"""
Autonomous component implementation agent.
Implements ASDM Level 6 workflow with systemspec-designsystem validation.
"""

import json
import subprocess
import sys
from pathlib import Path

class ComponentAgent:
    def __init__(self, spec_dir: str, component_id: str, baseline: str):
        self.spec_dir = Path(spec_dir)
        self.component_id = component_id
        self.baseline = baseline
        self.max_iterations = 5
        self.iteration = 0

    def run(self) -> bool:
        """Execute the agent loop until success or max iterations."""
        while self.iteration < self.max_iterations:
            self.iteration += 1
            print(f"\n=== Iteration {self.iteration} ===")

            # 1. Read specification
            spec = self.read_spec()
            llm_context = self.get_llm_context()

            # 2. Generate implementation
            code = self.generate_code(spec, llm_context)
            self.write_code(code)

            # 3. Build and deploy
            if not self.build():
                continue

            # 4. Validate visually
            result = self.validate()

            # 5. Check satisfaction
            if self.is_satisfied(result):
                print(f"\n✓ Component {self.component_id} satisfied in {self.iteration} iterations")
                return True
            else:
                print(f"\n✗ Validation failed, iterating...")
                self.feedback = result

        print(f"\n✗ Max iterations reached without satisfaction")
        return False

    def read_spec(self) -> dict:
        """Read component specification."""
        result = subprocess.run(
            ['dss', 'info', '--component', self.component_id, '--json', '-d', str(self.spec_dir)],
            capture_output=True, text=True
        )
        return json.loads(result.stdout)

    def get_llm_context(self) -> str:
        """Get LLM-optimized implementation guidance."""
        result = subprocess.run(
            ['dss', 'generate', '--llm', '--component', self.component_id, '-d', str(self.spec_dir)],
            capture_output=True, text=True
        )
        return result.stdout

    def generate_code(self, spec: dict, llm_context: str) -> str:
        """Generate component code using LLM."""
        prompt = self.build_prompt(spec, llm_context)

        # Use Claude, GPT-4, or other LLM
        # This is a placeholder - integrate your preferred LLM
        code = self.call_llm(prompt)
        return code

    def build_prompt(self, spec: dict, llm_context: str) -> str:
        """Build the generation prompt."""
        prompt = f"""
Implement a React component for {self.component_id}.

## Specification
{llm_context}

## Requirements
- TypeScript with proper types
- Use CSS modules or styled-components
- Implement all variants and states
- Follow accessibility requirements
- Use design tokens (do not hardcode colors/spacing)

## Output
Return only the component code. No explanations or markdown.
"""

        # Add feedback from previous iteration if available
        if hasattr(self, 'feedback') and self.feedback:
            failures = [r for r in self.feedback.get('results', []) if r['status'] == 'failed']
            if failures:
                prompt += f"""

## Previous Iteration Feedback
The following visual tests failed:
"""
                for f in failures:
                    prompt += f"- {f['testId']} @ {f['viewport']}: {f['diffPercent']*100:.2f}% diff\n"
                prompt += "\nAdjust the implementation to match the baseline more closely."

        return prompt

    def call_llm(self, prompt: str) -> str:
        """Call LLM API. Override with your LLM integration."""
        # Placeholder - integrate Claude API, OpenAI, etc.
        raise NotImplementedError("Integrate your LLM API here")

    def write_code(self, code: str):
        """Write generated code to file."""
        output_path = self.spec_dir / 'src' / 'components' / f'{self.component_id}.tsx'
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(code)

    def build(self) -> bool:
        """Build the component for testing."""
        result = subprocess.run(
            ['npm', 'run', 'build-storybook'],
            cwd=self.spec_dir,
            capture_output=True
        )
        return result.returncode == 0

    def validate(self) -> dict:
        """Run visual regression tests."""
        result = subprocess.run(
            ['dss', 'visual', 'test',
             '--tests', self.component_id,
             '--baseline', self.baseline,
             '--json',
             '-d', str(self.spec_dir)],
            capture_output=True, text=True
        )
        return json.loads(result.stdout)

    def is_satisfied(self, result: dict) -> bool:
        """Check if all scenarios are satisfied."""
        summary = result.get('summary', {})
        return summary.get('failed', 0) == 0 and summary.get('errors', 0) == 0


def main():
    if len(sys.argv) < 4:
        print("Usage: agent.py <spec-dir> <component-id> <baseline-version>")
        sys.exit(1)

    agent = ComponentAgent(
        spec_dir=sys.argv[1],
        component_id=sys.argv[2],
        baseline=sys.argv[3]
    )

    success = agent.run()
    sys.exit(0 if success else 1)


if __name__ == '__main__':
    main()
```

## CI/CD Integration

### GitHub Actions Workflow

```yaml
name: Autonomous Component Delivery

on:
  workflow_dispatch:
    inputs:
      component:
        description: 'Component to implement'
        required: true
      baseline:
        description: 'Baseline version'
        required: true
        default: 'latest'

jobs:
  implement:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup
        run: |
          npm ci
          go install github.com/plexusone/systemspec-designsystem/cmd/dss@latest

      - name: Run Agent Loop
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          python3 scripts/agent.py . ${{ inputs.component }} ${{ inputs.baseline }}

      - name: Create PR
        if: success()
        run: |
          git checkout -b agent/${{ inputs.component }}-$(date +%s)
          git add .
          git commit -m "feat(${{ inputs.component }}): autonomous implementation"
          git push -u origin HEAD
          gh pr create \
            --title "feat(${{ inputs.component }}): autonomous implementation" \
            --body "Implemented by autonomous agent. All visual scenarios satisfied." \
            --label "autonomous,auto-merge"

      - name: Auto-merge
        if: success()
        run: gh pr merge --auto --squash
```

### Merge Policy

```yaml
# .github/branch-protection.yaml
branches:
  main:
    required_status_checks:
      strict: true
      contexts:
        - visual-regression

    # No required reviewers for autonomous PRs
    required_pull_request_reviews:
      bypass_pull_request_actors:
        - autonomous-agent[bot]

    # Labels that bypass review
    auto_merge_labels:
      - autonomous
      - visual-satisfied
```

## Convergence Metrics

Track agent performance over time:

```yaml
# metrics.yaml
agent_runs:
  - component: button
    iterations: 3
    duration_minutes: 12
    token_cost: 0.45
    satisfied: true

  - component: card
    iterations: 5
    duration_minutes: 28
    token_cost: 1.20
    satisfied: false
    failure_reason: "baseline_missing"
```

**Key metrics**:

| Metric | Target | Description |
|--------|--------|-------------|
| Convergence rate | >90% | % of runs that satisfy within max iterations |
| Mean iterations | <3 | Average iterations to satisfaction |
| Cost per component | <$1 | Token cost per successful implementation |
| Cycle time | <30min | Time from start to PR merge |

## Escalation Policies

When autonomous delivery fails:

```yaml
# escalation.yaml
policies:
  - condition: iterations >= max_iterations
    action: create_issue
    assignee: "@design-system-team"
    labels: ["agent-escalation", "needs-spec-review"]

  - condition: visual_diff > 0.05  # >5% diff
    action: request_baseline_update
    message: "Significant visual change detected. Update baseline?"

  - condition: error_type == "baseline_missing"
    action: generate_baseline
    version: "auto-$(date)"
```

## Best Practices

### 1. Comprehensive Baseline Coverage

Ensure baselines cover all variants and states:

```bash
# Check coverage
dss visual baseline list --verbose

# Generate missing baselines
dss visual baseline generate v2.2.0 --tests button,card,input
```

### 2. Appropriate Thresholds

Configure thresholds based on component stability:

```yaml
# visual-tests.yaml
tests:
  - id: button-primary
    threshold: 0.001  # 0.1% - strict for stable components

  - id: chart-dynamic
    threshold: 0.01   # 1% - looser for dynamic content
```

### 3. Feedback Loop Quality

Provide detailed failure information to agents:

```json
{
  "testId": "button-primary-hover",
  "viewport": "desktop",
  "status": "failed",
  "diffPercent": 0.023,
  "threshold": 0.001,
  "diffPath": "results/button-primary-hover.diff.png",
  "analysis": "Color mismatch in hover state. Expected #0066CC, found #0055BB."
}
```

### 4. Incremental Adoption

Start with stable, well-specified components:

1. **Week 1-2**: Run agent on 1-2 simple components with human oversight
2. **Week 3-4**: Expand to 5-10 components, review aggregate metrics
3. **Month 2**: Enable auto-merge for high-confidence components
4. **Month 3+**: Full autonomous delivery for entire design system

## Summary

The agentic delivery loop with systemspec-designsystem enables:

1. **Specification-driven generation**: Agents read specs, not source code
2. **Visual validation**: Scenarios replace human code review
3. **Iterative convergence**: Agents self-correct based on diff feedback
4. **Autonomous merge**: Satisfied scenarios trigger auto-merge
5. **Escalation by exception**: Humans handle specification changes, not implementation review

This workflow implements ASDM Level 6 for UI components, removing the human review bottleneck while maintaining quality through specification-based validation.
