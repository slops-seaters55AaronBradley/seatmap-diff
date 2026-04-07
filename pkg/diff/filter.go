package diff

import (
	"path/filepath"
	"strings"
)

// FilterOptions controls which differences are included in results.
type FilterOptions struct {
	// IncludePaths limits results to differences at or under these dot-separated
	// paths (e.g. "spec.containers"). An empty slice means include everything.
	IncludePaths []string

	// ExcludePaths omits differences at or under these paths.
	ExcludePaths []string

	// Types restricts results to specific change types (Added, Removed, Modified).
	// An empty slice means include all types.
	Types []ChangeType

	// IgnoreKeys is a set of leaf key names whose changes are silently ignored
	// regardless of path (e.g. "updatedAt", "resourceVersion").
	IgnoreKeys []string
}

// Filter returns a new Result containing only the differences that satisfy f.
// If f is nil or contains no constraints, the original result is returned
// unchanged.
func Filter(r *Result, f *FilterOptions) *Result {
	if f == nil || f.isEmpty() {
		return r
	}

	filtered := NewResult()
	for _, diff := range r.Differences {
		if f.matches(diff) {
			filtered.Add(diff)
		}
	}
	return filtered
}

// isEmpty reports whether the FilterOptions impose any constraints.
func (f *FilterOptions) isEmpty() bool {
	return len(f.IncludePaths) == 0 &&
		len(f.ExcludePaths) == 0 &&
		len(f.Types) == 0 &&
		len(f.IgnoreKeys) == 0
}

// matches reports whether a single Difference satisfies all filter criteria.
func (f *FilterOptions) matches(d Difference) bool {
	// Check ignored leaf keys first — cheapest rejection.
	leaf := leafKey(d.Path)
	for _, k := range f.IgnoreKeys {
		if strings.EqualFold(k, leaf) {
			return false
		}
	}

	// Filter by change type.
	if len(f.Types) > 0 && !containsType(f.Types, d.Type) {
		return false
	}

	// Exclude paths take precedence over include paths.
	for _, pattern := range f.ExcludePaths {
		if pathMatches(d.Path, pattern) {
			return false
		}
	}

	// If include paths are specified the difference must match at least one.
	if len(f.IncludePaths) > 0 {
		matched := false
		for _, pattern := range f.IncludePaths {
			if pathMatches(d.Path, pattern) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

// pathMatches reports whether diffPath is equal to or nested under pattern.
// Both arguments use dot-separated notation. Glob wildcards are supported via
// filepath.Match after replacing dots with slashes.
func pathMatches(diffPath, pattern string) bool {
	// Exact match.
	if diffPath == pattern {
		return true
	}
	// Prefix match: pattern is an ancestor of diffPath.
	if strings.HasPrefix(diffPath, pattern+".") {
		return true
	}
	// Glob match using filepath semantics (dots → slashes for matching).
	slashPath := strings.ReplaceAll(diffPath, ".", "/")
	slashPattern := strings.ReplaceAll(pattern, ".", "/")
	ok, err := filepath.Match(slashPattern, slashPath)
	if err == nil && ok {
		return true
	}
	return false
}

// leafKey returns the last segment of a dot-separated path.
func leafKey(path string) string {
	if idx := strings.LastIndex(path, "."); idx >= 0 {
		return path[idx+1:]
	}
	return path
}

// containsType is a small linear search helper.
func containsType(types []ChangeType, t ChangeType) bool {
	for _, ct := range types {
		if ct == t {
			return true
		}
	}
	return false
}
