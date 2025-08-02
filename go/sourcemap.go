package istanbul

import (
	"fmt"

	"github.com/go-sourcemap/sourcemap"
)

// SourceMapTransformer handles source map transformations
type SourceMapTransformer struct {
	cache map[string]*sourcemap.Consumer
}

// NewSourceMapTransformer creates a new transformer
func NewSourceMapTransformer() *SourceMapTransformer {
	return &SourceMapTransformer{
		cache: make(map[string]*sourcemap.Consumer),
	}
}

// GetOriginalPosition maps a generated position to original position
func (smt *SourceMapTransformer) GetOriginalPosition(sm *SourceMap, pos Position) (*MappingResult, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("%s:%s", sm.File, sm.Mappings)

	// Check cache
	consumer, exists := smt.cache[cacheKey]
	if !exists {
		// Parse source map
		var err error
		consumer, err = sourcemap.Parse(sm.File, []byte(sourceMapToJSON(sm)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse source map: %w", err)
		}
		smt.cache[cacheKey] = consumer
	}

	// Get original position
	file, _, line, col, ok := consumer.Source(pos.Line, pos.Column)
	if !ok {
		return nil, fmt.Errorf("no mapping found for position %d:%d", pos.Line, pos.Column)
	}

	return &MappingResult{
		Source: file,
		Location: Location{
			Start: Position{Line: line, Column: col},
			End:   Position{Line: line, Column: col}, // For now, assume single position
		},
	}, nil
}

// MapLocation maps a generated location to original location
func (smt *SourceMapTransformer) MapLocation(sm *SourceMap, loc Location) (*MappingResult, error) {
	// Map start position
	startMapping, err := smt.GetOriginalPosition(sm, loc.Start)
	if err != nil {
		return nil, err
	}

	// Map end position
	endMapping, err := smt.GetOriginalPosition(sm, loc.End)
	if err != nil {
		// If end mapping fails, use start mapping
		endMapping = &MappingResult{
			Source: startMapping.Source,
			Location: Location{
				Start: startMapping.Location.Start,
				End:   startMapping.Location.Start,
			},
		}
	}

	// Ensure both positions map to the same source
	if startMapping.Source != endMapping.Source {
		return startMapping, nil // Use start mapping if sources differ
	}

	return &MappingResult{
		Source: startMapping.Source,
		Location: Location{
			Start: startMapping.Location.Start,
			End:   endMapping.Location.End,
		},
	}, nil
}

// sourceMapToJSON converts SourceMap struct to JSON string for parsing
func sourceMapToJSON(sm *SourceMap) string {
	// Create a minimal source map JSON
	return fmt.Sprintf(`{
		"version": %d,
		"sources": %s,
		"names": %s,
		"mappings": "%s",
		"file": "%s",
		"sourceRoot": "%s"
	}`,
		sm.Version,
		arrayToJSON(sm.Sources),
		arrayToJSON(sm.Names),
		sm.Mappings,
		sm.File,
		sm.SourceRoot,
	)
}

// arrayToJSON converts string array to JSON array string
func arrayToJSON(arr []string) string {
	if len(arr) == 0 {
		return "[]"
	}

	result := "["
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`"%s"`, s)
	}
	result += "]"
	return result
}
