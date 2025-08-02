# Istanbul Source Maps

多语言的Istanbul覆盖率数据source map转换器实现。

## 🌍 语言支持

- [Go](./go/) - 纯Go语言实现，高性能，无外部依赖

## 📦 安装

### Go 版本

```bash
go get github.com/canyon-project/go-istanbul-sourcemap
```

## 🚀 快速开始

### Go 示例

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
    
    // 转换覆盖率数据
    result, err := ist.TransformCoverage(coverageData)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("转换结果: %s\n", result)
}
```

## 📚 文档

- [Go 包文档](./go/README.md)
- [发布指南](./RELEASE.md)

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

MIT License
