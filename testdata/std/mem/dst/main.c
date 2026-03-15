#include "main.h"

// -- Forward declarations (functions and methods) --
static void withDefer(void);

// -- Implementation --

static void withDefer(void) {
    main_Point* p = mem_Alloc(main_Point, (mem_Allocator){0});
    p->x = 11;
    p->y = 22;
    if (p->x != 11 || p->y != 22) {
        mem_Free(main_Point, (mem_Allocator){0}, p);
        so_panic("unexpected value");
    }
    mem_Free(main_Point, (mem_Allocator){0}, p);
}

int main(void) {
    {
        // TryAlloc and Free.
        so_Result _res1 = mem_TryAlloc(main_Point, mem_System);
        main_Point* p = _res1.val.as_ptr;
        so_Error err = _res1.err;
        if (err != NULL) {
            so_panic("Alloc: allocation failed");
        }
        p->x = 11;
        p->y = 22;
        if (p->x != 11 || p->y != 22) {
            so_panic("Alloc: unexpected value");
        }
        mem_Free(main_Point, mem_System, p);
    }
    {
        // TryAllocSlice and FreeSlice.
        so_Result _res2 = mem_TryAllocSlice(so_int, mem_System, 3, 3);
        so_Slice slice = _res2.val.as_slice;
        so_Error err = _res2.err;
        if (err != NULL) {
            so_panic("AllocSlice: allocation failed");
        }
        so_at(so_int, slice, 0) = 11;
        so_at(so_int, slice, 1) = 22;
        so_at(so_int, slice, 2) = 33;
        if (so_at(so_int, slice, 0) != 11 || so_at(so_int, slice, 1) != 22 || so_at(so_int, slice, 2) != 33) {
            so_panic("AllocSlice: unexpected value");
        }
        mem_FreeSlice(so_int, mem_System, slice);
    }
    {
        // Alloc/Free with default allocator.
        main_Point* p = mem_Alloc(main_Point, (mem_Allocator){0});
        p->x = 11;
        p->y = 22;
        if (p->x != 11 || p->y != 22) {
            so_panic("New: unexpected value");
        }
        mem_Free(main_Point, (mem_Allocator){0}, p);
    }
    {
        // AllocSlice/FreeSlice with default allocator.
        so_Slice slice = mem_AllocSlice(so_int, (mem_Allocator){0}, 3, 3);
        so_at(so_int, slice, 0) = 11;
        so_at(so_int, slice, 1) = 22;
        so_at(so_int, slice, 2) = 33;
        if (so_at(so_int, slice, 0) != 11 || so_at(so_int, slice, 1) != 22 || so_at(so_int, slice, 2) != 33) {
            so_panic("NewSlice: unexpected value");
        }
        mem_FreeSlice(so_int, (mem_Allocator){0}, slice);
    }
    {
        // Free with nil or an empty slice.
        main_Point* p = NULL;
        mem_Free(main_Point, (mem_Allocator){0}, p);
        so_Slice empty = {0};
        mem_FreeSlice(so_int, (mem_Allocator){0}, empty);
    }
    {
        // Free string.
        so_Slice b = mem_AllocSlice(so_byte, (mem_Allocator){0}, 3, 3);
        so_at(so_byte, b, 0) = 'h';
        so_at(so_byte, b, 1) = 'i';
        so_at(so_byte, b, 2) = '!';
        so_String s = so_bytes_string(b);
        mem_FreeString((mem_Allocator){0}, s);
    }
    withDefer();
}
