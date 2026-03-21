package maps

const (
	wyp0 uint64 = 0xa0761d6478bd642f
	wyp1 uint64 = 0xe7037ed1a0b428db
)

// wyhash implements the wyhash algorithm,
// returning a 64-bit hash of the key and seed.
func wyhash(key []byte, seed uint64) uint64 {
	l := len(key)
	seed ^= wymum(seed^wyp0, wyp1)
	var a, b uint64
	if l > 16 {
		for i := 0; i+16 <= l; i += 16 {
			seed = wymum(wyr8(key[i:])^wyp1, wyr8(key[i+8:])^seed)
		}
		a = wyr8(key[l-16:])
		b = wyr8(key[l-8:])
	} else if l >= 4 {
		a = (wyr4(key) << 32) | wyr4(key[(l>>3)<<2:])
		b = (wyr4(key[l-4:]) << 32) | wyr4(key[l-4-((l>>3)<<2):])
	} else if l > 0 {
		a = uint64(key[0])<<16 | uint64(key[l>>1])<<8 | uint64(key[l-1])
	}
	return wymum(wyp1^uint64(l), wymum(a^wyp1, b^seed))
}

// wymum implements the 128-bit multiplication and mixing step of the
// wyhash algorithm, returning the mixed result as a 64-bit value.
func wymum(a, b uint64) uint64 {
	ha, la := a>>32, a&0xFFFFFFFF
	hb, lb := b>>32, b&0xFFFFFFFF
	rh := ha * hb
	rm0 := ha * lb
	rm1 := hb * la
	rl := la * lb
	t := rl + (rm0 << 32)
	c := uint64(0)
	if t < rl {
		c = 1
	}
	lo := t + (rm1 << 32)
	if lo < t {
		c++
	}
	hi := rh + (rm0 >> 32) + (rm1 >> 32) + c
	return hi ^ lo
}

// wyr8 reads 8 bytes from the slice as a little-endian uint64.
func wyr8(p []byte) uint64 {
	return uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24 |
		uint64(p[4])<<32 | uint64(p[5])<<40 | uint64(p[6])<<48 | uint64(p[7])<<56
}

// wyr4 reads 4 bytes from the slice as a little-endian uint64.
func wyr4(p []byte) uint64 {
	return uint64(p[0]) | uint64(p[1])<<8 | uint64(p[2])<<16 | uint64(p[3])<<24
}
