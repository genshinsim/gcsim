# CLAUDE.md

## Skills

Project skills live in `.claude/skills/` and work without any external plugin installed:

- **`implement-with-pr`** (`/implement-with-pr`): implement a ticket on its own branch and open a PR against `upstream`.
- **`tdd`** (`/tdd`): red → green test-first loop.
- **`review-changes`** (`/review-changes`): two-axis (Standards + Spec) review of a diff. Named to avoid colliding with the built-in `/code-review`.
- **`triage`** (`/triage`): move upstream issues and external PRs through the triage state machine.
- **`grilling`** (`/grilling`): stress-test a plan or decision with rounds of questions until shared understanding.
- **`domain-modeling`** (`/domain-modeling`): build and sharpen the glossary and ADRs (`CONTEXT.md`, `docs/adr/`).

These skills are vendored copies of the `mattpocock-skills` plugin. `tdd`, `review-changes` and `triage` are adapted for this repo's fork workflow and the `docs/agents/` config below; `grilling` and `domain-modeling` are verbatim.

## Agent skills

### Issue tracker

Issues are tracked as GitHub issues on the **upstream** repo `genshinsim/gcsim` (not the `origin` fork); all `gh` commands must pass `--repo genshinsim/gcsim`. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical triage labels (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`), created and applied on upstream `genshinsim/gcsim`. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: one `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.
