package atomic

//so:embed atomic.h
var atomic_h string

// load atomically loads and returns the value at p.
//
//so:extern so_atomic_load
func load[T any](p *T) T {
	return *p
}

// store atomically sets the value at p to v.
//
//so:extern so_atomic_store
func store[T any](p *T, v T) {
	*p = v
}

// add atomically adds delta to the value at p and returns the new value.
//
//so:extern so_atomic_add
func add[T any](p *T, delta T) T {
	_ = delta
	return *p
}

// swap atomically sets the value at p to v and returns the previous value.
//
//so:extern so_atomic_swap
func swap[T any](p *T, v T) T {
	_ = v
	return *p
}

// cas atomically sets the value at p to new if it equals old,
// reporting whether the swap happened.
//
//so:extern so_atomic_cas
func cas[T any](p *T, old, new T) bool {
	_, _, _ = p, old, new
	return false
}
