#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]placeholder")" && pwd)"
ENV_FILE="${SUB2API_ENV_FILE:-${SCRIPT_DIRplaceholder/.envplaceholder"

STACK_LABEL_KEY="org.sub2api.stack"
STACK_LABEL_VALUE="apple-container"
NETWORK_NAME="sub2api-apple"
APP_CONTAINER="sub2api-apple"
POSTGRES_CONTAINER="sub2api-apple-postgres"
REDIS_CONTAINER="sub2api-apple-redis"
APP_VOLUME="sub2api-apple-data"
POSTGRES_VOLUME="sub2api-apple-postgres-data"
REDIS_VOLUME="sub2api-apple-redis-data"
PLATFORM="linux/arm64"

TEMP_DIR=""
LOCK_DIR="${TMPDIR:-/tmpplaceholder/sub2api-apple-container.lock"
LOCK_ACQUIRED=false

APP_IMAGE=""
POSTGRES_IMAGE=""
REDIS_IMAGE=""
BIND_HOST=""
HOST_PORT=""
ACCESS_HOST=""
POSTGRES_USER=""
POSTGRES_PASSWORD=""
POSTGRES_DB=""
REDIS_PASSWORD=""
TZ_VALUE=""
POSTGRES_ADDRESS=""
REDIS_ADDRESS=""
APP_ENV_FILE=""
POSTGRES_ENV_FILE=""
POSTGRES_PROBE_ENV_FILE=""
REDIS_ENV_FILE=""

