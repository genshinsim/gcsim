#!/bin/bash
# Post-edit hook: runs turbo test for the ui-next package containing the edited file
INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

# Only process files within ui-next/
case "$FILE_PATH" in
  */ui-next/packages/*/src/*|*/ui-next/apps/*/src/*)
    ;;
  *)
    exit 0
    ;;
esac

# Extract package/app name from path
# e.g., .../ui-next/packages/viewer/src/... → viewer
# e.g., .../ui-next/apps/web/src/... → web
PKG_NAME=$(echo "$FILE_PATH" | sed -E 's|.*/ui-next/(packages|apps)/([^/]+)/.*|\2|')

if [ -z "$PKG_NAME" ]; then
  exit 0
fi

REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
cd "$REPO_ROOT/ui-next" || exit 0

turbo run test --filter="@gcsim/$PKG_NAME" 2>&1

exit 0
