#ifndef _WT_WGET_H
#define _WT_WGET_H
#include "wasm.h"

#ifdef __cplusplus
extern "C" {
#endif

  typedef void (*wt_wget_data_cb)(void * user, void * data, int sz);
  typedef void (*wt_wget_error_cb)(void * user);
  void WT_EXPORT wt_wget(const char * url, void * user, wt_wget_data_cb data, wt_wget_error_cb err);
  
#ifdef __cplusplus
}
#endif
#endif 
