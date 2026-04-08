#include "so/builtin/builtin.h"
#include "so/time/time.h"

volatile int64_t sinkInt = 0;
volatile so_String sinkStr = {0};
volatile time_Time sinkTime = {0};
