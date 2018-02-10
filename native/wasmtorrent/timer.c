#include "timer.h"

void wt_set_timeout(wt_set_timeout_cb cb, void * user, int ms)
{
#ifdef WASM
  emscripten_async_call(cb, user, ms);
#else
  abort();  
#endif
}

