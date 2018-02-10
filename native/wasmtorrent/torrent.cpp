#include "torrent.h"
#include "stream.h"
#include "timer.h"
#include "mainloop.h"
#include <chrono>
#include <thread>
#include <list>

struct wt_client
{

  wt_client() 
  {
  }

  ~wt_client()
  {
  }

  void Tick()
  {
    for ( auto & itr : streams)
    {
      if(wt_stream_should_poll(itr)) wt_stream_poll(itr);
    }
  }

  void Open()
  {
  }

  bool Continue()
  {
    return true;
  }
  
  void Close()
  {
    for (auto & itr : streams)
    {
      wt_stream * stream = itr;
      wt_stream_close(stream);
      wt_stream_free(&stream);
    }
    streams.clear();
  }

  wt_stream * AddStream(const char * url)
  {
    wt_stream * stream = nullptr;
    wt_stream_init(&stream, this, url);
    if(stream)
      streams.push_back(stream);
    return stream;
  }

  std::list<wt_stream *> streams;
};

/*
static void wt_poll(void * user)
{
  wt_client * cl = static_cast<wt_client*>(user);
  cl->Tick();
  if (cl->Continue())
    wt_set_timeout(&wt_poll, user, 100);
}
*/

extern "C" {

  void wt_client_init(struct wt_client ** cl)
  {
    *cl = new wt_client;
  }


  void wt_client_free(struct wt_client ** cl)
  {
    if(*cl)
    {
      delete *cl;
      *cl = nullptr;
    }
  }

  void wt_client_add_stream(struct wt_client * cl, const char * api_url, wt_stream_visitor v, void * user)
  {
    wt_stream * stream = cl->AddStream(api_url);
    v(stream, user);
  }
  
  void wt_run(struct wt_client * cl, wt_started_hook cb, void * user)
  {
    cl->Open();
    cb(cl, user);
    wt_mainloop(cl);
    cl->Close();
  }

  void wt_tick(struct wt_client * t)
  {
    t->Tick();
  }
}
