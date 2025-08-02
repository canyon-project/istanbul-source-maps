package istanbul

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	istanbul := New()
	assert.NotNil(t, istanbul)
	assert.NotNil(t, istanbul.transformer)
}

func TestGetVersion(t *testing.T) {
	istanbul := New()
	version := istanbul.GetVersion()
	assert.Equal(t, "1.0.0", version)
}

func TestGetPlatform(t *testing.T) {
	istanbul := New()
	platform := istanbul.GetPlatform()
	assert.Equal(t, "go/pure", platform)
}

func TestTransformCoverageWithoutSourceMap(t *testing.T) {
	istanbul := New()

	// Coverage data without source map
	coverageData := `{
		"test.js": {
			"path": "test.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 10}}
			},
			"fnMap": {},
			"branchMap": {},
			"s": {"0": 1},
			"f": {},
			"b": {}
		}
	}`

	result, err := istanbul.TransformCoverage(coverageData)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Parse result to verify structure
	var resultMap CoverageMap
	err = json.Unmarshal([]byte(result), &resultMap)
	require.NoError(t, err)
	assert.Contains(t, resultMap, "test.js")
}

func TestTransformCoverageWithSourceMap(t *testing.T) {
	istanbul := New()

	// Coverage data with source map
	coverageData := `{
		"dist/bundle.js": {
			"path": "dist/bundle.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 25}}
			},
			"fnMap": {
				"0": {
					"name": "testFunction",
					"decl": {"start": {"line": 1, "column": 9}, "end": {"line": 1, "column": 21}},
					"loc": {"start": {"line": 1, "column": 0}, "end": {"line": 3, "column": 1}}
				}
			},
			"branchMap": {},
			"s": {"0": 5},
			"f": {"0": 2},
			"b": {},
			"inputSourceMap": {
				"version": 3,
				"sources": ["src/main.ts"],
				"names": ["testFunction"],
				"mappings": "AAAA,SAASA",
				"file": "bundle.js"
			}
		}
	}`

	result, err := istanbul.TransformCoverage(coverageData)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Parse result to verify structure
	var resultMap CoverageMap
	err = json.Unmarshal([]byte(result), &resultMap)
	require.NoError(t, err)

	// Should have transformed to original source
	assert.NotEmpty(t, resultMap)
}

func TestTransformCoverageBytes(t *testing.T) {
	istanbul := New()

	coverageData := []byte(`{
		"test.js": {
			"path": "test.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 10}}
			},
			"fnMap": {},
			"branchMap": {},
			"s": {"0": 1},
			"f": {},
			"b": {}
		}
	}`)

	result, err := istanbul.TransformCoverageBytes(coverageData)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestPackageLevelFunctions(t *testing.T) {
	coverageData := `{
		"test.js": {
			"path": "test.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 10}}
			},
			"fnMap": {},
			"branchMap": {},
			"s": {"0": 1},
			"f": {},
			"b": {}
		}
	}`

	// Test string function
	result, err := TransformCoverageString(coverageData)
	require.NoError(t, err)
	assert.NotEmpty(t, result)

	// Test bytes function
	resultBytes, err := TransformCoverageBytes([]byte(coverageData))
	require.NoError(t, err)
	assert.NotEmpty(t, resultBytes)
}

func TestValidateCoverageData(t *testing.T) {
	// Valid data
	validData := []byte(`{
		"test.js": {
			"path": "test.js",
			"statementMap": {},
			"fnMap": {},
			"branchMap": {},
			"s": {},
			"f": {},
			"b": {}
		}
	}`)

	err := ValidateCoverageData(validData)
	assert.NoError(t, err)

	// Invalid JSON
	invalidJSON := []byte(`{invalid json}`)
	err = ValidateCoverageData(invalidJSON)
	assert.Error(t, err)

	// Missing required fields
	missingFields := []byte(`{
		"test.js": {
			"path": "test.js"
		}
	}`)
	err = ValidateCoverageData(missingFields)
	assert.Error(t, err)
}

func TestTransformCoverageInvalidInput(t *testing.T) {
	istanbul := New()

	// Invalid JSON
	result, err := istanbul.TransformCoverage("invalid json")
	assert.Error(t, err)
	assert.Empty(t, result)
}

// Benchmark tests
func BenchmarkTransformCoverage(b *testing.B) {
	istanbul := New()
	coverageData := `{
		"test.js": {
			"path": "test.js",
			"statementMap": {
				"0": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 10}},
				"1": {"start": {"line": 2, "column": 0}, "end": {"line": 2, "column": 15}},
				"2": {"start": {"line": 3, "column": 0}, "end": {"line": 3, "column": 20}}
			},
			"fnMap": {
				"0": {
					"name": "testFn",
					"decl": {"start": {"line": 1, "column": 0}, "end": {"line": 1, "column": 10}},
					"loc": {"start": {"line": 1, "column": 0}, "end": {"line": 3, "column": 1}}
				}
			},
			"branchMap": {},
			"s": {"0": 1, "1": 1, "2": 0},
			"f": {"0": 1},
			"b": {}
		}
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := istanbul.TransformCoverage(coverageData)
		if err != nil {
			b.Fatal(err)
		}
	}
}
