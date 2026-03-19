#!/bin/bash
# Cleanup stale git worktrees created during parallel agent development
# Worktree naming convention: worktree-phase{N}-{package-name}
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
if [ -z "$REPO_ROOT" ]; then
  echo "Error: not in a git repository" >&2
  exit 1
fi

echo "Checking for stale worktrees..."

# List all worktrees and find stale ones
STALE_COUNT=0
while IFS= read -r line; do
  WORKTREE_PATH=$(echo "$line" | awk '{print $1}')

  # Skip the main worktree
  if [ "$WORKTREE_PATH" = "$REPO_ROOT" ]; then
    continue
  fi

  # Only clean up worktrees matching our naming convention
  WORKTREE_NAME=$(basename "$WORKTREE_PATH")
  if [[ "$WORKTREE_NAME" != worktree-phase* ]]; then
    continue
  fi

  # Check if the worktree directory still exists
  if [ ! -d "$WORKTREE_PATH" ]; then
    echo "  Pruning missing worktree: $WORKTREE_NAME"
    git worktree prune
    STALE_COUNT=$((STALE_COUNT + 1))
    continue
  fi

  # Check if worktree has uncommitted changes
  if git -C "$WORKTREE_PATH" diff --quiet 2>/dev/null && \
     git -C "$WORKTREE_PATH" diff --cached --quiet 2>/dev/null; then
    echo "  Removing clean worktree: $WORKTREE_NAME ($WORKTREE_PATH)"
    git worktree remove "$WORKTREE_PATH"
    STALE_COUNT=$((STALE_COUNT + 1))
  else
    echo "  Skipping worktree with changes: $WORKTREE_NAME"
  fi
done < <(git worktree list --porcelain | grep "^worktree " | sed 's/^worktree //')

if [ "$STALE_COUNT" -eq 0 ]; then
  echo "No stale worktrees found."
else
  echo "Cleaned up $STALE_COUNT worktree(s)."
fi
