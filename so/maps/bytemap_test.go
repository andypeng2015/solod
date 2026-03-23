package maps

import (
	"encoding/binary"
	"testing"
)

func TestByteMap_SetGet(t *testing.T) {
	ksize, vsize := 4, 4
	key := make([]byte, ksize)
	val := make([]byte, vsize)
	m := NewByteMap(nil, 0, ksize, vsize)

	for i := range 100 {
		encode(i, key)
		encode(i*10, val)
		m.Set(key, val)
	}
	if m.Len() != 100 {
		t.Fatalf("expected length 100, got %d", m.Len())
	}

	out := make([]byte, vsize)
	for i := range 100 {
		encode(i, key)
		if !m.Get(key, out) {
			t.Fatalf("key %d not found", i)
		}
		got := decode(out)
		if got != i*10 {
			t.Fatalf("key %d: expected %d, got %d", i, i*10, got)
		}
	}
}

func TestByteMap_Delete(t *testing.T) {
	ksize, vsize := 4, 4
	key := make([]byte, ksize)
	val := make([]byte, vsize)
	m := NewByteMap(nil, 0, ksize, vsize)

	for i := range 100 {
		encode(i, key)
		encode(i, val)
		m.Set(key, val)
	}
	for i := range 50 {
		encode(i, key)
		if !m.Delete(key) {
			t.Fatalf("delete key %d failed", i)
		}
	}
	if m.Len() != 50 {
		t.Fatalf("expected length 50, got %d", m.Len())
	}

	out := make([]byte, vsize)
	for i := range 50 {
		encode(i, key)
		if m.Get(key, out) {
			t.Fatalf("key %d should not exist", i)
		}
	}
	for i := 50; i < 100; i++ {
		encode(i, key)
		out := make([]byte, vsize)
		if !m.Get(key, out) {
			t.Fatalf("key %d not found", i)
		}
	}
}

func TestByteMap_Overwrite(t *testing.T) {
	m := NewByteMap(nil, 0, 4, 4)
	key := make([]byte, 4)
	val := make([]byte, 4)

	encode(42, key)
	encode(100, val)
	m.Set(key, val)
	encode(200, val)
	m.Set(key, val)

	if m.Len() != 1 {
		t.Fatalf("expected length 1, got %d", m.Len())
	}

	out := make([]byte, 4)
	m.Get(key, out)
	if decode(out) != 200 {
		t.Fatalf("expected 200, got %d", decode(out))
	}
}

func TestByteMap_Missing(t *testing.T) {
	m := NewByteMap(nil, 0, 4, 4)
	key := make([]byte, 4)
	out := make([]byte, 4)

	encode(999, key)
	if m.Get(key, out) {
		t.Fatal("expected key not found")
	}
	if m.Delete(key) {
		t.Fatal("expected delete to return false")
	}
}

func TestByteMap_Grow(t *testing.T) {
	m := NewByteMap(nil, 0, 4, 4)
	key := make([]byte, 4)
	val := make([]byte, 4)

	// Insert enough to trigger multiple resizes
	n := 1000
	for i := range n {
		encode(i, key)
		encode(i, val)
		m.Set(key, val)
	}
	if m.Len() != n {
		t.Fatalf("expected length %d, got %d", n, m.Len())
	}

	out := make([]byte, 4)
	for i := range n {
		encode(i, key)
		if !m.Get(key, out) {
			t.Fatalf("key %d not found after grow", i)
		}
		if decode(out) != i {
			t.Fatalf("key %d: wrong value", i)
		}
	}
}

func encode(n int, b []byte) {
	binary.LittleEndian.PutUint32(b, uint32(n))
}

func decode(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}
