package atomic

// Int32 is an atomic int32. The zero value is zero.
// Int32 must not be copied after first use.
type Int32 struct {
	v int32
}

// Load atomically loads and returns the value stored in x.
func (x *Int32) Load() int32 {
	return load(&x.v)
}

// Store atomically stores val into x.
func (x *Int32) Store(val int32) {
	store(&x.v, val)
}

// Add atomically adds delta to x and returns the new value.
func (x *Int32) Add(delta int32) int32 {
	return add(&x.v, delta)
}

// Swap atomically stores new into x and returns the previous value.
func (x *Int32) Swap(new int32) int32 {
	return swap(&x.v, new)
}

// CompareAndSwap atomically sets x to new if it currently holds old,
// reporting whether the swap happened.
func (x *Int32) CompareAndSwap(old, new int32) bool {
	return cas(&x.v, old, new)
}
