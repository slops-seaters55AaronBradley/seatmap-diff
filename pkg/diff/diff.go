// Package diff provides functionality for comparing YAML and JSON
// infrastructure configuration files across environments.
package diff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeType represents the kind of change detected between two configs.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// Change represents a single detected difference between two config states.
type Change struct {
	Path     string     // dot-separated path to the changed key
	Type     ChangeType
	OldValue interface{}
	NewValue interface{}
}

// Result holds the full diff output for a pair of config files.
type Result struct {
	SourceFile string
	TargetFile string
	Changes    []Change
}

// HasChanges returns true if any differences were detected.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a human-readable summary string of the diff result.
func (r *Result) Summary() string {
	if !r.HasChanges() {
		return fmt.Sprintf("No differences found between %s and %s", r.SourceFile, r.TargetFile)
	}
	added, removed, modified := 0, 0, 0
	for _, c := range r.Changes {
		switch c.Type {
		case ChangeAdded:
			added++
		case ChangeRemoved:
			removed++
		case ChangeModified:
			modified++
		}
	}
	return fmt.Sprintf("%d added, %d removed, %d modified", added, removed, modified)
}

// Compare performs a deep comparison of two parsed config maps and returns
// a Result containing all detected changes.
func Compare(source, target map[string]interface{}, sourceFile, targetFile string) *Result {
	result := &Result{
		SourceFile: sourceFile,
		TargetFile: targetFile,
	}
	changes := []Change{}
	compareRecursive("", source, target, &changes)
	// Sort changes by path for deterministic output
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})
	result.Changes = changes
	return result
}

// compareRecursive walks both maps simultaneously, collecting changes at each level.
func compareRecursive(prefix string, source, target map[string]interface{}, changes *[]Change) {
	// Check for removed or modified keys
	for key, srcVal := range source {
		path := joinPath(prefix, key)
		tgtVal, exists := target[key]
		if !exists {
			*changes = append(*changes, Change{Path: path, Type: ChangeRemoved, OldValue: srcVal})
			continue
		}
		srcMap, srcIsMap := srcVal.(map[string]interface{})
		tgtMap, tgtIsMap := tgtVal.(map[string]interface{})
		if srcIsMap && tgtIsMap {
			compareRecursive(path, srcMap, tgtMap, changes)
		} else if fmt.Sprintf("%v", srcVal) != fmt.Sprintf("%v", tgtVal) {
			*changes = append(*changes, Change{Path: path, Type: ChangeModified, OldValue: srcVal, NewValue: tgtVal})
		}
	}
	// Check for added keys
	for key, tgtVal := range target {
		if _, exists := source[key]; !exists {
			path := joinPath(prefix, key)
			*changes = append(*changes, Change{Path: path, Type: ChangeAdded, NewValue: tgtVal})
		}
	}
}

// joinPath constructs a dot-separated config key path.
func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return strings.Join([]string{prefix, key}, ".")
}
