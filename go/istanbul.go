// Package istanbul provides Istanbul coverage source map transformation functionality
package istanbul

import (
	"encoding/json"
	"fmt"
)

// Istanbul provides the main API for coverage transformation
type Istanbul struct {
	transformer *CoverageTransformer
}

// New creates a new Istanbul instance
func New() *Istanbul {
	return &Istanbul{
		transformer: NewCoverageTransformer(),
	}
}

// TransformCoverage transforms Istanbul coverage data using source maps
func (i *Istanbul) TransformCoverage(coverageData string) (string, error) {
	// Parse input JSON
	coverage, err := ParseCoverageMap([]byte(coverageData))
	if err != nil {
		return "", fmt.Errorf("failed to parse coverage data: %w", err)
	}

	// Transform coverage
	transformed, err := i.transformer.Transform(coverage)
	if err != nil {
		return "", fmt.Errorf("failed to transform coverage: %w", err)
	}

	// Convert back to JSON
	result, err := transformed.ToJSON()
	if err != nil {
		return "", fmt.Errorf("failed to serialize result: %w", err)
	}

	return string(result), nil
}

// TransformCoverageBytes transforms Istanbul coverage data from bytes
func (i *Istanbul) TransformCoverageBytes(coverageData []byte) ([]byte, error) {
	result, err := i.TransformCoverage(string(coverageData))
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

// GetVersion returns the version of this library
func (i *Istanbul) GetVersion() string {
	return "1.0.0"
}

// GetPlatform returns platform information
func (i *Istanbul) GetPlatform() string {
	return "go/pure"
}

// Package-level convenience functions

// TransformCoverageString transforms coverage data from string
func TransformCoverageString(coverageData string) (string, error) {
	istanbul := New()
	return istanbul.TransformCoverage(coverageData)
}

// TransformCoverageBytes transforms coverage data from bytes
func TransformCoverageBytes(coverageData []byte) ([]byte, error) {
	istanbul := New()
	return istanbul.TransformCoverageBytes(coverageData)
}

// ValidateCoverageData validates that the input is valid Istanbul coverage data
func ValidateCoverageData(data []byte) error {
	var coverage CoverageMap
	if err := json.Unmarshal(data, &coverage); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Basic validation
	for filePath, fc := range coverage {
		if fc == nil {
			return fmt.Errorf("null coverage data for file: %s", filePath)
		}
		if fc.Path == "" {
			return fmt.Errorf("missing path in coverage data for file: %s", filePath)
		}
		if fc.StatementMap == nil {
			return fmt.Errorf("missing statementMap for file: %s", filePath)
		}
		if fc.FnMap == nil {
			return fmt.Errorf("missing fnMap for file: %s", filePath)
		}
		if fc.BranchMap == nil {
			return fmt.Errorf("missing branchMap for file: %s", filePath)
		}
	}

	return nil
}
