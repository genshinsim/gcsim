#!/bin/bash
# Pre-commit hook: runs biome check --write on staged ui-next files before git commit
INPUT=$(cat)
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // empty')

# Only intercept git commit commands
if [[ "$COMMAND" != *"git commit"* ]]; then
  exit 0
fi

# Get staged files within ui-next/
STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACMR -- 'ui-next/apps' 'ui-next/packages' 'ui-next/tooling' 2>/dev/null)

if [ -z "$STAGED_FILES" ]; then
  exit 0
fi

# Run biome check --write on staged files
# Strip ui-next/ prefix since we cd into that directory
cd "$(git rev-parse --show-toplevel)/ui-next" || exit 0
echo "$STAGED_FILES" | sed 's|^ui-next/||' | xargs npx biome check --write 2>&1

if [ $? -ne 0 ]; then
  echo "Biome check failed on staged files. Fix issues before committing." >&2
  exit 2
fi

exit 0
