#ifndef _WT_METAINFO_H
#define _WT_METAINFO_H
#include "wasm.h"
#include <stdint.h>
#include <stdbool.h>
#ifdef __cplusplus
extern "C" {
#endif

  struct wt_metainfo;

  void wt_metainfo_init(struct wt_metainfo ** meta);
  void wt_metainfo_free(struct wt_metainfo ** meta);

  /** @brief load metainfo fully formed from buffer, does not copy buffer */
  bool wt_metainfo_loadbuffer(struct wt_metainfo * meta, const void * buff, size_t sz);

  /**@brief calculate metainfo's infohash value */
  void wt_metainfo_calc_infohash(struct wt_metainfo * metainfo, uint8_t * ih);

#ifdef __cplusplus
}
#endif
#endif
