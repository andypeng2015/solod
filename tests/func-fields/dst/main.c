#include "main.h"
static so_int freshness(main_Movie m);

static so_int freshness(main_Movie m) {
    return m.year - 1970;
}

int main(void) {
    main_Movie m1 = (main_Movie){.year = 2020, .ratingFn = freshness};
    so_int s1 = m1.ratingFn(m1);
    so_println("%lld", s1);
    main_Movie m2 = (main_Movie){.year = 1995, .ratingFn = freshness};
    so_int s2 = m2.ratingFn(m2);
    so_println("%lld", s2);
}
