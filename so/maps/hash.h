#include "so/builtin/builtin.h"
#include <time.h>

#if defined(__APPLE__)
#include <stdlib.h>
#elif defined(__linux__)
#include <sys/random.h>
#endif

// seed returns a random 64-bit seed for hash randomization.
static inline uint64_t maps_seed(void) {
    uint64_t seed = 0;
#if defined(__APPLE__)
    arc4random_buf(&seed, sizeof(seed));
#elif defined(__linux__)
    if (getrandom(&seed, sizeof(seed), 0) != sizeof(seed)) {
        // Fallback to time-based seed.
        struct timespec ts;
        clock_gettime(CLOCK_MONOTONIC, &ts);
        seed ^= (uint64_t)ts.tv_nsec ^ (uint64_t)ts.tv_sec;
    }
#else
    seed = (uint64_t)time(NULL) ^ (uintptr_t)&seed;
#endif
    return seed;
}

// wymum performs 128-bit multiply-and-mix using hardware support.
static inline uint64_t _maps_wymum(uint64_t a, uint64_t b) {
    __uint128_t r = (__uint128_t)a * b;
    return (uint64_t)(r >> 64) ^ (uint64_t)r;
}

// wyr8 reads 8 bytes as a little-endian uint64.
static inline uint64_t _maps_wyr8(const uint8_t* p) {
    uint64_t v;
    memcpy(&v, p, 8);
    return v;
}

// wyr4 reads 4 bytes as a little-endian uint64.
static inline uint64_t _maps_wyr4(const uint8_t* p) {
    uint32_t v;
    memcpy(&v, p, 4);
    return (uint64_t)v;
}

// hash computes wyhash with a per-map seed.
static inline so_int maps_hash(const void* key, size_t len, uint64_t seed) {
    const uint8_t* p = (const uint8_t*)key;
    const uint64_t wyp0 = 0xa0761d6478bd642fULL;
    const uint64_t wyp1 = 0xe7037ed1a0b428dbULL;
    seed = _maps_wymum(seed ^ wyp0, wyp1);
    uint64_t a = 0, b = 0;
    if (len > 16) {
        for (size_t i = 0; i + 16 <= len; i += 16) {
            seed = _maps_wymum(_maps_wyr8(p + i) ^ wyp1,
                               _maps_wyr8(p + i + 8) ^ seed);
        }
        a = _maps_wyr8(p + len - 16);
        b = _maps_wyr8(p + len - 8);
    } else if (len >= 4) {
        a = (_maps_wyr4(p) << 32) | _maps_wyr4(p + ((len >> 3) << 2));
        b = (_maps_wyr4(p + len - 4) << 32) |
            _maps_wyr4(p + len - 4 - ((len >> 3) << 2));
    } else if (len > 0) {
        a = ((uint64_t)p[0] << 16) | ((uint64_t)p[len >> 1] << 8) |
            (uint64_t)p[len - 1];
    }
    uint64_t r = _maps_wymum(wyp1 ^ (uint64_t)len, _maps_wymum(a ^ wyp1, b ^ seed));
    return (so_int)(r >> 16);  // upper 48 bits is the hash value
}

// hashString hashes a string key by its content.
static inline so_int maps_hashString(void* key_ptr, uint64_t seed) {
    so_String* s = (so_String*)key_ptr;
    return maps_hash(s->ptr, s->len, seed);
}

// keyHash hashes a key, dispatching to string or inline hash.
#define maps_keyHash(K, key_ptr, seed) _Generic((K){0}, \
    so_String: maps_hashString(key_ptr, seed),          \
    default: maps_hash((key_ptr), sizeof(K), seed))
