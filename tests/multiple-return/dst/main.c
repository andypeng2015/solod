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
    so_int res = 0;
    res = x / y;
    *mod = x % y;
    return res;
}

int main(void) {
    so_int b;
    so_int a = vals(&b);
    b = swap(a, b, &a);
    (void)a;
    (void)b;
    so_int m;
    so_int d1 = divide(7, 3, &m);
    so_int d2 = divide(8, 3, &m);
    (void)d1;
    (void)d2;
    (void)m;
    so_int c1;
    vals(&c1);
    (void)c1;
    so_int _d1;
    so_int c2 = vals(&_d1);
    (void)c2;
}
