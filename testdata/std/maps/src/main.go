package main

import "solod.dev/so/maps"

func makeMap() maps.Map[string, int] {
	m := maps.New[string, int](nil, 0)
	m.Set("abc", 11)
	m.Set("def", 22)
	m.Set("xyz", 33)
	return m
}

func main() {
	mapTest()
	iterTest()
}
