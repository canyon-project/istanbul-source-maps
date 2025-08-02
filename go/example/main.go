package main

import (
	"fmt"
	"log"

	istanbul "github.com/canyon-project/go-istanbul-sourcemap"
)

func main() {
	fmt.Println("🧪 Go Istanbul Sourcemap 示例")

	// 创建Istanbul实例
	ist := istanbul.New()

	fmt.Printf("📚 版本: %s\n", ist.GetVersion())
	fmt.Printf("🖥️  平台: %s\n", ist.GetPlatform())

	// 示例1: 没有source map的覆盖率数据
	fmt.Println("\n--- 示例1: 无Source Map ---")
	simpleCoverage := `{
		"src/app.js": {
			"path": "src/app.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 25}},
				"1": {"start": {"line": 2, "column": 0}, "end": {"line": 2, "column": 20}}
			},
			"fnMap": {
				"0": {
					"name": "myFunction",
					"decl": {"start": {"line": 1, "column": 9}, "end": {"line": 1, "column": 19}},
					"loc": {"start": {"line": 1, "column": 0}, "end": {"line": 3, "column": 1}}
				}
			},
			"branchMap": {},
			"s": {"0": 5, "1": 3},
			"f": {"0": 2},
			"b": {}
		}
	}`

	result1, err := ist.TransformCoverage(simpleCoverage)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}
	fmt.Printf("✅ 转换成功 (长度: %d)\n", len(result1))

	// 示例2: 带source map的覆盖率数据
	fmt.Println("\n--- 示例2: 带Source Map ---")
	coverageWithSourceMap := `{
		"dist/bundle.js": {
			"path": "dist/bundle.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 25}},
				"1": {"start": {"line": 2, "column": 0}, "end": {"line": 2, "column": 20}},
				"2": {"start": {"line": 3, "column": 0}, "end": {"line": 3, "column": 15}}
			},
			"fnMap": {
				"0": {
					"name": "testFunction",
					"decl": {"start": {"line": 1, "column": 9}, "end": {"line": 1, "column": 21}},
					"loc": {"start": {"line": 1, "column": 0}, "end": {"line": 3, "column": 1}}
				}
			},
			"branchMap": {
				"0": {
					"type": "if",
					"loc": {"start": {"line": 2, "column": 0}, "end": {"line": 2, "column": 20}},
					"locations": [
						{"start": {"line": 2, "column": 0}, "end": {"line": 2, "column": 10}},
						{"start": {"line": 2, "column": 10}, "end": {"line": 2, "column": 20}}
					]
				}
			},
			"s": {"0": 5, "1": 3, "2": 0},
			"f": {"0": 2},
			"b": {"0": [3, 1]},
			"inputSourceMap": {
				"version": 3,
				"sources": ["src/main.ts", "src/utils.ts"],
				"names": ["testFunction", "console", "log"],
				"mappings": "AAAA,SAASA,aACP,OAAOC,QAAQC,IAAI",
				"file": "bundle.js",
				"sourceRoot": "",
				"sourcesContent": [
					"function testFunction() {\\n  if (condition) {\\n    return console.log('Hello');\\n  }\\n}",
					"export function helper() {\\n  return 'helper';\\n}"
				]
			}
		}
	}`

	result2, err := ist.TransformCoverage(coverageWithSourceMap)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}
	fmt.Printf("✅ 转换成功 (长度: %d)\n", len(result2))

	// 显示转换结果
	fmt.Println("\n--- 转换结果 ---")
	fmt.Printf("结果:\n%s\n", result2)

	// 示例3: 使用包级别函数
	fmt.Println("\n--- 示例3: 包级别函数 ---")
	result3, err := istanbul.TransformCoverageString(simpleCoverage)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}
	fmt.Printf("✅ 包级别函数转换成功 (长度: %d)\n", len(result3))

	// 示例4: 验证数据
	fmt.Println("\n--- 示例4: 数据验证 ---")
	if err := istanbul.ValidateCoverageData([]byte(simpleCoverage)); err != nil {
		fmt.Printf("❌ 验证失败: %v\n", err)
	} else {
		fmt.Println("✅ 数据验证通过")
	}

	fmt.Println("\n🎉 所有示例完成!")
}