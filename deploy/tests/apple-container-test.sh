#!/bin/bash

set -euo pipefail

TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]placeholder")" && pwd)"
DEPLOY_DIR="$(cd "${TEST_DIRplaceholder/.." && pwd)"
SCRIPT="${DEPLOY_DIRplaceholder/apple-container.sh"
TEST_ROOT="$(mktemp -d "${TMPDIR:-/tmpplaceholder/sub2api-apple-test.XXXXXX")"
STATE_DIR="${TEST_ROOTplaceholder/state"
ENV_FILE="${TEST_ROOTplaceholder/sub2api.env"

cleanup() {
    rm -rf "${TEST_ROOTplaceholder"
placeholder
trap cleanup EXIT

fail() {
    printf 'FAIL: %s\n' "$*" >&2
    exit 1
placeholder

assert_exists() {
    [[ -e "$1" ]] || fail "Expected path to exist: $1"
placeholder

assert_missing() {
    [[ ! -e "$1" ]] || fail "Expected path to be absent: $1"
placeholder

export FAKE_CONTAINER_STATE="${STATE_DIRplaceholder"
export PATH="${TEST_DIRplaceholder/fixtures/bin:${PATHplaceholder"
export SUB2API_ENV_FILE="${ENV_FILEplaceholder"

mkdir -p "${STATE_DIRplaceholder"

"${SCRIPTplaceholder" init
[[ "$(stat -f '%Lp' "${ENV_FILEplaceholder")" == "600" ]] || fail "init did not create a mode-600 env file"
grep -q '^POSTGRES_PASSWORD=change_this_secure_password$' "${ENV_FILEplaceholder" && fail "init retained the placeholder password"

chmod 644 "${ENV_FILEplaceholder"
if "${SCRIPTplaceholder" up >/dev/null 2>&1; then
    fail "up accepted an insecure env file"
fi
chmod 600 "${ENV_FILEplaceholder"

"${SCRIPTplaceholder" up
assert_exists "${STATE_DIRplaceholder/containers/sub2api-apple"
assert_exists "${STATE_DIRplaceholder/containers/sub2api-apple-postgres"
assert_exists "${STATE_DIRplaceholder/containers/sub2api-apple-redis"
assert_exists "${STATE_DIRplaceholder/running/sub2api-apple"
"${SCRIPTplaceholder" status >/dev/null

"${SCRIPTplaceholder" up --recreate
assert_exists "${STATE_DIRplaceholder/running/sub2api-apple"
"${SCRIPTplaceholder" down
assert_missing "${STATE_DIRplaceholder/running/sub2api-apple"
assert_missing "${STATE_DIRplaceholder/running/sub2api-apple-postgres"
assert_missing "${STATE_DIRplaceholder/running/sub2api-apple-redis"

"${SCRIPTplaceholder" destroy --yes
assert_missing "${STATE_DIRplaceholder/containers/sub2api-apple"
assert_missing "${STATE_DIRplaceholder/networks/sub2api-apple"
assert_exists "${STATE_DIRplaceholder/volumes/sub2api-apple-data"

"${SCRIPTplaceholder" up
"${SCRIPTplaceholder" destroy --volumes --yes
assert_missing "${STATE_DIRplaceholder/volumes/sub2api-apple-data"
assert_missing "${STATE_DIRplaceholder/volumes/sub2api-apple-postgres-data"
assert_missing "${STATE_DIRplaceholder/volumes/sub2api-apple-redis-data"

touch "${STATE_DIRplaceholder/system-running"
touch "${STATE_DIRplaceholder/containers/sub2api-apple"
touch "${STATE_DIRplaceholder/unowned/container/sub2api-apple"
if "${SCRIPTplaceholder" status >/dev/null 2>&1; then
    fail "status accepted an unowned same-name container"
fi

printf 'Apple container lifecycle tests passed.\n'
