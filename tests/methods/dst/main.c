#include "main.h"

so_int main_Rect_Area(void* self) {
    main_Rect* r = (main_Rect*)self;
    return r->width * r->height;
}

static so_int main_Rect_perim(void* self, so_int n) {
    main_Rect* r = (main_Rect*)self;
    return n * (2 * r->width + 2 * r->height);
}

typedef struct circle {
    so_int radius;
} circle;

static so_int circle_area(void* self) {
    circle* c = (circle*)self;
    return 3 * c->radius * c->radius;
}

typedef struct discard {
} discard;

static so_int discard_Write(void* self, so_Slice p, so_Error* err) {
    (void)self;
    *err = NULL;
    return so_len(p);
}

int main(void) {
    main_Rect r = {.width = 10, .height = 5};
    so_int rArea = main_Rect_Area(&r);
    (void)rArea;
    so_int rPerim = main_Rect_perim(&r, 2);
    (void)rPerim;
    main_Rect* rp = &r;
    so_int rpArea = main_Rect_Area(rp);
    (void)rpArea;
    so_int rpPerim = main_Rect_perim(rp, 2);
    (void)rpPerim;
    circle c = {.radius = 7};
    so_int cArea = circle_area(&c);
    (void)cArea;
    discard d = {};
    so_Slice p = {(uint8_t[2]){'h', 'i'}, 2, 2};
    so_Error err;
    so_int n = discard_Write(&d, p, &err);
    (void)n;
    (void)err;
}
