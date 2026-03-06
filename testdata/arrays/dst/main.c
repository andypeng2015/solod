#include "main.h"

// -- Implementation --

int main(void) {
    so_Slice a = (so_Slice){(so_int[5]){0}, 5, 5};
    (void)a;
    so_at(so_int, a, 4) = 100;
    so_int x = so_at(so_int, a, 4);
    (void)x;
    so_int l = so_len(a);
    (void)l;
    so_Slice b = (so_Slice){(so_int[5]){1, 2, 3, 4, 5}, 5, 5};
    (void)b;
    so_Slice c = (so_Slice){(so_int[5]){1, 2, 3, 4, 5}, 5, 5};
    (void)c;
    so_Slice d = (so_Slice){(so_int[5]){100, [3] = 400, 500}, 5, 5};
    (void)d;
}
