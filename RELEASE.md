# 发布指南

本文档说明如何发布 `go-istanbul-sourcemap` Go 包。

## 🚀 自动发布流程

### 方式一：通过 Git 标签触发（推荐）

1. **确保代码已准备就绪**
   ```bash
   # 运行测试
   cd go && go test -v ./...
   
   # 检查代码质量
   go vet ./...
   ```

2. **使用发布脚本**
   ```bash
   # 使用提供的脚本（推荐）
   ./scripts/release.sh v1.0.0
   ```
   
   或者手动创建标签：
   ```bash
   # 创建并推送标签
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. **GitHub Action 自动执行**
   - 运行完整测试套件
   - 构建包
   - 创建 GitHub Release
   - 触发 Go proxy 更新
   - 验证包可用性

### 方式二：手动触发 GitHub Action

1. 访问 [GitHub Actions](https://github.com/canyon-project/go-istanbul-sourcemap/actions)
2. 选择 "Release Go Package" workflow
3. 点击 "Run workflow"
4. 输入版本号（如 `v1.0.0`）
5. 点击 "Run workflow"

## 📋 版本号规范

遵循 [语义化版本](https://semver.org/lang/zh-CN/) 规范：

- **主版本号**：不兼容的 API 修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

### 版本格式示例

- `v1.0.0` - 正式版本
- `v1.0.1` - 修复版本
- `v1.1.0` - 功能版本
- `v2.0.0` - 重大更新
- `v1.0.0-beta` - 预发布版本
- `v1.0.0-alpha.1` - 内测版本

## 🔄 发布流程详解

### 1. 预发布检查

发布前会自动执行以下检查：

- ✅ 运行所有单元测试
- ✅ 执行基准测试
- ✅ 代码质量检查（go vet）
- ✅ 依赖验证
- ✅ 跨平台兼容性测试

### 2. 发布执行

- 📦 创建 GitHub Release
- 📝 自动生成更新日志
- 🔄 触发 Go proxy 缓存更新
- 🌐 更新 Go checksum 数据库

### 3. 发布验证

- ✅ 验证包在 Go proxy 中可用
- ✅ 测试包安装和导入
- ✅ 确认版本信息正确

## 📦 用户安装方式

发布成功后，用户可以通过以下方式安装：

```bash
# 安装最新版本
go get github.com/canyon-project/go-istanbul-sourcemap

# 安装特定版本
go get github.com/canyon-project/go-istanbul-sourcemap@v1.0.0

# 安装预发布版本
go get github.com/canyon-project/go-istanbul-sourcemap@v1.0.0-beta
```

## 🛠️ 本地开发

### 运行测试

```bash
cd go
go test -v ./...
go test -bench=. -benchmem ./...
```

### 代码质量检查

```bash
cd go
go vet ./...
golangci-lint run
```

### 构建示例

```bash
cd go
go build ./...
go run example/main.go
```

## 🔍 故障排除

### 常见问题

1. **标签已存在**
   ```bash
   # 删除本地标签
   git tag -d v1.0.0
   
   # 删除远程标签（谨慎操作）
   git push origin :refs/tags/v1.0.0
   ```

2. **测试失败**
   ```bash
   # 查看详细测试输出
   cd go && go test -v -race ./...
   ```

3. **Go proxy 更新延迟**
   - Go proxy 可能需要几分钟来更新
   - 可以手动触发：`GOPROXY=direct go get github.com/canyon-project/go-istanbul-sourcemap@v1.0.0`

4. **权限问题**
   - 确保有仓库的写权限
   - 检查 GitHub token 权限

### 查看发布状态

- [GitHub Actions](https://github.com/canyon-project/go-istanbul-sourcemap/actions)
- [GitHub Releases](https://github.com/canyon-project/go-istanbul-sourcemap/releases)
- [Go Packages](https://pkg.go.dev/github.com/canyon-project/go-istanbul-sourcemap)

## 📞 获取帮助

如果遇到发布问题：

1. 检查 [GitHub Actions 日志](https://github.com/canyon-project/go-istanbul-sourcemap/actions)
2. 查看 [Issues](https://github.com/canyon-project/go-istanbul-sourcemap/issues)
3. 联系维护者

## 🔐 安全注意事项

- 只有维护者可以创建发布
- 所有发布都会经过自动化测试
- 使用 GitHub 的安全扫描
- 遵循最小权限原则