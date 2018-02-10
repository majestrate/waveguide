#ifndef _WT_STREAM_H
#define _WT_STREAM_H
#include "wasm.h"
#include "torrent.h"
#include "metainfo.h"

#ifdef __cplusplus
extern "C" {
#endif

  /** a wt_stream is a metainfo input stream that obtains in order metainfos for streaming media */
  struct wt_stream;

  /* all of these are not to be exported */
  void wt_stream_init(struct wt_stream ** s, struct wt_client * cl, const char * url);
  void wt_stream_free(struct wt_stream ** s);
  bool wt_stream_should_poll(struct wt_stream * s);
  void wt_stream_poll(struct wt_stream * s);

  
  /** wt_stream visitor hook */
  typedef void (*wt_stream_visitor)(struct wt_stream *, void *);

  /** create a new stream using api endpoint at api_url, visit the newly created stream in a callback v */
  void WT_EXPORT wt_client_add_stream(struct wt_client * cl, const char * api_url, wt_stream_visitor v, void * user);

  struct wt_stream_ev_listener
  {
    /** called when the input stream is starting up */
    void (*start)(void * user);
    /** we got a new metainfo for segment N */
    void (*metainfo)(struct wt_metainfo * info, size_t n, void * user);
    /** no more metainfos */
    void (*eos)(void * user);
  };

  /** add event listener for this stream */
  void WT_EXPORT wt_stream_add_event_listener(struct wt_stream * stream, struct wt_stream_ev_listener l, void * user);

  void wt_stream_close(struct wt_stream * stream);
  
#ifdef __cplusplus
}
#endif
#endif
