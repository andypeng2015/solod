// TryAppend appends elements to a heap-allocated slice, growing it if needed.
// Returns a result with the updated slice or an error if reallocation fails.
// If the allocator is nil, uses the system allocator.
#define slices_TryAppend(T, a, s, ...) ({                 \
    T _vals[] = {__VA_ARGS__};                            \
    size_t _n = sizeof(_vals) / sizeof(T);                \
    slices_tryExtend((a), (s),                            \
                     (so_Slice){(so_byte*)_vals, _n, _n}, \
                     sizeof(T), alignof(so_typeof(T)));   \
})

// Append appends elements to a heap-allocated slice, growing it if needed.
// Returns the updated slice or panics on allocation failure.
// If the allocator is nil, uses the system allocator.
#define slices_Append(T, a, s, ...) ({                 \
    T _vals[] = {__VA_ARGS__};                         \
    size_t _n = sizeof(_vals) / sizeof(T);             \
    slices_extend((a), (s),                            \
                  (so_Slice){(so_byte*)_vals, _n, _n}, \
                  sizeof(T), alignof(so_typeof(T)));   \
})

// TryExtend appends all elements from another slice, growing if needed.
// Returns a result with the updated slice or an error if reallocation fails.
// If the allocator is nil, uses the system allocator.
#define slices_TryExtend(T, a, s, other) \
    slices_tryExtend((a), (s), (other), sizeof(T), alignof(so_typeof(T)))

// Extend appends all elements from another slice, growing if needed.
// Returns the updated slice or panics on allocation failure.
// If the allocator is nil, uses the system allocator.
#define slices_Extend(T, a, s, other) \
    slices_extend((a), (s), (other), sizeof(T), alignof(so_typeof(T)))

// Header returns the Slice header for a given slice.
#define slices_Header(T, s) (s)
