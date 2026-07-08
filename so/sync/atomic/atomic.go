// Package atomic provides low-level atomic memory primitives
// useful for implementing synchronization algorithms.
//
// These functions require great care to be used correctly. Except for special,
// low-level applications, synchronization is better done with conc.Chan or the
// facilities of the sync package.
//
// Each type's zero value is a usable, zeroed atomic; no initialization is
// needed. All operations use sequentially consistent ordering. The types must
// not be copied after first use.
package atomic
