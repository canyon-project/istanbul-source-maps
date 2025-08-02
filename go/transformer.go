package istanbul

import (
	"fmt"
	"strconv"
)

// CoverageTransformer transforms Istanbul coverage data using source maps
type CoverageTransformer struct {
	sourceMapTransformer *SourceMapTransformer
}

// NewCoverageTransformer creates a new coverage transformer
func NewCoverageTransformer() *CoverageTransformer {
	return &CoverageTransformer{
		sourceMapTransformer: NewSourceMapTransformer(),
	}
}

// Transform transforms coverage data using source maps
func (ct *CoverageTransformer) Transform(coverage CoverageMap) (CoverageMap, error) {
	result := make(CoverageMap)

	for filePath, fileCoverage := range coverage {
		if fileCoverage.InputSourceMap == nil {
			// No source map, keep original
			result[filePath] = fileCoverage
			continue
		}

		// Transform this file's coverage
		transformedFiles, err := ct.transformFile(fileCoverage)
		if err != nil {
			return nil, fmt.Errorf("failed to transform file %s: %w", filePath, err)
		}

		// Merge transformed files into result
		for path, fc := range transformedFiles {
			if existing, exists := result[path]; exists {
				// Merge with existing coverage
				ct.mergeCoverage(existing, fc)
			} else {
				result[path] = fc
			}
		}
	}

	return result, nil
}

// transformFile transforms a single file's coverage data
func (ct *CoverageTransformer) transformFile(fc *FileCoverage) (map[string]*FileCoverage, error) {
	if fc.InputSourceMap == nil {
		return map[string]*FileCoverage{fc.Path: fc}, nil
	}

	result := make(map[string]*FileCoverage)

	// Transform statements
	if err := ct.transformStatements(fc, result); err != nil {
		return nil, err
	}

	// Transform functions
	if err := ct.transformFunctions(fc, result); err != nil {
		return nil, err
	}

	// Transform branches
	if err := ct.transformBranches(fc, result); err != nil {
		return nil, err
	}

	return result, nil
}

// transformStatements transforms statement coverage
func (ct *CoverageTransformer) transformStatements(fc *FileCoverage, result map[string]*FileCoverage) error {
	for stmtID, loc := range fc.StatementMap {
		hits, exists := fc.S[stmtID]
		if !exists {
			continue
		}

		// Map location to original source
		mapping, err := ct.sourceMapTransformer.MapLocation(fc.InputSourceMap, loc)
		if err != nil {
			continue // Skip unmappable statements
		}

		// Get or create file coverage for the original source
		targetFC := ct.getOrCreateFileCoverage(result, mapping.Source)

		// Add statement to target file
		newStmtID := ct.getNextID(targetFC.StatementMap)
		targetFC.StatementMap[newStmtID] = mapping.Location
		targetFC.S[newStmtID] = hits
	}

	return nil
}

// transformFunctions transforms function coverage
func (ct *CoverageTransformer) transformFunctions(fc *FileCoverage, result map[string]*FileCoverage) error {
	for fnID, fnMeta := range fc.FnMap {
		hits, exists := fc.F[fnID]
		if !exists {
			continue
		}

		// Map function declaration location
		declMapping, err := ct.sourceMapTransformer.MapLocation(fc.InputSourceMap, fnMeta.Decl)
		if err != nil {
			continue // Skip unmappable functions
		}

		// Map function body location
		locMapping, err := ct.sourceMapTransformer.MapLocation(fc.InputSourceMap, fnMeta.Loc)
		if err != nil {
			locMapping = declMapping // Use declaration mapping if body mapping fails
		}

		// Ensure both mappings are to the same source
		if declMapping.Source != locMapping.Source {
			locMapping = declMapping
		}

		// Get or create file coverage for the original source
		targetFC := ct.getOrCreateFileCoverage(result, declMapping.Source)

		// Add function to target file
		newFnID := ct.getNextID(targetFC.FnMap)
		targetFC.FnMap[newFnID] = FunctionMeta{
			Name: fnMeta.Name,
			Decl: declMapping.Location,
			Loc:  locMapping.Location,
		}
		targetFC.F[newFnID] = hits
	}

	return nil
}

