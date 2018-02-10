#include "mainloop.h"
#include "torrent.h"
#include <stdbool.h>

static void _run_mainloop(void * user)
{
  struct wt_client * cl = (struct wt_client *) user;
  wt_tick(cl);
}

void wt_mainloop(struct wt_client * cl)
{
#ifdef WASM
  emscripten_set_main_loop_arg(_run_mainloop, cl, 0, true);
#else
  _run_mainloop(cl);
#endif
}