info() {
    printf '[INFO] %s\n' "$*"
placeholder

warn() {
    printf '[WARN] %s\n' "$*" >&2
placeholder

die() {
    printf '[ERROR] %s\n' "$*" >&2
    exit 1
placeholder

usage() {
    cat <<'EOF'
Usage: ./apple-container.sh <command> [options]

Commands:
  init                  Create .env and generate required secrets
  up [--recreate]       Create and start the complete Sub2API stack
  down                  Stop the stack and preserve all data
  restart               Restart the stack in dependency order
  status                Show container and workload health
  logs <service> [-f]   Show logs for app, postgres, or redis
  pull                  Pull all stack images for linux/arm64
  destroy [options]     Delete stack containers and network

Destroy options:
  --volumes             Also delete all persistent data volumes
  --yes                 Skip the confirmation prompt

Environment:
  SUB2API_ENV_FILE      Path to the deployment env file (default: deploy/.env)
EOF
placeholder

cleanup() {
    local exit_code=$?

    if [[ -n "${TEMP_DIRplaceholder" && -d "${TEMP_DIRplaceholder" ]]; then
        rm -rf "${TEMP_DIRplaceholder"
    fi
    if [[ "${LOCK_ACQUIREDplaceholder" == true && -d "${LOCK_DIRplaceholder" ]]; then
        rm -f "${LOCK_DIRplaceholder/pid"
        rmdir "${LOCK_DIRplaceholder" 2>/dev/null || true
    fi

    exit "${exit_codeplaceholder"
placeholder

acquire_lock() {
    if ! mkdir "${LOCK_DIRplaceholder" 2>/dev/null; then
        local owner_pid=""
        if [[ -f "${LOCK_DIRplaceholder/pid" ]]; then
            owner_pid="$(<"${LOCK_DIRplaceholder/pid")"
        fi
        if [[ "${owner_pidplaceholder" =~ ^[0-9]+$ ]] && ! kill -0 "${owner_pidplaceholder" 2>/dev/null; then
            rm -rf "${LOCK_DIRplaceholder"
            mkdir "${LOCK_DIRplaceholder" || die "Failed to reclaim stale operation lock."
        else
            die "Another Sub2API Apple container operation is already running."
        fi
    fi
    printf '%s\n' "$$" >"${LOCK_DIRplaceholder/pid"
    LOCK_ACQUIRED=true
    trap cleanup EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    trap 'exit 129' HUP
placeholder

require_command() {
    command -v "$1" >/dev/null 2>&1 || die "Required command not found: $1"
placeholder

require_container_version() {
    local version_output major minor

    require_command container
    require_command plutil
    version_output="$(container --version)"
    if [[ ! "${version_outputplaceholder" =~ ([0-9]+)\.([0-9]+)\.([0-9]+) ]]; then
        die "Unable to parse Apple container version: ${version_outputplaceholder"
    fi

    major="${BASH_REMATCH[1]placeholder"
    minor="${BASH_REMATCH[2]placeholder"
    if (( major < 1 || (major == 1 && minor < 1) )); then
        die "Apple container 1.1.0 or newer is required; found ${version_outputplaceholder."
    fi
placeholder

system_is_running() {
    container system status >/dev/null 2>&1
placeholder

start_system() {
    if ! system_is_running; then
        info "Starting Apple container services..."
        container system start --enable-kernel-install
    fi
placeholder

list_resource_ids() {
    case "$1" in
        container) container list --all --quiet ;;
        network) container network list --quiet ;;
        volume) container volume list --quiet ;;
        *) die "Unknown resource type: $1" ;;
    esac
placeholder

resource_exists() {
    local resource_type=$1
    local resource_name=$2
    local output line

    if ! output="$(list_resource_ids "${resource_typeplaceholder")"; then
        die "Failed to list Apple container ${resource_typeplaceholder resources."
    fi

    while IFS= read -r line; do
        if [[ "${lineplaceholder" == "${resource_nameplaceholder" ]]; then
            return 0
        fi
    done <<<"${outputplaceholder"

    return 1
placeholder

inspect_resource() {
    case "$1" in
        container) container inspect "$2" ;;
        network) container network inspect "$2" ;;
        volume) container volume inspect "$2" ;;
        *) die "Unknown resource type: $1" ;;
    esac
placeholder

assert_resource_owned() {
    local resource_type=$1
    local resource_name=$2
    local inspection compact

    inspection="$(inspect_resource "${resource_typeplaceholder" "${resource_nameplaceholder" | \
        plutil -extract 0.configuration.labels json -o - -)" || \
        die "Failed to inspect ${resource_typeplaceholder ${resource_nameplaceholder."
    compact="$(printf '%s' "${inspectionplaceholder" | tr -d '[:space:]')"
    if [[ "${compactplaceholder" != *"\"${STACK_LABEL_KEYplaceholder\":\"${STACK_LABEL_VALUEplaceholder\""* ]]; then
        die "Refusing to manage existing ${resource_typeplaceholder '${resource_nameplaceholder' because it is not owned by this stack."
    fi
placeholder

preflight_stack_ownership() {
    local resource_name

    for resource_name in "${APP_CONTAINERplaceholder" "${REDIS_CONTAINERplaceholder" "${POSTGRES_CONTAINERplaceholder"; do
        if resource_exists container "${resource_nameplaceholder"; then
            assert_resource_owned container "${resource_nameplaceholder"
        fi
    done
    if resource_exists network "${NETWORK_NAMEplaceholder"; then
        assert_resource_owned network "${NETWORK_NAMEplaceholder"
    fi
    for resource_name in "${APP_VOLUMEplaceholder" "${REDIS_VOLUMEplaceholder" "${POSTGRES_VOLUMEplaceholder"; do
        if resource_exists volume "${resource_nameplaceholder"; then
            assert_resource_owned volume "${resource_nameplaceholder"
        fi
    done
placeholder

ensure_network() {
    if resource_exists network "${NETWORK_NAMEplaceholder"; then
        assert_resource_owned network "${NETWORK_NAMEplaceholder"
        return
    fi

    info "Creating network ${NETWORK_NAMEplaceholder..."
    container network create \
        --label "${STACK_LABEL_KEYplaceholder=${STACK_LABEL_VALUEplaceholder" \
        "${NETWORK_NAMEplaceholder" >/dev/null
placeholder

ensure_volume() {
    local volume_name=$1

    if resource_exists volume "${volume_nameplaceholder"; then
        assert_resource_owned volume "${volume_nameplaceholder"
        return
    fi

    info "Creating volume ${volume_nameplaceholder..."
    container volume create \
        --label "${STACK_LABEL_KEYplaceholder=${STACK_LABEL_VALUEplaceholder" \
        "${volume_nameplaceholder" >/dev/null
placeholder

ensure_image_available() {
    local image=$1

    if container image inspect "${imageplaceholder" >/dev/null 2>&1; then
        return
    fi
    info "Pulling ${imageplaceholder..."
    container image pull --platform "${PLATFORMplaceholder" "${imageplaceholder"
placeholder

container_is_running() {
    local container_name=$1
    local output line

    output="$(container list --quiet)" || die "Failed to list running Apple containers."
    while IFS= read -r line; do
        if [[ "${lineplaceholder" == "${container_nameplaceholder" ]]; then
            return 0
        fi
    done <<<"${outputplaceholder"

    return 1
placeholder

ensure_system() {
    require_container_version
    require_command curl
    start_system
placeholder

container_ipv4_address() {
    local container_name=$1
    local address

    address="$(container inspect "${container_nameplaceholder" | \
        plutil -extract 0.status.networks.0.ipv4Address raw -o - -)" || \
        die "Unable to read the network address for ${container_nameplaceholder."
    address="${address%%/*placeholder"
    [[ "${addressplaceholder" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]] || \
        die "Apple container returned an invalid IPv4 address for ${container_nameplaceholder: ${addressplaceholder"
    printf '%s\n' "${addressplaceholder"
placeholder

read_env_value() {
    local key=$1
    local fallback=${2-placeholder

    awk -v wanted="${keyplaceholder" -v fallback="${fallbackplaceholder" '
        BEGIN { found = 0 placeholder
        /^[[:space:]]*#/ || /^[[:space:]]*$/ { next placeholder
        {
            separator = index($0, "=")
            if (separator == 0) { next placeholder
            key = substr($0, 1, separator - 1)
            if (key == wanted) {
                value = substr($0, separator + 1)
                sub(/\r$/, "", value)
                found = 1
            placeholder
        placeholder
        END {
            if (found) { print value placeholder
            else { print fallback placeholder
        placeholder
    ' "${ENV_FILEplaceholder"
placeholder

replace_env_value() {
    local key=$1
    local value=$2
    local target_file=${3:-${ENV_FILEplaceholderplaceholder
    local temp_file="${target_fileplaceholder.tmp.$$"

    awk -v wanted="${keyplaceholder" -v replacement="${valueplaceholder" '
        BEGIN { replaced = 0 placeholder
        {
            separator = index($0, "=")
            key = separator == 0 ? "" : substr($0, 1, separator - 1)
            if (key == wanted) {
                if (!replaced) { print wanted "=" replacement placeholder
                replaced = 1
                next
            placeholder
            print
        placeholder
        END {
            if (!replaced) { print wanted "=" replacement placeholder
        placeholder
    ' "${target_fileplaceholder" >"${temp_fileplaceholder"
    chmod 600 "${temp_fileplaceholder"
    mv "${temp_fileplaceholder" "${target_fileplaceholder"
placeholder

generate_secret() {
    openssl rand -hex 32
placeholder

cmd_init() {
    local env_dir temp_file postgres_secret jwt_secret totp_secret

    require_command openssl

    if [[ -e "${ENV_FILEplaceholder" ]]; then
        die "Environment file already exists: ${ENV_FILEplaceholder"
    fi

    postgres_secret="$(generate_secret)" || die "Failed to generate PostgreSQL password."
    jwt_secret="$(generate_secret)" || die "Failed to generate JWT secret."
    totp_secret="$(generate_secret)" || die "Failed to generate TOTP encryption key."
    [[ -n "${postgres_secretplaceholder" && -n "${jwt_secretplaceholder" && -n "${totp_secretplaceholder" ]] || \
        die "Secret generation returned an empty value."

    env_dir="$(dirname "${ENV_FILEplaceholder")"
    temp_file="${ENV_FILEplaceholder.init.tmp.$$"
    mkdir -p "${env_dirplaceholder"
    cp "${SCRIPT_DIRplaceholder/.env.example" "${temp_fileplaceholder"
    chmod 600 "${temp_fileplaceholder"
    replace_env_value POSTGRES_PASSWORD "${postgres_secretplaceholder" "${temp_fileplaceholder"
    replace_env_value JWT_SECRET "${jwt_secretplaceholder" "${temp_fileplaceholder"
    replace_env_value TOTP_ENCRYPTION_KEY "${totp_secretplaceholder" "${temp_fileplaceholder"
    mv "${temp_fileplaceholder" "${ENV_FILEplaceholder"

    info "Created ${ENV_FILEplaceholder with generated secrets."
    info "Review the file, then run: SUB2API_ENV_FILE='${ENV_FILEplaceholder' ${SCRIPT_DIRplaceholder/apple-container.sh up"
placeholder

validate_port() {
    local port=$1
    local decimal_port

    [[ "${portplaceholder" =~ ^[0-9]+$ ]] || die "SERVER_PORT must be numeric: ${portplaceholder"
    decimal_port=$((10#${portplaceholder))
    (( decimal_port >= 1025 && decimal_port <= 65535 )) || \
        die "SERVER_PORT must be between 1025 and 65535 for Apple container port forwarding."
placeholder

validate_ipv4_address() {
    local address=$1
    local first second third fourth extra octet

    IFS=. read -r first second third fourth extra <<<"${addressplaceholder"
    [[ -n "${firstplaceholder" && -n "${secondplaceholder" && -n "${thirdplaceholder" && -n "${fourthplaceholder" && -z "${extraplaceholder" ]] || \
        die "BIND_HOST must be a valid IPv4 address: ${addressplaceholder"
    for octet in "${firstplaceholder" "${secondplaceholder" "${thirdplaceholder" "${fourthplaceholder"; do
        [[ "${octetplaceholder" =~ ^[0-9]+$ ]] || die "BIND_HOST must be a valid IPv4 address: ${addressplaceholder"
        (( 10#${octetplaceholder <= 255 )) || die "BIND_HOST must be a valid IPv4 address: ${addressplaceholder"
    done
placeholder

validate_env_file_security() {
    local owner mode permissions

    [[ -f "${ENV_FILEplaceholder" ]] || die "Environment file not found: ${ENV_FILEplaceholder. Run '$0 init' first."
    owner="$(stat -f '%u' "${ENV_FILEplaceholder")" || die "Unable to read owner for ${ENV_FILEplaceholder."
    mode="$(stat -f '%Lp' "${ENV_FILEplaceholder")" || die "Unable to read permissions for ${ENV_FILEplaceholder."
    [[ "${ownerplaceholder" == "${EUIDplaceholder" ]] || die "Environment file must be owned by the current user: ${ENV_FILEplaceholder"
    [[ "${modeplaceholder" =~ ^[0-7]+$ ]] || die "Unable to parse permissions for ${ENV_FILEplaceholder: ${modeplaceholder"
    permissions=$((8#${modeplaceholder))
    (( (permissions & 077) == 0 )) || \
        die "Environment file must not be readable by group or others. Run: chmod 600 '${ENV_FILEplaceholder'"
placeholder

prepare_environment() {
    validate_env_file_security

    APP_IMAGE="$(read_env_value APPLE_CONTAINER_SUB2API_IMAGE weishaw/sub2api:latest)"
    POSTGRES_IMAGE="$(read_env_value APPLE_CONTAINER_POSTGRES_IMAGE postgres:18-alpine)"
    REDIS_IMAGE="$(read_env_value APPLE_CONTAINER_REDIS_IMAGE redis:8-alpine)"
    BIND_HOST="$(read_env_value BIND_HOST 0.0.0.0)"
    HOST_PORT="$(read_env_value SERVER_PORT 8080)"
    POSTGRES_USER="$(read_env_value POSTGRES_USER sub2api)"
    POSTGRES_PASSWORD="$(read_env_value POSTGRES_PASSWORD)"
    POSTGRES_DB="$(read_env_value POSTGRES_DB sub2api)"
    REDIS_PASSWORD="$(read_env_value REDIS_PASSWORD)"
    TZ_VALUE="$(read_env_value TZ Asia/Shanghai)"

    [[ -n "${BIND_HOSTplaceholder" ]] || die "BIND_HOST must not be empty."
    validate_ipv4_address "${BIND_HOSTplaceholder"
    validate_port "${HOST_PORTplaceholder"
    if [[ "${BIND_HOSTplaceholder" == "0.0.0.0" ]]; then
        ACCESS_HOST="127.0.0.1"
    else
        ACCESS_HOST="${BIND_HOSTplaceholder"
    fi
    [[ -n "${POSTGRES_USERplaceholder" ]] || die "POSTGRES_USER must not be empty."
    [[ -n "${POSTGRES_DBplaceholder" ]] || die "POSTGRES_DB must not be empty."
    if [[ -z "${POSTGRES_PASSWORDplaceholder" || "${POSTGRES_PASSWORDplaceholder" == "change_this_secure_password" ]]; then
        die "Set a secure POSTGRES_PASSWORD in ${ENV_FILEplaceholder."
    fi

    TEMP_DIR="$(mktemp -d "${TMPDIR:-/tmpplaceholder/sub2api-apple.XXXXXX")"
    APP_ENV_FILE="${TEMP_DIRplaceholder/app.env"
    POSTGRES_ENV_FILE="${TEMP_DIRplaceholder/postgres.env"
    POSTGRES_PROBE_ENV_FILE="${TEMP_DIRplaceholder/postgres-probe.env"
    REDIS_ENV_FILE="${TEMP_DIRplaceholder/redis.env"

    cat >"${POSTGRES_ENV_FILEplaceholder" <<EOF
POSTGRES_USER=${POSTGRES_USERplaceholder
POSTGRES_PASSWORD=${POSTGRES_PASSWORDplaceholder
POSTGRES_DB=${POSTGRES_DBplaceholder
TZ=${TZ_VALUEplaceholder
EOF

    cat >"${POSTGRES_PROBE_ENV_FILEplaceholder" <<EOF
PGPASSWORD=${POSTGRES_PASSWORDplaceholder
EOF

    cat >"${REDIS_ENV_FILEplaceholder" <<EOF
REDIS_PASSWORD=${REDIS_PASSWORDplaceholder
TZ=${TZ_VALUEplaceholder
EOF
    if [[ -n "${REDIS_PASSWORDplaceholder" ]]; then
        printf 'REDISCLI_AUTH=%s\n' "${REDIS_PASSWORDplaceholder" >>"${REDIS_ENV_FILEplaceholder"
    fi

    chmod 600 "${POSTGRES_ENV_FILEplaceholder" "${POSTGRES_PROBE_ENV_FILEplaceholder" "${REDIS_ENV_FILEplaceholder"
placeholder

prepare_app_environment() {
    [[ -n "${POSTGRES_ADDRESSplaceholder" && -n "${REDIS_ADDRESSplaceholder" ]] || \
        die "Dependency network addresses are not available."

    cp "${ENV_FILEplaceholder" "${APP_ENV_FILEplaceholder"
    cat >>"${APP_ENV_FILEplaceholder" <<EOF

AUTO_SETUP=true
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
DATABASE_HOST=${POSTGRES_ADDRESSplaceholder
DATABASE_PORT=5432
DATABASE_USER=${POSTGRES_USERplaceholder
DATABASE_PASSWORD=${POSTGRES_PASSWORDplaceholder
DATABASE_DBNAME=${POSTGRES_DBplaceholder
DATABASE_SSLMODE=disable
REDIS_HOST=${REDIS_ADDRESSplaceholder
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORDplaceholder
DATA_DIR=/app/storage/data
EOF
    chmod 600 "${APP_ENV_FILEplaceholder"
placeholder

create_postgres_container() {
    info "Creating PostgreSQL container..."
    container create \
        --name "${POSTGRES_CONTAINERplaceholder" \
        --label "${STACK_LABEL_KEYplaceholder=${STACK_LABEL_VALUEplaceholder" \
        --network "${NETWORK_NAMEplaceholder" \
        --platform "${PLATFORMplaceholder" \
        --ulimit nofile=100000:100000 \
        --env-file "${POSTGRES_ENV_FILEplaceholder" \
        --volume "${POSTGRES_VOLUMEplaceholder:/var/lib/postgresql" \
        "${POSTGRES_IMAGEplaceholder" >/dev/null
placeholder

create_redis_container() {
    info "Creating Redis container..."
    container create \
        --name "${REDIS_CONTAINERplaceholder" \
        --label "${STACK_LABEL_KEYplaceholder=${STACK_LABEL_VALUEplaceholder" \
        --network "${NETWORK_NAMEplaceholder" \
        --platform "${PLATFORMplaceholder" \
        --ulimit nofile=100000:100000 \
        --env-file "${REDIS_ENV_FILEplaceholder" \
        --volume "${REDIS_VOLUMEplaceholder:/var/lib/redis" \
        "${REDIS_IMAGEplaceholder" \
        sh -c 'set -e; mkdir -p /var/lib/redis/data; chown redis:redis /var/lib/redis/data; exec /usr/local/bin/docker-entrypoint.sh redis-server --dir /var/lib/redis/data --save 60 1 --appendonly yes --appendfsync everysec ${REDIS_PASSWORD:+--requirepass "$REDIS_PASSWORD"placeholder' \
        >/dev/null
placeholder

create_app_container() {
    info "Creating Sub2API container..."
    container create \
        --name "${APP_CONTAINERplaceholder" \
        --label "${STACK_LABEL_KEYplaceholder=${STACK_LABEL_VALUEplaceholder" \
        --network "${NETWORK_NAMEplaceholder" \
        --platform "${PLATFORMplaceholder" \
        --ulimit nofile=100000:100000 \
        --publish "${BIND_HOSTplaceholder:${HOST_PORTplaceholder:8080/tcp" \
        --env-file "${APP_ENV_FILEplaceholder" \
        --volume "${APP_VOLUMEplaceholder:/app/storage" \
        --entrypoint /bin/sh \
        "${APP_IMAGEplaceholder" \
        -c 'set -e; mkdir -p "$DATA_DIR"; chown -R sub2api:sub2api "$DATA_DIR"; exec su-exec sub2api /app/sub2api' \
        >/dev/null
placeholder

ensure_container() {
    local container_name=$1
    local create_function=$2

    if resource_exists container "${container_nameplaceholder"; then
        assert_resource_owned container "${container_nameplaceholder"
        return
    fi

    "${create_functionplaceholder"
placeholder

start_container_if_needed() {
    local container_name=$1

    if container_is_running "${container_nameplaceholder"; then
        return
    fi

    info "Starting ${container_nameplaceholder..."
    container start "${container_nameplaceholder" >/dev/null
placeholder

stop_container_if_running() {
    local container_name=$1

    if ! resource_exists container "${container_nameplaceholder"; then
        return
    fi
    assert_resource_owned container "${container_nameplaceholder"
    if container_is_running "${container_nameplaceholder"; then
        info "Stopping ${container_nameplaceholder..."
        container stop --time 30 "${container_nameplaceholder" >/dev/null
    fi
placeholder

delete_container_if_present() {
    local container_name=$1

    if ! resource_exists container "${container_nameplaceholder"; then
        return
    fi
    assert_resource_owned container "${container_nameplaceholder"
    if container_is_running "${container_nameplaceholder"; then
        container stop --time 30 "${container_nameplaceholder" >/dev/null
    fi
    info "Deleting ${container_nameplaceholder..."
    container delete "${container_nameplaceholder" >/dev/null
placeholder

wait_for_probe() {
    local description=$1
    local attempts=$2
    shift 2

    local attempt
    for ((attempt = 1; attempt <= attempts; attempt++)); do
        if "$@" >/dev/null 2>&1; then
            info "${descriptionplaceholder is ready."
            return 0
        fi
        sleep 1
    done

    return 1
placeholder

probe_postgres() {
    container exec --env-file "${POSTGRES_PROBE_ENV_FILEplaceholder" \
        "${POSTGRES_CONTAINERplaceholder" \
        psql -h 127.0.0.1 -U "${POSTGRES_USERplaceholder" -d "${POSTGRES_DBplaceholder" \
        -v ON_ERROR_STOP=1 -tAc 'SELECT 1'
placeholder

probe_redis() {
    container exec --env-file "${REDIS_ENV_FILEplaceholder" \
        "${REDIS_CONTAINERplaceholder" \
        redis-cli ping
placeholder

probe_app() {
    container exec "${APP_CONTAINERplaceholder" \
        wget -q -T 5 -O /dev/null http://localhost:8080/health
placeholder

probe_host_app() {
    curl --fail --silent --show-error --max-time 5 \
        "http://${ACCESS_HOSTplaceholder:${HOST_PORTplaceholder/health"
placeholder

show_failure_logs() {
    local container_name=$1

    warn "Last logs from ${container_nameplaceholder:"
    container logs -n 50 "${container_nameplaceholder" >&2 || true
placeholder

start_dependencies() {
    start_container_if_needed "${POSTGRES_CONTAINERplaceholder"
    if ! wait_for_probe "PostgreSQL" 90 probe_postgres; then
        show_failure_logs "${POSTGRES_CONTAINERplaceholder"
        die "PostgreSQL did not become ready."
    fi

    start_container_if_needed "${REDIS_CONTAINERplaceholder"
    if ! wait_for_probe "Redis" 60 probe_redis; then
        show_failure_logs "${REDIS_CONTAINERplaceholder"
        die "Redis did not become ready."
    fi
placeholder

start_app() {
    start_container_if_needed "${APP_CONTAINERplaceholder"
    if ! wait_for_probe "Sub2API" 180 probe_app; then
        show_failure_logs "${APP_CONTAINERplaceholder"
        die "Sub2API did not become ready."
    fi
    if ! wait_for_probe "Sub2API host port" 15 probe_host_app; then
        die "Host port forwarding failed. In System Settings > Privacy & Security > Local Network, allow container-runtime-linux; restart Apple container services; then run 'apple-container.sh up' again."
    fi
placeholder

cmd_up() {
    local recreate=false

    if [[ $# -gt 1 || ($# -eq 1 && "${1-placeholder" != "--recreate") ]]; then
        usage
        exit 2
    fi
    if [[ $# -eq 1 ]]; then
        recreate=true
    fi

    ensure_system
    prepare_environment
    preflight_stack_ownership
    ensure_network
    ensure_volume "${APP_VOLUMEplaceholder"
    ensure_volume "${POSTGRES_VOLUMEplaceholder"
    ensure_volume "${REDIS_VOLUMEplaceholder"
    ensure_image_available "${APP_IMAGEplaceholder"
    ensure_image_available "${POSTGRES_IMAGEplaceholder"
    ensure_image_available "${REDIS_IMAGEplaceholder"

    if [[ "${recreateplaceholder" == true ]]; then
        delete_container_if_present "${APP_CONTAINERplaceholder"
        delete_container_if_present "${REDIS_CONTAINERplaceholder"
        delete_container_if_present "${POSTGRES_CONTAINERplaceholder"
    fi

    ensure_container "${POSTGRES_CONTAINERplaceholder" create_postgres_container
    ensure_container "${REDIS_CONTAINERplaceholder" create_redis_container
    start_dependencies
    POSTGRES_ADDRESS="$(container_ipv4_address "${POSTGRES_CONTAINERplaceholder")"
    REDIS_ADDRESS="$(container_ipv4_address "${REDIS_CONTAINERplaceholder")"
    prepare_app_environment
    # The dependency IPs may change whenever their lightweight VMs restart.
    delete_container_if_present "${APP_CONTAINERplaceholder"
    create_app_container
    start_app

    info "Sub2API is available at http://${ACCESS_HOSTplaceholder:${HOST_PORTplaceholder"
placeholder

cmd_down() {
    require_container_version
    if ! system_is_running; then
        info "Apple container services are already stopped."
        return
    fi
    preflight_stack_ownership
    stop_container_if_running "${APP_CONTAINERplaceholder"
    stop_container_if_running "${REDIS_CONTAINERplaceholder"
    stop_container_if_running "${POSTGRES_CONTAINERplaceholder"
    info "Sub2API stack stopped; persistent volumes were preserved."
placeholder

cmd_restart() {
    cmd_down
    cmd_up
placeholder

print_container_status() {
    local service=$1
    local container_name=$2

    if ! resource_exists container "${container_nameplaceholder"; then
        printf '%-12s %s\n' "${serviceplaceholder" "missing"
    elif container_is_running "${container_nameplaceholder"; then
        printf '%-12s %s\n' "${serviceplaceholder" "running"
    else
        printf '%-12s %s\n' "${serviceplaceholder" "stopped"
    fi
placeholder

cmd_status() {
    local failed=0

    require_container_version
    if ! system_is_running; then
        printf '%-12s %s\n' "system" "stopped"
        return 1
    fi

    printf '%-12s %s\n' "system" "running"
    preflight_stack_ownership
    print_container_status app "${APP_CONTAINERplaceholder"
    print_container_status postgres "${POSTGRES_CONTAINERplaceholder"
    print_container_status redis "${REDIS_CONTAINERplaceholder"

    if [[ -f "${ENV_FILEplaceholder" ]]; then
        prepare_environment
        if container_is_running "${POSTGRES_CONTAINERplaceholder" && probe_postgres >/dev/null 2>&1; then
            printf '%-12s %s\n' "postgres" "healthy"
        else
            printf '%-12s %s\n' "postgres" "unhealthy"
            failed=1
        fi
        if container_is_running "${REDIS_CONTAINERplaceholder" && probe_redis >/dev/null 2>&1; then
            printf '%-12s %s\n' "redis" "healthy"
        else
            printf '%-12s %s\n' "redis" "unhealthy"
            failed=1
        fi
        if container_is_running "${APP_CONTAINERplaceholder" && probe_app >/dev/null 2>&1; then
            printf '%-12s %s\n' "app" "healthy"
        else
            printf '%-12s %s\n' "app" "unhealthy"
            failed=1
        fi
        if container_is_running "${APP_CONTAINERplaceholder" && probe_host_app >/dev/null 2>&1; then
            printf '%-12s %s\n' "host-port" "healthy"
        else
            printf '%-12s %s\n' "host-port" "unhealthy"
            failed=1
        fi
    else
        warn "Health probes require ${ENV_FILEplaceholder."
        failed=1
    fi

    return "${failedplaceholder"
placeholder

cmd_logs() {
    local service=${1-placeholder
    local follow=${2-placeholder
    local container_name

    [[ $# -ge 1 && $# -le 2 ]] || { usage; exit 2; placeholder
    if [[ -n "${followplaceholder" && "${followplaceholder" != "-f" && "${followplaceholder" != "--follow" ]]; then
        usage
        exit 2
    fi

    case "${serviceplaceholder" in
        app|sub2api) container_name="${APP_CONTAINERplaceholder" ;;
        postgres) container_name="${POSTGRES_CONTAINERplaceholder" ;;
        redis) container_name="${REDIS_CONTAINERplaceholder" ;;
        *) die "Unknown service '${serviceplaceholder'. Use app, postgres, or redis." ;;
    esac

    require_container_version
    system_is_running || die "Apple container services are not running."
    resource_exists container "${container_nameplaceholder" || die "Container not found: ${container_nameplaceholder"
    assert_resource_owned container "${container_nameplaceholder"
    if [[ -n "${followplaceholder" ]]; then
        container logs --follow "${container_nameplaceholder"
    else
        container logs "${container_nameplaceholder"
    fi
placeholder

cmd_pull() {
    ensure_system
    prepare_environment
    info "Pulling ${APP_IMAGEplaceholder..."
    container image pull --platform "${PLATFORMplaceholder" "${APP_IMAGEplaceholder"
    info "Pulling ${POSTGRES_IMAGEplaceholder..."
    container image pull --platform "${PLATFORMplaceholder" "${POSTGRES_IMAGEplaceholder"
    info "Pulling ${REDIS_IMAGEplaceholder..."
    container image pull --platform "${PLATFORMplaceholder" "${REDIS_IMAGEplaceholder"
placeholder

confirm_destroy() {
    local include_volumes=$1
    local answer

    if [[ "${include_volumesplaceholder" == true ]]; then
        printf 'Delete the Sub2API stack and all persistent data? [y/N] '
    else
        printf 'Delete the Sub2API containers and network, preserving volumes? [y/N] '
    fi
    read -r answer
    [[ "${answerplaceholder" == "y" || "${answerplaceholder" == "Y" ]]
placeholder

delete_volume_if_present() {
    local volume_name=$1

    if resource_exists volume "${volume_nameplaceholder"; then
        assert_resource_owned volume "${volume_nameplaceholder"
        info "Deleting volume ${volume_nameplaceholder..."
        container volume delete "${volume_nameplaceholder" >/dev/null
    fi
placeholder

cmd_destroy() {
    local include_volumes=false
    local assume_yes=false
    local argument

    for argument in "$@"; do
        case "${argumentplaceholder" in
            --volumes) include_volumes=true ;;
            --yes) assume_yes=true ;;
            *) usage; exit 2 ;;
        esac
    done

    require_container_version
    start_system
    preflight_stack_ownership
    if [[ "${assume_yesplaceholder" != true ]] && ! confirm_destroy "${include_volumesplaceholder"; then
        info "Cancelled."
        return
    fi

    delete_container_if_present "${APP_CONTAINERplaceholder"
    delete_container_if_present "${REDIS_CONTAINERplaceholder"
    delete_container_if_present "${POSTGRES_CONTAINERplaceholder"

    if resource_exists network "${NETWORK_NAMEplaceholder"; then
        assert_resource_owned network "${NETWORK_NAMEplaceholder"
        info "Deleting network ${NETWORK_NAMEplaceholder..."
        container network delete "${NETWORK_NAMEplaceholder" >/dev/null
    fi

    if [[ "${include_volumesplaceholder" == true ]]; then
        delete_volume_if_present "${APP_VOLUMEplaceholder"
        delete_volume_if_present "${REDIS_VOLUMEplaceholder"
        delete_volume_if_present "${POSTGRES_VOLUMEplaceholder"
        info "Sub2API stack and persistent data deleted."
    else
        info "Sub2API stack deleted; persistent volumes were preserved."
    fi
placeholder

main() {
    local command=${1-placeholder
    if [[ $# -gt 0 ]]; then
        shift
    fi

    case "${commandplaceholder" in
        init)
            [[ $# -eq 0 ]] || { usage; exit 2; placeholder
            acquire_lock
            cmd_init
            ;;
        up)
            acquire_lock
            cmd_up "$@"
            ;;
        down)
            [[ $# -eq 0 ]] || { usage; exit 2; placeholder
            acquire_lock
            cmd_down
            ;;
        restart)
            [[ $# -eq 0 ]] || { usage; exit 2; placeholder
            acquire_lock
            cmd_restart
            ;;
        status)
            [[ $# -eq 0 ]] || { usage; exit 2; placeholder
            trap cleanup EXIT
            cmd_status
            ;;
        logs)
            cmd_logs "$@"
            ;;
        pull)
            [[ $# -eq 0 ]] || { usage; exit 2; placeholder
            acquire_lock
            cmd_pull
            ;;
        destroy)
            acquire_lock
            cmd_destroy "$@"
            ;;
        help|-h|--help)
            usage
            ;;
        *)
            usage
            exit 2
            ;;
    esac
placeholder

main "$@"
