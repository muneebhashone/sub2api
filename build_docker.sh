#!/bin/bash
# =============================================================================
# Sub2API Docker 镜像构建脚本
# =============================================================================
#
# 用法:
#   ./build_docker.sh [选项]
#
# 选项:
#   -t, --tag TAG          指定镜像标签 (默认: latest)
#   -r, --registry REG     指定镜像仓库地址 (默认: 无)
#   -b, --base-registry REG 指定基础镜像仓库地址 (默认: 使用 docker.io)
#   -p, --push             构建后推送镜像到仓库
#   --no-cache             不使用 Docker 构建缓存
#   -h, --help             显示帮助信息
#
# 示例:
#   ./build_docker.sh                           # 构建 sub2api:latest
#   ./build_docker.sh -t v1.0.0                 # 构建 sub2api:v1.0.0
#   ./build_docker.sh -r ghcr.io/user -p        # 构建并推送到 GitHub Container Registry
#   ./build_docker.sh -t v1.0.0 --no-cache      # 不使用缓存构建
#
# =============================================================================

set -e  # 遇到错误立即退出

# =============================================================================
# 颜色定义
# =============================================================================
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# =============================================================================
# 默认配置
# =============================================================================
IMAGE_NAME="sub2api"
TAG="latest"
REGISTRY=""
BASE_IMAGE_REGISTRY="${BASE_IMAGE_REGISTRY:-placeholder"
PUSH=false
NO_CACHE=""

# =============================================================================
# 辅助函数
# =============================================================================

# 打印带颜色的信息
info() {
    echo -e "${BLUEplaceholder[INFO]${NCplaceholder $1"
placeholder

success() {
    echo -e "${GREENplaceholder[SUCCESS]${NCplaceholder $1"
placeholder

warn() {
    echo -e "${YELLOWplaceholder[WARN]${NCplaceholder $1"
placeholder

error() {
    echo -e "${REDplaceholder[ERROR]${NCplaceholder $1"
    exit 1
placeholder

# 显示帮助信息
show_help() {
    head -30 "$0" | tail -25
    exit 0
placeholder

# =============================================================================
# 参数解析
# =============================================================================
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        -b|--base-registry)
            BASE_IMAGE_REGISTRY="$2"
            shift 2
            ;;
        -p|--push)
            PUSH=true
            shift
            ;;
        --no-cache)
            NO_CACHE="--no-cache"
            shift
            ;;
        -h|--help)
            show_help
            ;;
        *)
            error "未知选项: $1\n使用 -h 或 --help 查看帮助"
            ;;
    esac
done

# =============================================================================
# 获取构建信息
# =============================================================================

