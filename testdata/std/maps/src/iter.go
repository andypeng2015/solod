package main

import "solod.dev/so/maps"

func iterTest() {
	{
		// Iterate over map.
		m := makeMap()
		seen := make(map[string]bool, m.Len())
		it := m.Iter()
		for it.Next() {
			k, v := it.Key(), it.Value()
			if m.Get(k) != v {
				panic("invalid key-value pair")
			}
			if seen[k] {
				panic("duplicate key")
			}
			seen[k] = true
		}
		if len(seen) != m.Len() {
			panic("missing keys")
		}
		m.Free()
	}
	{
		// Iterate over empty map.
		m := maps.New[string, int](nil, 0)
		it := m.Iter()
		if it.Next() {
			panic("expected no elements")
		}
		m.Free()
	}
}
