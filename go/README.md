# Go Istanbul Sourcemap

纯Go语言实现的Istanbul覆盖率数据source map转换器，无需依赖外部动态库。

## ✨ 特性

- 🚀 **纯Go实现** - 无需CGO或外部依赖
- ⚡ **高性能** - 优化的算法和数据结构
- 🔄 **完整支持** - 支持语句、函数和分支覆盖率转换
- 🛠️ **易于使用** - 简洁的API接口
- 🧪 **完整测试** - 全面的测试覆盖
- 📦 **轻量级** - 最小化的依赖

## 📦 安装

```bash
go get github.com/canyon-project/go-istanbul-sourcemap
```

## 🚀 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    
    istanbul "github.com/canyon-project/go-istanbul-sourcemap"
)

func main() {
    // 创建Istanbul实例
    ist := istanbul.New()
    
    // Istanbul覆盖率数据
    coverageData := `{
        "dist/app.js": {
            "path": "dist/app.js",
            "statementMap": {
                "0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 25}}
            },
            "fnMap": {},
            "branchMap": {},
            "s": {"0": 1},
            "f": {},
            "b": {},
            "inputSourceMap": {
                "version": 3,
                "sources": ["src/app.ts"],
                "names": [],
                "mappings": "AAAA,SAASA"
            }
        }
    }`
    
    // 转换覆盖率数据
    result, err := ist.TransformCoverage(coverageData)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("转换结果: %s\n", result)
}
```

### 使用便捷函数

```go
package main

import (
    "fmt"
    "log"
    
    istanbul "github.com/canyon-project/go-istanbul-sourcemap"
)

func main() {
    // 直接使用包级别函数
    result, err := istanbul.TransformCoverageString(coverageData)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("结果: %s\n", result)
}
```

## 📋 API 参考

### Istanbul 类型

#### New() *Istanbul
创建新的Istanbul实例。

#### (*Istanbul) TransformCoverage(coverageData string) (string, error)
转换Istanbul覆盖率数据，应用source map映射。

**参数:**
- `coverageData`: JSON格式的Istanbul覆盖率数据

**返回:**
- `string`: 转换后的覆盖率数据（JSON格式）
- `error`: 错误信息（如果有）

#### (*Istanbul) TransformCoverageBytes(coverageData []byte) ([]byte, error)
转换Istanbul覆盖率数据（字节版本）。

#### (*Istanbul) GetVersion() string
获取库版本号。

#### (*Istanbul) GetPlatform() string
获取平台信息。

### 包级别函数

#### TransformCoverageString(coverageData string) (string, error)
便捷函数，等同于创建实例后调用TransformCoverage方法。

#### TransformCoverageBytes(coverageData []byte) ([]byte, error)
便捷函数，处理字节数据。

#### ValidateCoverageData(data []byte) error
验证Istanbul覆盖率数据格式是否正确。

## 🏗️ 数据结构

### 主要类型

```go
// 覆盖率映射
type CoverageMap map[string]*FileCoverage

// 文件覆盖率数据
type FileCoverage struct {
    Path           string                    `json:"path"`
    StatementMap   map[string]Location       `json:"statementMap"`
    FnMap          map[string]FunctionMeta   `json:"fnMap"`
    BranchMap      map[string]BranchMeta     `json:"branchMap"`
    S              map[string]int            `json:"s"`              // 语句命中次数
    F              map[string]int            `json:"f"`              // 函数命中次数
    B              map[string][]int          `json:"b"`              // 分支命中次数
    InputSourceMap *SourceMap               `json:"inputSourceMap,omitempty"`
}

// 位置信息
type Location struct {
    Start Position `json:"start"`
    End   Position `json:"end"`
}

type Position struct {
    Line   int `json:"line"`
    Column int `json:"column"`
}
```

## 🧪 运行示例

```bash
# 克隆仓库
git clone https://github.com/canyon-project/go-istanbul-sourcemap.git
cd go-istanbul-sourcemap

# 运行示例
go run example/main.go

# 运行测试
go test -v

# 运行基准测试
go test -bench=.
```

## 🔧 高级用法

### 自定义转换器

```go
package main

import (
    istanbul "github.com/canyon-project/go-istanbul-sourcemap"
)

func main() {
    // 创建自定义转换器
    transformer := istanbul.NewCoverageTransformer()
    
    // 解析覆盖率数据
    coverage, err := istanbul.ParseCoverageMap([]byte(coverageData))
    if err != nil {
        panic(err)
    }
    
    // 执行转换
    result, err := transformer.Transform(coverage)
    if err != nil {
        panic(err)
    }
    
    // 转换为JSON
    jsonResult, err := result.ToJSON()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("结果: %s\n", jsonResult)
}
```

### 批量处理

```go
func processCoverageFiles(files []string) error {
    ist := istanbul.New()
    
    for _, file := range files {
        data, err := os.ReadFile(file)
        if err != nil {
            return err
        }
        
        result, err := ist.TransformCoverageBytes(data)
        if err != nil {
            return err
        }
        
        // 保存结果
        outputFile := strings.Replace(file, ".json", "_transformed.json", 1)
        if err := os.WriteFile(outputFile, result, 0644); err != nil {
            return err
        }
    }
    
    return nil
}
```

## 🎯 性能特点

- **内存效率**: 优化的数据结构，最小化内存使用
- **处理速度**: 高效的算法实现
- **并发安全**: 所有公共方法都是并发安全的
- **零依赖**: 除了source map解析库外无其他依赖

## 🆚 与其他实现的对比

| 特性 | Go版本 | Rust+CGO版本 | JavaScript版本 |
|------|--------|--------------|----------------|
| 安装复杂度 | 🟢 简单 | 🔴 复杂 | 🟢 简单 |
| 运行时依赖 | 🟢 无 | 🔴 动态库 | 🟡 Node.js |
| 性能 | 🟢 高 | 🟢 高 | 🟡 中等 |
| 内存使用 | 🟢 低 | 🟢 低 | 🔴 高 |
| 跨平台 | 🟢 优秀 | 🟡 需编译 | 🟢 优秀 |
| 维护成本 | 🟢 低 | 🔴 高 | 🟡 中等 |

## 🐛 故障排除

### 常见问题

1. **解析错误**: 确保输入是有效的Istanbul覆盖率JSON格式
2. **Source Map错误**: 检查source map格式是否正确
3. **内存使用**: 对于大型项目，考虑分批处理

### 调试技巧

```go
// 启用详细错误信息
if err := istanbul.ValidateCoverageData(data); err != nil {
    fmt.Printf("数据验证失败: %v\n", err)
}

// 检查转换前后的数据
fmt.Printf("转换前: %d 个文件\n", len(originalCoverage))
fmt.Printf("转换后: %d 个文件\n", len(transformedCoverage))
```

## 🤝 贡献

欢迎提交Issue和Pull Request！

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/canyon-project/go-istanbul-sourcemap.git
cd go-istanbul-sourcemap

# 安装依赖
go mod tidy

# 运行测试
go test -v

# 运行示例
go run example/main.go
```

## 📄 许可证

MIT License

## 🔄 更新日志

### v1.0.0
- 初始版本
- 完整的Istanbul覆盖率转换功能
- 纯Go实现，无外部依赖
- 完整的测试套件
- 详细的文档和示例