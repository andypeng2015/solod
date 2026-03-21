package main

import _ "embed"

//so:embed main.h
var header string

//so:extern
func newObj[T any]() *T {
	return nil
}

//so:extern
func freeObj[T any](ptr *T) {
}

//so:extern
func newMap[K comparable, V any](size int) int {
	return size
}

func main() {
	{
		// Generic extern function (single type parameter).
		var v *int = newObj[int]()
		*v = 42
		if *v != 42 {
			panic("unexpected value")
		}
		freeObj(v)
	}
	{
		// Generic extern function (multiple type parameters).
		m := newMap[string, int](10)
		if m != 10 {
			panic("unexpected map size")
		}
	}
}