# 获取 Git commit hash
get_commit() {
    if git rev-parse --git-dir > /dev/null 2>&1; then
        git rev-parse --short HEAD 2>/dev/null || echo "unknown"
    else
        echo "unknown"
    fi
placeholder

# 获取版本号 (优先使用 git tag)
get_version() {
    if git rev-parse --git-dir > /dev/null 2>&1; then
        # 尝试获取最近的 tag
        local tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
        if [[ -n "$tag" ]]; then
            echo "$tag"
        else
            # 没有 tag，使用 commit hash
            echo "dev-$(get_commit)"
        fi
    else
        echo "dev"
    fi
placeholder

# 获取构建日期 (ISO 8601 格式)
get_date() {
    date -u +"%Y-%m-%dT%H:%M:%SZ"
placeholder

# =============================================================================
# 主逻辑
# =============================================================================

# 切换到项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]placeholder")" && pwd)"
cd "$SCRIPT_DIR"

# 检查 Dockerfile 是否存在
if [[ ! -f "Dockerfile" ]]; then
    error "Dockerfile 不存在，请确保在项目根目录运行此脚本"
fi

# 获取构建参数
VERSION=$(get_version)
COMMIT=$(get_commit)
DATE=$(get_date)

# 解析基础镜像
NODE_IMAGE="node:24-alpine"
GOLANG_IMAGE="golang:1.25-alpine"
ALPINE_IMAGE="alpine:3.19"
if [[ -n "$BASE_IMAGE_REGISTRY" ]]; then
    BASE_IMAGE_REGISTRY="${BASE_IMAGE_REGISTRY%/placeholder"
    NODE_IMAGE="${BASE_IMAGE_REGISTRYplaceholder/library/node:24-alpine"
    GOLANG_IMAGE="${BASE_IMAGE_REGISTRYplaceholder/library/golang:1.25-alpine"
    ALPINE_IMAGE="${BASE_IMAGE_REGISTRYplaceholder/library/alpine:3.19"
fi

# 构建完整镜像名称
if [[ -n "$REGISTRY" ]]; then
    FULL_IMAGE_NAME="${REGISTRYplaceholder/${IMAGE_NAMEplaceholder:${TAGplaceholder"
else
    FULL_IMAGE_NAME="${IMAGE_NAMEplaceholder:${TAGplaceholder"
fi

# 打印构建信息
echo ""
echo "=============================================="
echo "  Sub2API Docker 镜像构建"
echo "=============================================="
info "镜像名称: ${FULL_IMAGE_NAMEplaceholder"
info "版本: ${VERSIONplaceholder"
info "Commit: ${COMMITplaceholder"
info "构建时间: ${DATEplaceholder"
if [[ -n "$BASE_IMAGE_REGISTRY" ]]; then
    info "基础镜像仓库: ${BASE_IMAGE_REGISTRYplaceholder"
fi
if [[ -n "$NO_CACHE" ]]; then
    info "缓存: 禁用"
fi
echo "=============================================="
echo ""

# 执行构建
info "开始构建 Docker 镜像..."

docker build \
    ${NO_CACHEplaceholder \
    --build-arg VERSION="${VERSIONplaceholder" \
    --build-arg COMMIT="${COMMITplaceholder" \
    --build-arg DATE="${DATEplaceholder" \
    --build-arg NODE_IMAGE="${NODE_IMAGEplaceholder" \
    --build-arg GOLANG_IMAGE="${GOLANG_IMAGEplaceholder" \
    --build-arg ALPINE_IMAGE="${ALPINE_IMAGEplaceholder" \
    -t "${FULL_IMAGE_NAMEplaceholder" \
    -f Dockerfile \
    .

# 检查构建结果
if [[ $? -eq 0 ]]; then
    success "镜像构建成功: ${FULL_IMAGE_NAMEplaceholder"
else
    error "镜像构建失败"
fi

# 如果指定了 latest 以外的 tag，同时打上 latest 标签
if [[ "$TAG" != "latest" && -z "$REGISTRY" ]]; then
    info "同时标记为 latest..."
    docker tag "${FULL_IMAGE_NAMEplaceholder" "${IMAGE_NAMEplaceholder:latest"
fi

# 推送镜像
if [[ "$PUSH" == true ]]; then
    if [[ -z "$REGISTRY" ]]; then
        warn "未指定镜像仓库 (-r)，跳过推送"
    else
        info "推送镜像到仓库..."
        docker push "${FULL_IMAGE_NAMEplaceholder"

        # 如果不是 latest，也推送 latest 标签
        if [[ "$TAG" != "latest" ]]; then
            LATEST_IMAGE="${REGISTRYplaceholder/${IMAGE_NAMEplaceholder:latest"
            docker tag "${FULL_IMAGE_NAMEplaceholder" "${LATEST_IMAGEplaceholder"
            docker push "${LATEST_IMAGEplaceholder"
        fi

        success "镜像推送成功"
    fi
fi

# 显示镜像信息
echo ""
echo "=============================================="
info "构建完成！"
echo "=============================================="
echo ""
info "运行容器示例:"
echo "  docker run -d \\"
echo "    --name sub2api \\"
echo "    -p 8080:8080 \\"
echo "    -e AUTO_SETUP=true \\"
echo "    -e DATABASE_HOST=your-db-host \\"
echo "    -e DATABASE_PASSWORD=your-password \\"
echo "    -e REDIS_HOST=your-redis-host \\"
echo "    ${FULL_IMAGE_NAMEplaceholder"
echo ""
info "查看镜像大小:"
docker images "${IMAGE_NAMEplaceholder" --format "table {{.Repositoryplaceholderplaceholder\t{{.Tagplaceholderplaceholder\t{{.Sizeplaceholderplaceholder"
echo ""
