package main

type Movie struct {
	year     int
	ratingFn func(m Movie) int
}

func freshness(m Movie) int {
	return m.year - 1970
}

func main() {
	m1 := Movie{year: 2020, ratingFn: freshness}
	s1 := m1.ratingFn(m1) // 50
	if s1 != 50 {
		panic("unexpected s1")
	}

	m2 := Movie{year: 1995, ratingFn: freshness}
	s2 := m2.ratingFn(m2) // 25
	if s2 != 25 {
		panic("unexpected s2")
	}
}
