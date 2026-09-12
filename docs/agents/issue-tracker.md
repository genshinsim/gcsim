# Issue tracker: GitHub

Issues and specs for this repo live as GitHub issues on the **upstream** repository, **`genshinsim/gcsim`** — NOT the `origin` fork (`srliao/gcsim`). Use the `gh` CLI for all operations and always scope commands to upstream with `--repo genshinsim/gcsim`.

> This is a fork checkout: `origin` points at a personal fork and `upstream` at the canonical repo. All issues, labels, and triage live upstream. Without `--repo genshinsim/gcsim`, `gh` may default to the fork and act on the wrong repository.

## Conventions

- **Create an issue**: `gh issue create --repo genshinsim/gcsim --title "..." --body "..."`. Use a heredoc for multi-line bodies.
- **Read an issue**: `gh issue view <number> --repo genshinsim/gcsim --comments`, filtering comments by `jq` and also fetching labels.
- **List issues**: `gh issue list --repo genshinsim/gcsim --state open --json number,title,body,labels,comments --jq '[.[] | {number, title, body, labels: [.labels[].name], comments: [.comments[].body]}]'` with appropriate `--label` and `--state` filters.
- **Comment on an issue**: `gh issue comment <number> --repo genshinsim/gcsim --body "..."`
- **Apply / remove labels**: `gh issue edit <number> --repo genshinsim/gcsim --add-label "..."` / `--remove-label "..."`
- **Close**: `gh issue close <number> --repo genshinsim/gcsim --comment "..."`

Always pass `--repo genshinsim/gcsim`. Do not rely on `gh`'s auto-detection here, because it can resolve to the `origin` fork.

## Pull requests as a triage surface

**PRs as a request surface: no.** _(Set to `yes` if this repo treats external PRs as feature requests; `/gcsim-triage` reads this flag.)_

When set to `yes`, PRs run through the same labels and states as issues, using the `gh pr` equivalents (all scoped with `--repo genshinsim/gcsim`):

- **Read a PR**: `gh pr view <number> --repo genshinsim/gcsim --comments` and `gh pr diff <number> --repo genshinsim/gcsim` for the diff.
- **List external PRs for triage**: `gh pr list --repo genshinsim/gcsim --state open --json number,title,body,labels,author,authorAssociation,comments` then keep only `authorAssociation` of `CONTRIBUTOR`, `FIRST_TIME_CONTRIBUTOR`, or `NONE` (drop `OWNER`/`MEMBER`/`COLLABORATOR`).
- **Comment / label / close**: `gh pr comment`, `gh pr edit --add-label`/`--remove-label`, `gh pr close` (each with `--repo genshinsim/gcsim`).

GitHub shares one number space across issues and PRs, so a bare `#42` may be either: resolve with `gh pr view 42 --repo genshinsim/gcsim` and fall back to `gh issue view 42 --repo genshinsim/gcsim`.

## When a skill says "publish to the issue tracker"

Create a GitHub issue on `genshinsim/gcsim`.

## When a skill says "fetch the relevant ticket"

Run `gh issue view <number> --repo genshinsim/gcsim --comments`.

## Wayfinding operations

Used by `/wayfinder`. The **map** is a single issue with **child** issues as tickets. All operations target `genshinsim/gcsim`.

- **Map**: a single issue labelled `wayfinder:map`, holding the Notes / Decisions-so-far / Fog body. `gh issue create --repo genshinsim/gcsim --label wayfinder:map`.
- **Child ticket**: an issue linked to the map as a GitHub sub-issue (`gh api` on the sub-issues endpoint). Where sub-issues aren't enabled, add the child to a task list in the map body and put `Part of #<map>` at the top of the child body. Labels: `wayfinder:<type>` (`research`/`prototype`/`grilling`/`task`). Once claimed, the ticket is assigned to the driving dev.
- **Blocking**: GitHub's **native issue dependencies**, the canonical, UI-visible representation. Add an edge with `gh api --method POST repos/genshinsim/gcsim/issues/<child>/dependencies/blocked_by -F issue_id=<blocker-db-id>`, where `<blocker-db-id>` is the blocker's numeric **database id** (`gh api repos/genshinsim/gcsim/issues/<n> --jq .id`, _not_ the `#number` or `node_id`). GitHub reports `issue_dependencies_summary.blocked_by` (open blockers only, the live gate). Where dependencies aren't available, fall back to a `Blocked by: #<n>, #<n>` line at the top of the child body. A ticket is unblocked when every blocker is closed.
- **Frontier query**: list the map's open children (`gh issue list --repo genshinsim/gcsim --state open`, scoped to the map's sub-issues / task list), drop any with an open blocker (`issue_dependencies_summary.blocked_by > 0`, or an open issue in the `Blocked by` line) or an assignee; first in map order wins.
- **Claim**: `gh issue edit <n> --repo genshinsim/gcsim --add-assignee @me`, the session's first write.
- **Resolve**: `gh issue comment <n> --repo genshinsim/gcsim --body "<answer>"`, then `gh issue close <n> --repo genshinsim/gcsim`, then append a context pointer (gist + link) to the map's Decisions-so-far.
