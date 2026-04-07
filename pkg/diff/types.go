package diff

// DiffType represents the kind of change detected between two configs.
type DiffType string

const (
	// DiffAdded indicates a key/value present in the new config but not the old.
	DiffAdded DiffType = "added"
	// DiffRemoved indicates a key/value present in the old config but not the new.
	DiffRemoved DiffType = "removed"
	// DiffModified indicates a key present in both configs but with a changed value.
	DiffModified DiffType = "modified"
)

// Change represents a single detected difference between two config states.
type Change struct {
	// Path is the dot-separated key path to the changed field (e.g. "server.port").
	Path string `json:"path" yaml:"path"`
	// Type describes the nature of the change: added, removed, or modified.
	Type DiffType `json:"type" yaml:"type"`
	// OldValue holds the previous value; nil for added fields.
	OldValue interface{} `json:"old_value,omitempty" yaml:"old_value,omitempty"`
	// NewValue holds the updated value; nil for removed fields.
	NewValue interface{} `json:"new_value,omitempty" yaml:"new_value,omitempty"`
}

// Result aggregates all changes produced by a single diff operation.
type Result struct {
	// Changes is the ordered list of detected differences.
	Changes []Change `json:"changes" yaml:"changes"`
	// Summary provides a quick breakdown of change counts by type.
	Summary Summary `json:"summary" yaml:"summary"`
}

// Summary holds aggregate counts for a diff result.
type Summary struct {
	Added    int `json:"added" yaml:"added"`
	Removed  int `json:"removed" yaml:"removed"`
	Modified int `json:"modified" yaml:"modified"`
	Total    int `json:"total" yaml:"total"`
}

// AuditEntry records a single audited diff event, associating metadata with
// the changes detected between two environments or config revisions.
type AuditEntry struct {
	// ID is a unique identifier for this audit record (typically a UUID or timestamp hash).
	ID string `json:"id" yaml:"id"`
	// Timestamp is the RFC3339-formatted time the audit was performed.
	Timestamp string `json:"timestamp" yaml:"timestamp"`
	// Environment identifies the target deployment environment (e.g. "staging", "prod").
	Environment string `json:"environment" yaml:"environment"`
	// BaseFile is the path or label of the baseline config file.
	BaseFile string `json:"base_file" yaml:"base_file"`
	// TargetFile is the path or label of the config being compared against the base.
	TargetFile string `json:"target_file" yaml:"target_file"`
	// Result contains the full diff result for this audit entry.
	Result Result `json:"result" yaml:"result"`
}

// NewResult constructs a Result from a slice of Changes and computes its Summary.
func NewResult(changes []Change) Result {
	s := Summary{Total: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case DiffAdded:
			s.Added++
		case DiffRemoved:
			s.Removed++
		case DiffModified:
			s.Modified++
		}
	}
	return Result{
		Changes: changes,
		Summary: s,
	}
}
