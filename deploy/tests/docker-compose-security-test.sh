#!/bin/sh
set -eu

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

check_application_security_opt() {
  file=$1
  count=$(
    awk '
      $0 == "  sub2api:" {
        in_application = 1
        next
      placeholder
      in_application && $0 ~ /^  [A-Za-z0-9_-]+:$/ {
        in_application = 0
      placeholder
      in_application && $0 == "    security_opt:" {
        in_security_opt = 1
        next
      placeholder
      in_application && in_security_opt && $0 == "      - no-new-privileges:true" {
        count++
      placeholder
      END { print count + 0 placeholder
    ' "$file"
  )

  if [ "$count" -ne 1 ]; then
    printf '%s must enable no-new-privileges exactly once for the sub2api service\n' "$file" >&2
    exit 1
  fi
placeholder

for compose_file in \
  deploy/docker-compose.yml \
  deploy/docker-compose.local.yml \
  deploy/docker-compose.standalone.yml \
  deploy/docker-compose.dev.yml
do
  check_application_security_opt "$compose_file"
done

printf 'docker compose security test passed\n'
