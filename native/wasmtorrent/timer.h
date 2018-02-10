#ifndef _WT_TIMER_H
#define _WT_TIMER_H
#include "wasm.h"

#ifdef __cplusplus
extern "C" {
#endif

  typedef void (*wt_set_timeout_cb)(void *);
  
  void wt_set_timeout(wt_set_timeout_cb vb, void * user, int ms);
  
#ifdef __cplusplus
}
#endif
#endif
