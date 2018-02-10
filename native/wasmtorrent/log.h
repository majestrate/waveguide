#ifndef _WT_LOG_H
#define _WT_LOG_H
#include "wasm.h"
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

  struct wt_logger;

  extern struct wt_logger wt_log;
  
  bool WT_EXPORT wt_log_init();

  void WT_EXPORT wt_info(const char * msg);
  void WT_EXPORT wt_warn(const char * msg);
  void WT_EXPORT wt_error(const char * msg);
  
#ifdef __cplusplus
}
#endif
#endif
