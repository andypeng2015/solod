package main

// Composite literals passed to function-like C macros (len, cap, append,
// copy, clear, indexing, slicing, map access) emit braced initializers whose
// commas would be misread by the preprocessor as macro argument separators.
// Each such argument must be wrapped in parentheses.

type point struct {
	x, y int
}

func main() {
	// len/cap of a slice literal.
	if len([]int{1, 2, 3}) != 3 {
		panic("len")
	}
	if cap([]int{1, 2, 3}) != 3 {
		panic("cap")
	}

	// index and slice of a slice literal.
	if []int{10, 20, 30}[1] != 20 {
		panic("index")
	}
	if len([]int{1, 2, 3, 4}[1:3]) != 2 {
		panic("slice")
	}
	if len([]int{1, 2, 3, 4}[2:]) != 2 {
		panic("slice open-ended")
	}

	// address of an element of a slice literal.
	_ = &[]int{5, 6, 7}[0]

	// slice-to-array conversion of a literal.
	arr := [2]int([]int{7, 8})
	if arr[0] != 7 {
		panic("slice-to-array")
	}

	// byte slice literal to string conversion.
	if string([]byte{'h', 'i'}) != "hi" {
		panic("byte slice to string")
	}
	if string([]byte{byte(97)}) != "a" {
		panic("byte slice to string")
	}

	// copy from a slice literal.
	dst := make([]int, 3)
	copy(dst, []int{1, 2, 3})
	if dst[2] != 3 {
		panic("copy")
	}

	// append a composite-literal value.
	pts := make([]point, 0, 2)
	pts = append(pts, point{1, 2})
	if pts[0].y != 2 {
		panic("append value")
	}

	// clear a slice literal (exercises the macro; no observable effect).
	clear([]int{1, 2, 3})

	// map with a composite-literal value.
	mv := make(map[int]point, 1)
	mv[0] = point{3, 4}
	if mv[0].x != 3 {
		panic("map value")
	}

	// map with a composite-literal pointer value.
	mp := make(map[int]*point, 1)
	mp[1] = &point{8, 9}
	if mp[1].y != 9 {
		panic("map pointer value")
	}

	// map with a composite-literal key.
	mk := make(map[point]int, 1)
	mk[point{1, 2}] = 42
	if mk[point{1, 2}] != 42 {
		panic("map key")
	}
	v, ok := mk[point{1, 2}]
	if !ok || v != 42 {
		panic("map key comma-ok")
	}
}
