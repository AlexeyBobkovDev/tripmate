#!/usr/bin/env bash

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=scripts/lib/die.sh
. "$SCRIPT_DIR/lib/die.sh"

GIT_DIR=$(git rev-parse --git-dir)

if [[ -d "$GIT_DIR/rebase-merge" || -d "$GIT_DIR/rebase-apply" || -f "$GIT_DIR/MERGE_HEAD" ]]; then
	exit 0
fi

readonly COMMIT_MSG_FILE=$1
COMMIT_MSG="$(sed "/^#/d" "$COMMIT_MSG_FILE")"
readonly COMMIT_MSG

readonly JIRA_TICKET_REGEX="([A-Z]+-[0-9]+)"
BRANCH="$(git symbolic-ref --short --quiet HEAD || true)"
readonly BRANCH

[[ -z "$BRANCH" ]] && die "Detached HEAD is not supported"
[[ "$BRANCH" =~ $JIRA_TICKET_REGEX ]] || die "INVALID BRANCH NAME. SHOULD BE LIKE
feature/TODO-123-short-description
WHERE TODO-123 is Jira ticket"
readonly JIRA_TICKET="${BASH_REMATCH[0]}"

COMMIT_HEADER_AND_BODY_SEPARATOR="$(tail -n +2 <(echo "$COMMIT_MSG") | head -n1)"
readonly COMMIT_HEADER_AND_BODY_SEPARATOR
body="$(tail -n +3 <(echo "$COMMIT_MSG"))"

if [[ -n "$body" && -n "$COMMIT_HEADER_AND_BODY_SEPARATOR" ]]; then
	die "Invalid commit format. Proper commit format:
==================================
feat(TODO-123, scope): brief explanation

detailed explanation

==============OR================

feat(TODO-123, scope): brief explanation
=================================="
fi

readonly COMMIT_HEADER_REGEX="^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)\(${JIRA_TICKET}, .{,40}\): .{1,50}$"
COMMIT_HEADER="$(head -n1 <(echo "$COMMIT_MSG"))"
readonly COMMIT_HEADER
[[ "$COMMIT_HEADER" =~ $COMMIT_HEADER_REGEX ]] || die "Invalid commit header format. Should be like:
feat(TODO-123, scope): brief explanation"
