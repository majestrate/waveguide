#include "wget.h"
#include "log.h"

void wt_wget(const char * url, void * user, wt_wget_data_cb data, wt_wget_error_cb err)
{
  wt_info(url);
#ifdef WASM
  emscripten_async_wget_data(url, user, data, err);
#else
  abort();
#endif
}
