#include "main.h"

// -- Types --

typedef struct point point;

// Composite literals passed to function-like C macros (len, cap, append,
// copy, clear, indexing, slicing, map access) emit braced initializers whose
// commas would be misread by the preprocessor as macro argument separators.
// Each such argument must be wrapped in parentheses.
typedef struct point {
    so_int x;
    so_int y;
} point;

// -- Implementation --

int main(void) {
    // len/cap of a slice literal.
    if (so_len(((so_Slice){(so_int[3]){1, 2, 3}, 3, 3})) != 3) {
        so_panic("len");
    }
    if (so_cap(((so_Slice){(so_int[3]){1, 2, 3}, 3, 3})) != 3) {
        so_panic("cap");
    }
    // index and slice of a slice literal.
    if (so_at(so_int, ((so_Slice){(so_int[3]){10, 20, 30}, 3, 3}), 1) != 20) {
        so_panic("index");
    }
    if (so_len(so_slice(so_int, ((so_Slice){(so_int[4]){1, 2, 3, 4}, 4, 4}), 1, 3)) != 2) {
        so_panic("slice");
    }
    if (so_len(so_slice(so_int, ((so_Slice){(so_int[4]){1, 2, 3, 4}, 4, 4}), 2, ((so_Slice){(so_int[4]){1, 2, 3, 4}, 4, 4}).len)) != 2) {
        so_panic("slice open-ended");
    }
    // address of an element of a slice literal.
    (void)&so_at(so_int, ((so_Slice){(so_int[3]){5, 6, 7}, 3, 3}), 0);
    // slice-to-array conversion of a literal.
    so_int arr[2];
    memcpy(arr, so_slice_array(((so_Slice){(so_int[2]){7, 8}, 2, 2}), 2), sizeof(arr));
    if (arr[0] != 7) {
        so_panic("slice-to-array");
    }
    // byte slice literal to string conversion.
    if (so_string_ne(so_bytes_string(((so_Slice){(so_byte[2]){'h', 'i'}, 2, 2})), so_str("hi"))) {
        so_panic("byte slice to string");
    }
    if (so_string_ne(so_bytes_string(((so_Slice){(so_byte[1]){(so_byte)(97)}, 1, 1})), so_str("a"))) {
        so_panic("byte slice to string");
    }
    // copy from a slice literal.
    so_Slice dst = so_make_slice(so_int, 3, 3);
    so_copy(so_int, dst, ((so_Slice){(so_int[3]){1, 2, 3}, 3, 3}));
    if (so_at(so_int, dst, 2) != 3) {
        so_panic("copy");
    }
    // append a composite-literal value.
    so_Slice pts = so_make_slice(point, 0, 2);
    pts = so_append(point, pts, ((point){1, 2}));
    if (so_at(point, pts, 0).y != 2) {
        so_panic("append value");
    }
    // clear a slice literal (exercises the macro; no observable effect).
    so_clear(so_int, ((so_Slice){(so_int[3]){1, 2, 3}, 3, 3}));
    // map with a composite-literal value.
    so_Map* mv = so_make_map(so_int, point, 1);
    so_map_set(so_int, point, mv, 0, ((point){3, 4}));
    if (so_map_get(so_int, point, mv, 0).x != 3) {
        so_panic("map value");
    }
    // map with a composite-literal pointer value.
    so_Map* mp = so_make_map(so_int, point*, 1);
    so_map_set(so_int, point*, mp, 1, (&(point){8, 9}));
    if (so_map_get(so_int, point*, mp, 1)->y != 9) {
        so_panic("map pointer value");
    }
    // map with a composite-literal key.
    so_Map* mk = so_make_map(point, so_int, 1);
    so_map_set(point, so_int, mk, ((point){1, 2}), 42);
    if (so_map_get(point, so_int, mk, ((point){1, 2})) != 42) {
        so_panic("map key");
    }
    so_int v = so_map_get(point, so_int, mk, ((point){1, 2}));
    bool ok = so_map_has(point, mk, ((point){1, 2}));
    if (!ok || v != 42) {
        so_panic("map key comma-ok");
    }
    return 0;
}
