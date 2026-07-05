package main

import (
	"solod.dev/so/maps"
	"solod.dev/so/mem"
	"solod.dev/so/testing"
)

func TestIter(t *testing.T) {
	m := makeMap()
	defer m.Free()

	seen := make(map[string]bool, m.Len())
	it := m.Iter()
	for it.Next() {
		k, v := it.Key(), it.Value()
		if m.Get(k) != v {
			t.Error("invalid key-value pair")
		}
		if seen[k] {
			t.Error("duplicate key")
		}
		seen[k] = true
	}
	if len(seen) != m.Len() {
		t.Error("missing keys")
	}
}

func TestIter_Empty(t *testing.T) {
	m := maps.New[string, int](mem.System, 0)
	defer m.Free()

	it := m.Iter()
	if it.Next() {
		t.Error("expected no elements")
	}
}
