set -euo pipefail

GIT_DIR=$(git rev-parse --git-dir)

if [[ -d "$GIT_DIR/rebase-merge" || -d "$GIT_DIR/rebase-apply" || -f "$GIT_DIR/MERGE_HEAD" ]]; then
	exit 0
fi

die() {
	echo "$1" >&2
	exit 1
}

readonly COMMIT_MSG_FILE=$1
readonly COMMIT_MSG="$(sed "/^#/d" "$COMMIT_MSG_FILE")"

readonly JIRA_TICKET_REGEX="([A-Z]+-[0-9]+)"
readonly BRANCH="$(git symbolic-ref --short --quiet HEAD || true)"
[[ -z "$BRANCH" ]] && die "Detached HEAD is not supported"
[[ "$BRANCH" =~ $JIRA_TICKET_REGEX ]] || die "INVALID BRANCH NAME. SHOULD BE LIKE
feature/TODO-123-short-description
WHERE TODO-123 is Jira ticket"
readonly JIRA_TICKET="${BASH_REMATCH[0]}"

readonly COMMIT_HEADER_AND_BODY_SEPARATOR="$(tail -n +2 "$COMMIT_MSG_FILE" | head -n1)"
body="$(tail -n +3 "$COMMIT_MSG_FILE")"

if [[ -n "$body" && -n "$COMMIT_HEADER_AND_BODY_SEPARATOR" ]]; then
	die "Invalid commit format. Proper commit format:
==================================
feat(TODO-123): brief explanation

detailed explanation

==============OR================

feat(TODO-123): brief explanation
=================================="
fi

readonly COMMIT_HEADER_REGEX="^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)\($JIRA_TICKET\): .{1,50}$"
readonly COMMIT_HEADER="$(head -n1 "$COMMIT_MSG_FILE")"
[[ "$COMMIT_HEADER" =~ $COMMIT_HEADER_REGEX ]] || die "Invalid commit header format. Should be like:
feat(TODO-123): brief explanation"
