package diff

import (
	"testing"
)

// TestFilter_ExcludePath verifies that changes matching an excluded path are removed.
func TestFilter_ExcludePath(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Path: "metadata.labels.version", Type: ChangeModified, OldValue: "v1", NewValue: "v2"},
			{Path: "spec.replicas", Type: ChangeModified, OldValue: 1, NewValue: 3},
		},
	}

	opts := FilterOptions{
		ExcludePaths: []string{"metadata.labels.version"},
	}

	filtered := Filter(result, opts)

	if len(filtered.Changes) != 1 {
		t.Fatalf("expected 1 change after filtering, got %d", len(filtered.Changes))
	}
	if filtered.Changes[0].Path != "spec.replicas" {
		t.Errorf("expected remaining change at 'spec.replicas', got '%s'", filtered.Changes[0].Path)
	}
}

// TestFilter_IncludePath verifies that only changes matching an included path are kept.
func TestFilter_IncludePath(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Path: "metadata.name", Type: ChangeModified, OldValue: "foo", NewValue: "bar"},
			{Path: "spec.replicas", Type: ChangeModified, OldValue: 1, NewValue: 3},
			{Path: "spec.template.image", Type: ChangeModified, OldValue: "nginx:1.19", NewValue: "nginx:1.21"},
		},
	}

	opts := FilterOptions{
		IncludePaths: []string{"spec"},
	}

	filtered := Filter(result, opts)

	if len(filtered.Changes) != 2 {
		t.Fatalf("expected 2 changes after include filter, got %d", len(filtered.Changes))
	}
}

// TestFilter_ByChangeType verifies that filtering by change type works correctly.
func TestFilter_ByChangeType(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Path: "a", Type: ChangeAdded, NewValue: "x"},
			{Path: "b", Type: ChangeRemoved, OldValue: "y"},
			{Path: "c", Type: ChangeModified, OldValue: "old", NewValue: "new"},
		},
	}

	opts := FilterOptions{
		ChangeTypes: []ChangeType{ChangeAdded, ChangeRemoved},
	}

	filtered := Filter(result, opts)

	if len(filtered.Changes) != 2 {
		t.Fatalf("expected 2 changes after type filter, got %d", len(filtered.Changes))
	}
	for _, ch := range filtered.Changes {
		if ch.Type == ChangeModified {
			t.Errorf("unexpected ChangeModified in filtered result")
		}
	}
}

// TestFilter_NoOptions verifies that an empty FilterOptions returns all changes.
func TestFilter_NoOptions(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Path: "a", Type: ChangeAdded, NewValue: "x"},
			{Path: "b", Type: ChangeModified, OldValue: "1", NewValue: "2"},
		},
	}

	filtered := Filter(result, FilterOptions{})

	if len(filtered.Changes) != len(result.Changes) {
		t.Errorf("expected %d changes with no filter options, got %d", len(result.Changes), len(filtered.Changes))
	}
}

// TestFilter_WildcardPath verifies that wildcard path patterns match nested keys.
func TestFilter_WildcardPath(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Path: "spec.containers.env.SECRET_KEY", Type: ChangeModified, OldValue: "old", NewValue: "new"},
			{Path: "spec.replicas", Type: ChangeModified, OldValue: 1, NewValue: 2},
		},
	}

	opts := FilterOptions{
		ExcludePaths: []string{"spec.containers.env.*"},
	}

	filtered := Filter(result, opts)

	if len(filtered.Changes) != 1 {
		t.Fatalf("expected 1 change after wildcard exclude, got %d", len(filtered.Changes))
	}
	if filtered.Changes[0].Path != "spec.replicas" {
		t.Errorf("expected 'spec.replicas' to remain, got '%s'", filtered.Changes[0].Path)
	}
}

// TestFilter_CombinedOptions verifies that include paths and change types can be combined.
func TestFilter_CombinedOptions(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Path: "spec.image", Type: ChangeModified, OldValue: "v1", NewValue: "v2"},
			{Path: "spec.replicas", Type: ChangeAdded, NewValue: 3},
			{Path: "metadata.name", Type: ChangeModified, OldValue: "a", NewValue: "b"},
		},
	}

	opts := FilterOptions{
		IncludePaths: []string{"spec"},
		ChangeTypes:  []ChangeType{ChangeModified},
	}

	filtered := Filter(result, opts)

	if len(filtered.Changes) != 1 {
		t.Fatalf("expected 1 change with combined options, got %d", len(filtered.Changes))
	}
	if filtered.Changes[0].Path != "spec.image" {
		t.Errorf("expected 'spec.image', got '%s'", filtered.Changes[0].Path)
	}
}
