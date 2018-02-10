#ifndef _WT_TORRENT_H
#define _WT_TORRENT_H
#include "wasm.h"
#include "metainfo.h"
#include <stdbool.h>
#ifdef __cplusplus
extern "C" {
#endif

  struct wt_client;

  void WT_EXPORT wt_client_init(struct wt_client ** t);
  void WT_EXPORT wt_client_free(struct wt_client ** t);

  bool WT_EXPORT wt_client_has_torrent(struct wt_client * t, const uint8_t * infohash);
  /** copies meta info */
  void WT_EXPORT wt_client_add_torrent(struct wt_client * t, struct wt_metainfo * info);
  
  typedef void (*wt_started_hook)(struct wt_client *, void *);
  
  void WT_EXPORT wt_run(struct wt_client * t, wt_started_hook cb, void * user);

  void wt_tick(struct wt_client * t);
  
#ifdef __cplusplus
}
#endif
#endif
