#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# 检查是否在正确的目录
if [ ! -f "go/go.mod" ]; then
    print_error "请在项目根目录运行此脚本"
    exit 1
fi

# 检查是否有未提交的更改
if [ -n "$(git status --porcelain)" ]; then
    print_error "有未提交的更改，请先提交所有更改"
    git status --short
    exit 1
fi

# 获取版本号
if [ -z "$1" ]; then
    print_error "请提供版本号"
    echo "用法: $0 <version>"
    echo "示例: $0 v1.0.0"
    exit 1
fi

VERSION=$1

# 验证版本格式
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    print_error "无效的版本格式: $VERSION"
    print_info "期望格式: v1.0.0 或 v1.0.0-beta"
    exit 1
fi

print_info "准备发布版本: $VERSION"

# 检查标签是否已存在
if git tag -l | grep -q "^$VERSION$"; then
    print_error "标签 $VERSION 已存在"
    exit 1
fi

# 切换到go目录进行测试
cd go

print_info "运行测试..."
if ! go test -v ./...; then
    print_error "测试失败"
    exit 1
fi
print_success "所有测试通过"

print_info "运行代码检查..."
if ! go vet ./...; then
    print_error "代码检查失败"
    exit 1
fi
print_success "代码检查通过"

print_info "检查依赖..."
go mod tidy
if [ -n "$(git status --porcelain go.mod go.sum)" ]; then
    print_warning "go.mod 或 go.sum 有更改，正在提交..."
    git add go.mod go.sum
    git commit -m "chore: update go.mod and go.sum for release $VERSION"
fi

# 回到根目录
cd ..

print_info "创建标签..."
git tag -a "$VERSION" -m "Release $VERSION"

print_info "推送标签到远程仓库..."
git push origin "$VERSION"

print_success "版本 $VERSION 发布成功！"
print_info "GitHub Action 将自动处理剩余的发布流程"
print_info "你可以在以下地址查看发布状态:"
print_info "https://github.com/canyon-project/go-istanbul-sourcemap/actions"

print_info "发布完成后，用户可以通过以下命令安装:"
echo -e "${GREEN}go get github.com/canyon-project/go-istanbul-sourcemap@$VERSION${NC}"