#include <stdio.h>
#include <stdint.h>
#include "so/builtin/builtin.h"

typedef struct {
    so_String name;
    int64_t balance;
    so_Slice flags;
} account;

int64_t account_inc_balance(account* a, int64_t amount);