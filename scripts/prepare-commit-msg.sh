set -euo pipefail

GIT_DIR=$(git rev-parse --git-dir)

if [[ -d "$GIT_DIR/rebase-merge" || -d "$GIT_DIR/rebase-apply" || -f "$GIT_DIR/MERGE_HEAD" ]]; then
	exit 0
fi

die() {
	echo "$1" >&2
	exit 1
}

readonly BRANCH="$(git symbolic-ref --short --quiet HEAD || true)"
[[ -z "$BRANCH" ]] && die "Detached HEAD is not supported"
readonly COMMIT_MSG_FILE=$1
readonly COMMIT_SOURCE=${2:-}
readonly COMMIT_TYPE="^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)\("
readonly JIRA_TICKET_REGEX="([A-Z]+-[0-9]+)"

COMMIT_MSG="$(sed "/^#/d" "$COMMIT_MSG_FILE")"
[[ "$COMMIT_MSG" =~ $COMMIT_TYPE ]] && exit 0

case "$COMMIT_SOURCE" in
merge | squash)
	exit 0
	;;
esac

[[ $BRANCH =~ $JIRA_TICKET_REGEX ]] || die "Missing Jira ticket in branch name"
JIRA_TICKET="${BASH_REMATCH[0]}"

OUTPUT="feat(${JIRA_TICKET}): $COMMIT_MSG"

echo "$OUTPUT" >"$COMMIT_MSG_FILE"
