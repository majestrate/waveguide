#ifndef _WT_MAINLOOP_H
#define _WT_MAINLOOP_H
#include "wasm.h"

#ifdef __cplusplus
extern "C" {
#endif

  /* foward declare */
  struct wt_client;
  
  void wt_mainloop(struct wt_client * cl);
  
#ifdef __cplusplus
}
#endif
#endif
