---
name: gcsim-implement
description: "Implement a ticket on its own branch and open a PR linked to the issue. Use when the work should land as a reviewable PR rather than commits on the current branch."
disable-model-invocation: true
---

# Implement with PR

Implements a ticket, with the branch and PR steps around it.

This repo uses the **fork workflow**. Your checkout has two remotes: `origin` is your personal
fork, and `upstream` is the canonical repo. Always: branch off **`upstream/main`**, do the work
on a branch in **your fork (`origin`)**, and open the PR **against `upstream`**. See
`docs/agents/issue-tracker.md`.

## 1. Read the ticket

The user names an issue (a number, a URL, or a description). Issues live on the `upstream`
repo. Read it in full with `gh issue view <number> --comments`, including the issues it is
blocked by. If `gh` targets your fork by default, pass `--repo <upstream owner>/<repo>` so it
reads from upstream.

If it has an open blocker, say so and stop.

## 2. Branch off `upstream/main`

Before writing any code:

```
git fetch upstream
git switch --create <type>/<issue-number>-<slug> upstream/main
```

`<type>` is `feat`, `fix`, `chore`, `refactor` or `test` — whichever matches the ticket.
`<slug>` is two to four words from the title, kebab-case.

Always branch off the freshly fetched `upstream/main`, not `origin/main` — your fork can lag
behind upstream. Never implement on `main`.

## 3. Implement

Implement the work described in the ticket.

Use `/gcsim-tdd` where possible, at pre-agreed seams.

Run typechecking regularly, single test files regularly, and the full test suite once at the
end.

Commit regularly to the branch you created in step 2. Keep commits small and standalone where possible, allowing for easy review by others commit by commit.

Write commit messages tight: a `<type>(<scope>): <subject>` subject line, then point-form bullets of what changed — short, to the point, no prose paragraphs. Keep the `Co-Authored-By` trailer; drop any Claude session link.

Once done, use `/gcsim-review-changes` to review the work.

## 4. Open the PR

Before creating the PR, always check again if `upstream/main` has changed. If it has, rebase your branch onto it (resolve any conflicts) and re-run the tests:

```
git fetch upstream
git rebase upstream/main
```

Do not merge `main` into your branch — always rebase.

Push the branch to **your fork (`origin`)** and open a PR against **`upstream`**:

```
git push --set-upstream origin <branch>
gh pr create --base main --title "..." --body-file <file>
```

Run `gh pr create` from the fork clone. If it prompts for a base repo, pick `upstream`; if it
prompts where to push, pick `origin`.

Write the body like game patch notes: **New**, **Changed** and **Removed** sections (omit any that is empty), each a short bullet list of what a reviewer can now observe — not files touched. Keep it brief; add detail under a bullet only where absolutely necessary.

The body must also:

- carry `Closes #<issue-number>` on its own line, so merging closes the ticket
- name anything in the ticket's acceptance criteria you did **not** do, and why

Drop any Claude session link from the PR footer.

Keep the PR small. If the work outgrew the ticket, stop and say so rather than widening the PR.

## 5. Report

Give the user the PR URL and the issue it closes. Do not merge it, and do not close or edit
the issue — the merge does that.
