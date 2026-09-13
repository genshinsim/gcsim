# CLAUDE.md

## Skills

Project skills live in `.claude/skills/` and work without any external plugin installed:

- **`gcsim-implement`** (`/gcsim-implement`): implement a ticket on its own branch and open a PR against `upstream`.
- **`gcsim-tdd`** (`/gcsim-tdd`): red → green test-first loop.
- **`gcsim-review-changes`** (`/gcsim-review-changes`): two-axis (Standards + Spec) review of a diff. Named to avoid colliding with the built-in `/code-review`.
- **`gcsim-triage`** (`/gcsim-triage`): move upstream issues and external PRs through the triage state machine.
- **`gcsim-grilling`** (`/gcsim-grilling`): stress-test a plan or decision with rounds of questions until shared understanding.
- **`gcsim-domain-modeling`** (`/gcsim-domain-modeling`): build and sharpen the glossary and ADRs (`CONTEXT.md`, `docs/adr/`).

These skills are vendored copies of the `mattpocock-skills` plugin. `gcsim-tdd`, `gcsim-review-changes` and `gcsim-triage` are adapted for this repo's fork workflow and the `docs/agents/` config below; `gcsim-grilling` and `gcsim-domain-modeling` are verbatim.

## Agent skills

### Issue tracker

Issues are tracked as GitHub issues on the **upstream** repo `genshinsim/gcsim` (not the `origin` fork); all `gh` commands must pass `--repo genshinsim/gcsim`. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical triage labels (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`), created and applied on upstream `genshinsim/gcsim`. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: one `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.
