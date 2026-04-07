package diff_test

import (
	"testing"

	"github.com/seatmap-diff/pkg/diff"
)

// TestCompare_AddedKey verifies that a key present in the new config but not
// the old is reported as an Added change.
func TestCompare_AddedKey(t *testing.T) {
	old := map[string]interface{}{
		"replicas": 2,
	}
	new := map[string]interface{}{
		"replicas": 2,
		"timeout":  30,
	}

	result := diff.Compare(old, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added {
		t.Errorf("expected change type Added, got %s", result.Changes[0].Type)
	}
	if result.Changes[0].Path != "timeout" {
		t.Errorf("expected path 'timeout', got %s", result.Changes[0].Path)
	}
}

// TestCompare_RemovedKey verifies that a key present in the old config but not
// the new is reported as a Removed change.
func TestCompare_RemovedKey(t *testing.T) {
	old := map[string]interface{}{
		"replicas": 2,
		"timeout":  30,
	}
	new := map[string]interface{}{
		"replicas": 2,
	}

	result := diff.Compare(old, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Removed {
		t.Errorf("expected change type Removed, got %s", result.Changes[0].Type)
	}
}

// TestCompare_ModifiedValue verifies that a changed scalar value is reported
// as a Modified change with correct old and new values.
func TestCompare_ModifiedValue(t *testing.T) {
	old := map[string]interface{}{"replicas": 2}
	new := map[string]interface{}{"replicas": 5}

	result := diff.Compare(old, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != diff.Modified {
		t.Errorf("expected Modified, got %s", c.Type)
	}
	if c.OldValue != 2 {
		t.Errorf("expected OldValue 2, got %v", c.OldValue)
	}
	if c.NewValue != 5 {
		t.Errorf("expected NewValue 5, got %v", c.NewValue)
	}
}

// TestCompare_NestedChange verifies that changes inside nested maps are
// reported with a dot-separated path.
func TestCompare_NestedChange(t *testing.T) {
	old := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}
	new := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "db.prod.internal",
			"port": 5432,
		},
	}

	result := diff.Compare(old, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Path != "database.host" {
		t.Errorf("expected path 'database.host', got %s", result.Changes[0].Path)
	}
}

// TestCompare_NoChanges verifies that identical configs produce an empty
// change set and HasChanges returns false.
func TestCompare_NoChanges(t *testing.T) {
	cfg := map[string]interface{}{
		"replicas": 3,
		"image":    "nginx:1.25",
	}

	result := diff.Compare(cfg, cfg)
	if result.HasChanges() {
		t.Errorf("expected no changes for identical configs, got %d", len(result.Changes))
	}
}

// TestCompare_NilInputs verifies graceful handling when one or both inputs
// are nil.
func TestCompare_NilInputs(t *testing.T) {
	new := map[string]interface{}{"key": "value"}

	result := diff.Compare(nil, new)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change for nil old config, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added {
		t.Errorf("expected Added, got %s", result.Changes[0].Type)
	}
}