// transformBranches transforms branch coverage
func (ct *CoverageTransformer) transformBranches(fc *FileCoverage, result map[string]*FileCoverage) error {
	for branchID, branchMeta := range fc.BranchMap {
		hits, exists := fc.B[branchID]
		if !exists {
			continue
		}

		// Map branch location
		locMapping, err := ct.sourceMapTransformer.MapLocation(fc.InputSourceMap, branchMeta.Loc)
		if err != nil {
			continue // Skip unmappable branches
		}

		// Map branch locations
		var mappedLocations []Location
		for _, branchLoc := range branchMeta.Locations {
			branchMapping, err := ct.sourceMapTransformer.MapLocation(fc.InputSourceMap, branchLoc)
			if err != nil {
				continue // Skip unmappable branch locations
			}
			if branchMapping.Source == locMapping.Source {
				mappedLocations = append(mappedLocations, branchMapping.Location)
			}
		}

		if len(mappedLocations) == 0 {
			continue // Skip if no locations could be mapped
		}

		// Get or create file coverage for the original source
		targetFC := ct.getOrCreateFileCoverage(result, locMapping.Source)

		// Add branch to target file
		newBranchID := ct.getNextID(targetFC.BranchMap)
		targetFC.BranchMap[newBranchID] = BranchMeta{
			Type:      branchMeta.Type,
			Loc:       locMapping.Location,
			Locations: mappedLocations,
		}
		targetFC.B[newBranchID] = hits
	}

	return nil
}

// getOrCreateFileCoverage gets or creates a FileCoverage for the given path
func (ct *CoverageTransformer) getOrCreateFileCoverage(result map[string]*FileCoverage, path string) *FileCoverage {
	if fc, exists := result[path]; exists {
		return fc
	}

	fc := &FileCoverage{
		Path:         path,
		StatementMap: make(map[string]Location),
		FnMap:        make(map[string]FunctionMeta),
		BranchMap:    make(map[string]BranchMeta),
		S:            make(map[string]int),
		F:            make(map[string]int),
		B:            make(map[string][]int),
	}
	result[path] = fc
	return fc
}

// getNextID generates the next available ID for a map
func (ct *CoverageTransformer) getNextID(m interface{}) string {
	var maxID int

	switch typedMap := m.(type) {
	case map[string]Location:
		for id := range typedMap {
			if intID, err := strconv.Atoi(id); err == nil && intID > maxID {
				maxID = intID
			}
		}
	case map[string]FunctionMeta:
		for id := range typedMap {
			if intID, err := strconv.Atoi(id); err == nil && intID > maxID {
				maxID = intID
			}
		}
	case map[string]BranchMeta:
		for id := range typedMap {
			if intID, err := strconv.Atoi(id); err == nil && intID > maxID {
				maxID = intID
			}
		}
	}

	return strconv.Itoa(maxID + 1)
}

// mergeCoverage merges two FileCoverage objects
func (ct *CoverageTransformer) mergeCoverage(target, source *FileCoverage) {
	// Merge statements
	for id, loc := range source.StatementMap {
		newID := ct.getNextID(target.StatementMap)
		target.StatementMap[newID] = loc
		if hits, exists := source.S[id]; exists {
			target.S[newID] = hits
		}
	}

	// Merge functions
	for id, fn := range source.FnMap {
		newID := ct.getNextID(target.FnMap)
		target.FnMap[newID] = fn
		if hits, exists := source.F[id]; exists {
			target.F[newID] = hits
		}
	}

	// Merge branches
	for id, branch := range source.BranchMap {
		newID := ct.getNextID(target.BranchMap)
		target.BranchMap[newID] = branch
		if hits, exists := source.B[id]; exists {
			target.B[newID] = hits
		}
	}
}

// TransformCoverage is a convenience function to transform coverage data
func TransformCoverage(coverageJSON []byte) ([]byte, error) {
	// Parse input
	coverage, err := ParseCoverageMap(coverageJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse coverage data: %w", err)
	}

	// Transform
	transformer := NewCoverageTransformer()
	transformed, err := transformer.Transform(coverage)
	if err != nil {
		return nil, fmt.Errorf("failed to transform coverage: %w", err)
	}

	// Convert back to JSON
	return transformed.ToJSON()
}
