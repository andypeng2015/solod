#pragma once
#include "so/builtin/builtin.h"

// -- Embeds --

#define newObj(T) (alloca(sizeof(T)))
#define freeObj(T, ptr) ((void)(ptr))
#define newMap(K, V, size) (size)