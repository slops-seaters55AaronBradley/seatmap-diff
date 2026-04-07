package diff

import "fmt"

// ChangeType represents the kind of change detected between two configs.
type ChangeType string

const (
	// ChangeAdded indicates a key/value was added in the new config.
	ChangeAdded ChangeType = "added"
	// ChangeRemoved indicates a key/value was removed from the old config.
	ChangeRemoved ChangeType = "removed"
	// ChangeModified indicates a key/value was changed between configs.
	ChangeModified ChangeType = "modified"
)

// Change represents a single detected difference between two config files.
type Change struct {
	// Path is the dot-separated key path to the changed field (e.g. "server.port").
	Path string
	// Type is the kind of change (added, removed, modified).
	Type ChangeType
	// OldValue is the previous value (nil for added changes).
	OldValue interface{}
	// NewValue is the updated value (nil for removed changes).
	NewValue interface{}
}

// String returns a human-readable summary of the change.
func (c Change) String() string {
	switch c.Type {
	case ChangeAdded:
		return fmt.Sprintf("[+] %s: %v", c.Path, c.NewValue)
	case ChangeRemoved:
		return fmt.Sprintf("[-] %s: %v", c.Path, c.OldValue)
	case ChangeModified:
		return fmt.Sprintf("[~] %s: %v -> %v", c.Path, c.OldValue, c.NewValue)
	default:
		return fmt.Sprintf("[?] %s", c.Path)
	}
}

// Result holds the full set of changes produced by a diff operation.
type Result struct {
	// Changes is the ordered list of detected differences.
	Changes []Change
	// SourceLabel is a display name for the "old" config (e.g. environment name or file path).
	SourceLabel string
	// TargetLabel is a display name for the "new" config.
	TargetLabel string
}

// NewResult initialises an empty Result with the given source and target labels.
func NewResult(sourceLabel, targetLabel string) *Result {
	return &Result{
		Changes:     make([]Change, 0),
		SourceLabel: sourceLabel,
		TargetLabel: targetLabel,
	}
}

// Add appends a change to the result set.
func (r *Result) Add(c Change) {
	r.Changes = append(r.Changes, c)
}

// HasChanges returns true when at least one difference was detected.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a brief string describing the total number of changes
// broken down by type.
func (r *Result) Summary() string {
	var added, removed, modified int
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
