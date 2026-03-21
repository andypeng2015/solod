#include "main.h"

// -- Implementation --

int main(void) {
    {
        // SetGet: insert 3 entries, verify all values
        maps_Map m = maps_New(so_String, so_int, (mem_Allocator){0}, 0);
        maps_Map_Set(so_String, so_int, &m, so_str("abc"), 11);
        maps_Map_Set(so_String, so_int, &m, so_str("def"), 22);
        maps_Map_Set(so_String, so_int, &m, so_str("xyz"), 33);
        if (maps_Map_Get(so_String, so_int, &m, so_str("abc")) != 11) {
            so_panic("want abc = 11");
        }
        so_String key = so_str("abc");
        if (maps_Map_Get(so_String, so_int, &m, key) != 11) {
            so_panic("want abc = 11 for key = abc");
        }
        if (maps_Map_Get(so_String, so_int, &m, so_str("def")) != 22) {
            so_panic("want def = 22");
        }
        if (maps_Map_Get(so_String, so_int, &m, so_str("xyz")) != 33) {
            so_panic("want xyz = 33");
        }
        maps_Map_Free(so_String, so_int, &m);
    }
    {
        // Delete: insert 3 entries, delete one, verify
        maps_Map m = maps_New(so_String, so_int, (mem_Allocator){0}, 0);
        maps_Map_Set(so_String, so_int, &m, so_str("abc"), 11);
        maps_Map_Set(so_String, so_int, &m, so_str("def"), 22);
        maps_Map_Set(so_String, so_int, &m, so_str("xyz"), 33);
        maps_Map_Delete(so_String, so_int, &m, so_str("def"));
        if (maps_Map_Get(so_String, so_int, &m, so_str("def")) != 0) {
            so_panic("want def = 0 after delete");
        }
        if (maps_Map_Get(so_String, so_int, &m, so_str("abc")) != 11) {
            so_panic("want abc = 11 after delete");
        }
        if (maps_Map_Get(so_String, so_int, &m, so_str("xyz")) != 33) {
            so_panic("want xyz = 33 after delete");
        }
        maps_Map_Free(so_String, so_int, &m);
    }
    {
        // Overwrite: set same key twice, verify latest value wins
        maps_Map m = maps_New(so_String, so_int, (mem_Allocator){0}, 0);
        maps_Map_Set(so_String, so_int, &m, so_str("key"), 100);
        maps_Map_Set(so_String, so_int, &m, so_str("key"), 200);
        if (maps_Map_Get(so_String, so_int, &m, so_str("key")) != 200) {
            so_panic("want key = 200 after overwrite");
        }
        maps_Map_Free(so_String, so_int, &m);
    }
    {
        // Missing: get non-existent key returns zero value
        maps_Map m = maps_New(so_String, so_int, (mem_Allocator){0}, 0);
        if (maps_Map_Get(so_String, so_int, &m, so_str("missing")) != 0) {
            so_panic("want missing = 0");
        }
        maps_Map_Free(so_String, so_int, &m);
    }
    {
        // Grow: insert 100 int-keyed entries, verify all retrievable
        maps_Map m = maps_New(so_int, so_int, (mem_Allocator){0}, 0);
        for (so_int i = 0; i < 100; i++) {
            maps_Map_Set(so_int, so_int, &m, i, i * 10);
        }
        for (so_int i = 0; i < 100; i++) {
            if (maps_Map_Get(so_int, so_int, &m, i) != i * 10) {
                so_panic("wrong value after grow");
            }
        }
        maps_Map_Free(so_int, so_int, &m);
    }
}
