#include "main.h"

static so_int vals(so_int* _r1) {
    *_r1 = 7;
    return 3;
}

static so_int swap(so_int x, so_int y, so_int* _r1) {
    *_r1 = x;
    return y;
}

static so_int divide(so_int x, so_int y, so_int* mod) {
    *mod = x % y;
    return x / y;
}

int main(void) {
    so_int a, b;
    a = vals(&b);
    b = swap(a, b, &a);
    (void)a;
    (void)b;
    so_int d1, m;
    d1 = divide(7, 3, &m);
    so_int d2;
    d2 = divide(8, 3, &m);
    (void)d1;
    (void)d2;
    (void)m;
    so_int c;
    vals(&c);
    (void)c;
}
