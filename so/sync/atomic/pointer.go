package atomic

// Pointer is an atomic pointer of type *T. The zero value is a nil *T.
// Pointer must not be copied after first use.
//
//so:extern atomic_Pointer
type Pointer[T any] struct {
	v *T
}

// Load atomically loads and returns the pointer stored in x.
//
//so:extern
func (x *Pointer[T]) Load() *T {
	return x.v
}

// Store atomically stores val into x.
//
//so:extern
func (x *Pointer[T]) Store(val *T) {
	x.v = val
}

// Swap atomically stores new into x and returns the previous value.
//
//so:extern
func (x *Pointer[T]) Swap(new *T) *T {
	old := x.v
	x.v = new
	return old
}

// CompareAndSwap atomically sets x to new if it currently holds old,
// reporting whether the swap happened.
//
//so:extern
func (x *Pointer[T]) CompareAndSwap(old, new *T) bool {
	if x.v == old {
		x.v = new
		return true
	}
	return false
}
