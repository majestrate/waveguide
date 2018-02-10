#include "stream.h"
#include "log.h"
#include "wget.h"
#include "time.hpp"
#include "str.hpp"
#include <list>

namespace wasmtorrent
{
  struct ev_listener
  {
    void * user;
    wt_stream_ev_listener hooks;

    void End() const
    {
      hooks.eos(user);
    }

    void MetaInfo(wt_metainfo * info, size_t n) const
    {
      hooks.metainfo(info, n, user);
    }

    void Start() const
    {
      hooks.start(user);
    }
    
  };
}

void wt_stream_handle_fetch_url_data(void * stream, void * data, int sz);

void wt_stream_handle_fetch_url_error(void * stream);

void wt_stream_handle_fetch_metainfo_data(void * stream, void * data, int sz);

void wt_stream_handle_fetch_metainfo_error(void * stream);



struct wt_stream
{
  wt_stream(wt_client * cl, const char * url) :
    m_Client(cl),
    m_FetchingMetaInfo(nullptr),
    m_URL(url),
    m_PollInterval(5),
    m_LastPoll(0)
  {
  };

  void AddListener(wt_stream_ev_listener l, void * user)
  {
    listeners.push_back(wasmtorrent::ev_listener{user, l});
  }

  void Close()
  {
    for(const auto & l : listeners)
    {
      l.End();
    }
    listeners.clear();
  }
  
  ~wt_stream()
  {
  }

  bool ShouldPoll()
  {
    auto now = wasmtorrent::NowSeconds();
    return now - m_LastPoll >= m_PollInterval;
  }

  void Poll()
  {
    wt_wget(m_URL, this, wt_stream_handle_fetch_url_data,  wt_stream_handle_fetch_url_error);
  }

  void FetchedURLData(void * data, int sz)
  {
    if(sz >= sizeof(m_LastURL))
    {
      /** url too long ? */
      wt_warn("fetched url was too long");
      return;
    }
    /* got url data */
    char * nextURL = static_cast<char *>(data);
    nextURL[sz] = 0;
    if(!wasmtorrent::StrEq(nextURL, m_LastURL, sizeof(m_LastURL)))
    {
      /* new url */
      m_LastPoll = wasmtorrent::NowSeconds();
      memcpy(m_LastURL, nextURL, sz);
      wt_wget(m_LastURL, this, &wt_stream_handle_fetch_metainfo_data, &wt_stream_handle_fetch_metainfo_error);
    }
  }

  void FetchedMetaInfoData(void * data, int sz)
  {
    wt_metainfo * info = nullptr;
    wt_metainfo_init(&info);
    if(wt_metainfo_loadbuffer(info, data, sz))
    {
      /* parse success */
      uint8_t ih[20];
      wt_metainfo_calc_infohash(info, ih);
      if(!wt_client_has_torrent(m_Client, ih))
      {
        wt_info("adding new torrent");
        wt_client_add_torrent(m_Client, info);
      }
    }
    else
    {
      /* parse fail */
      wt_error("metainfo is invalid format");
    }
    wt_metainfo_free(&info);
  }

  void FetchError()
  {
    wt_error("fetch error");
  }

  std::list<wasmtorrent::ev_listener> listeners;
  wt_client * m_Client;
  wt_metainfo * m_FetchingMetaInfo;
  const char * m_URL;
  uint32_t m_PollInterval;
  uint32_t m_LastPoll;
  char m_LastURL[512];
};

void wt_stream_handle_fetch_url_data(void * stream, void * data, int sz)
{
  static_cast<wt_stream *>(stream)->FetchedURLData(data, sz);
}

void wt_stream_handle_fetch_url_error(void * stream)
{
  static_cast<wt_stream *>(stream)->FetchError();
}

void wt_stream_handle_fetch_metainfo_data(void * stream, void * data, int sz)
{
  static_cast<wt_stream *>(stream)->FetchedMetaInfoData(data, sz);
}

void wt_stream_handle_fetch_metainfo_error(void * stream)
{
  static_cast<wt_stream *>(stream)->FetchError();
}


extern "C" {

  void wt_stream_init(struct wt_stream ** s, struct wt_client * cl, const char * url)
  {
    *s = new wt_stream(cl, url);
  }

  void wt_stream_free(struct wt_stream ** s)
  {
    if (*s) delete *s;
    *s = nullptr;
  }
  
  void wt_stream_add_event_listener(struct wt_stream * stream, struct wt_stream_ev_listener l, void * user)
  {
    stream->AddListener(l, user);
  }

  void wt_stream_close(struct wt_stream * stream)
  {
    stream->Close();
  }

  bool wt_stream_should_poll(struct wt_stream * stream)
  {
    return stream->ShouldPoll();
  }

  void wt_stream_poll(struct wt_stream * stream)
  {
    stream->Poll();
  }

}
