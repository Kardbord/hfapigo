package chatstream

import "testing"

func TestToolCallAccumulator_Merge(t *testing.T) {
	t.Parallel()

	acc := ToolCallAccumulator{}

	// Initial update should store and return provided values.
	id, typ, name := acc.Merge(0, 0, "call_1", "function", "fn")
	if id != "call_1" || typ != "function" || name != "fn" {
		t.Fatalf("unexpected first merge: %q %q %q", id, typ, name)
	}

	// Missing fields should fall back to cached state.
	id, typ, name = acc.Merge(0, 0, "", "", "")
	if id != "call_1" || typ != "function" || name != "fn" {
		t.Fatalf("missing fields should reuse cache: %q %q %q", id, typ, name)
	}

	// Updating only one field should preserve others.
	id, typ, name = acc.Merge(0, 0, "", "", "fn_override")
	if id != "call_1" || typ != "function" || name != "fn_override" {
		t.Fatalf("partial update failed: %q %q %q", id, typ, name)
	}

	// Separate choice/call indices should keep isolated caches.
	idB, typB, nameB := acc.Merge(1, 2, "call_2", "function", "fn2")
	if idB != "call_2" || typB != "function" || nameB != "fn2" {
		t.Fatalf("unexpected second call: %q %q %q", idB, typB, nameB)
	}
	id, typ, name = acc.Merge(0, 0, "", "", "")
	if name != "fn_override" {
		t.Fatalf("first call cache should remain untouched: %q %q %q", id, typ, name)
	}
}

func TestToolCallAccumulator_MergeWithoutMetadata(t *testing.T) {
	t.Parallel()

	acc := ToolCallAccumulator{}

	id, typ, name := acc.Merge(0, 0, "", "", "")
	if id != "" || typ != "" || name != "" {
		t.Fatalf("expected empty metadata, got %q %q %q", id, typ, name)
	}

	id, typ, name = acc.Merge(0, 0, "call_1", "", "")
	if id != "call_1" || typ != "" || name != "" {
		t.Fatalf("unexpected selective metadata: %q %q %q", id, typ, name)
	}

	id, typ, name = acc.Merge(0, 0, "", "", "")
	if id != "call_1" || typ != "" || name != "" {
		t.Fatalf("expected cached id only: %q %q %q", id, typ, name)
	}
}
