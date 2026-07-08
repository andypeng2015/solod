package atomic

// Bool is an atomic boolean value. The zero value is false.
// Bool must not be copied after first use.
type Bool struct {
	v bool
}

// Load atomically loads and returns the value stored in x.
func (x *Bool) Load() bool {
	return load(&x.v)
}

// Store atomically stores val into x.
func (x *Bool) Store(val bool) {
	store(&x.v, val)
}

// Swap atomically stores new into x and returns the previous value.
func (x *Bool) Swap(new bool) bool {
	return swap(&x.v, new)
}

// CompareAndSwap atomically sets x to new if it currently holds old,
// reporting whether the swap happened.
func (x *Bool) CompareAndSwap(old, new bool) bool {
	return cas(&x.v, old, new)
}
