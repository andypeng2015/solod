#include "main.h"
#include "person.ext.h"

int main(void) {
    Account acc = (Account){.name = so_strlit("Alice"), .balance = 100, .flags = {(uint8_t[1]){42}, 1, 1}};
    int64_t balBefore = account_inc_balance(&acc, 50);
    so_println("%s %s %s %lld %lld %s %u", "name =", acc.name.ptr, "balance =", balBefore, acc.balance, "flags[0] =", so_index(acc.flags, uint8_t, 0));
}
