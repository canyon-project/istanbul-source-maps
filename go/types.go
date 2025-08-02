package istanbul

import "encoding/json"

// Position represents a position in source code (line, column)
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Location represents a range in source code
type Location struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// FunctionMeta represents function metadata
type FunctionMeta struct {
	Name string   `json:"name"`
	Decl Location `json:"decl"`
	Loc  Location `json:"loc"`
}

// BranchMeta represents branch metadata
type BranchMeta struct {
	Type      string     `json:"type"`
	Loc       Location   `json:"loc"`
	Locations []Location `json:"locations"`
}

// SourceMap represents a source map
type SourceMap struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	Names          []string `json:"names"`
	Mappings       string   `json:"mappings"`
	File           string   `json:"file,omitempty"`
	SourceRoot     string   `json:"sourceRoot,omitempty"`
	SourcesContent []string `json:"sourcesContent,omitempty"`
}

// FileCoverage represents coverage data for a single file
type FileCoverage struct {
	Path           string                  `json:"path"`
	StatementMap   map[string]Location     `json:"statementMap"`
	FnMap          map[string]FunctionMeta `json:"fnMap"`
	BranchMap      map[string]BranchMeta   `json:"branchMap"`
	S              map[string]int          `json:"s"` // statement hits
	F              map[string]int          `json:"f"` // function hits
	B              map[string][]int        `json:"b"` // branch hits
	InputSourceMap *SourceMap              `json:"inputSourceMap,omitempty"`
}

// CoverageMap represents coverage data for multiple files
type CoverageMap map[string]*FileCoverage

// MappingResult represents a mapping result from source map
type MappingResult struct {
	Source   string
	Location Location
}

// ParseCoverageMap parses JSON coverage data into CoverageMap
func ParseCoverageMap(data []byte) (CoverageMap, error) {
	var coverage CoverageMap
	if err := json.Unmarshal(data, &coverage); err != nil {
		return nil, err
	}
	return coverage, nil
}

// ToJSON converts CoverageMap to JSON
func (cm CoverageMap) ToJSON() ([]byte, error) {
	return json.MarshalIndent(cm, "", "  ")
}
