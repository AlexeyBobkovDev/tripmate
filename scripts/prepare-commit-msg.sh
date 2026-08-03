#!/usr/bin/env bash

# shellcheck source=lib/die.sh
. lib/die.sh

GIT_DIR=$(git rev-parse --git-dir)

if [[ -d "$GIT_DIR/rebase-merge" || -d "$GIT_DIR/rebase-apply" || -f "$GIT_DIR/MERGE_HEAD" ]]; then
	exit 0
fi

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

declare -A service_names=()
while IFS= read -r file; do
	echo "$file"
	if [[ "$file" =~ ^services/([^/]+) ]]; then
		service="${BASH_REMATCH[1]}"
	else
		service="service"
	fi

	service_names["$service"]=1
	echo "${!service_names[@]}"
done < <(git diff --cached --name-only)

if (( "${#service_names[@]}" != 1 )); then
	die "You have to split your commits into small chunks!"
fi

SERVICE_NAME="${!service_names[@]}"

OUTPUT="feat(${JIRA_TICKET}, ${SERVICE_NAME}): $COMMIT_MSG"

echo "$OUTPUT" >"$COMMIT_MSG_FILE"
