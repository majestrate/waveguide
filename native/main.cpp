#include "wasmtorrent/wasmtorrent.h"
#include "gitgud/player.hpp"
#include <cassert>

static void end_of_stream(void * user)
{
  wt_info("end_of_stream");
  gitgud::Player * player = static_cast<gitgud::Player *>(user);
  player->EOS();
}

static void got_stream(wt_stream * stream, void * user)
{
  wt_info("got_stream");
  gitgud::Player * player = static_cast<gitgud::Player *>(user);
  wt_stream_ev_listener l;
  l.eos = &end_of_stream;
  wt_stream_add_event_listener(stream, l, player);
}

static void on_started(wt_client * wt, void * user)
{
  wt_info("on_started");
  gitgud::Player * player = static_cast<gitgud::Player *>(user);
  wt_client_add_stream(wt, player->ApiURL(), &got_stream, player);
}

int main(int argc, char * argv[])
{
  assert(wt_log_init());

  /** default url */
  const char * url = "https://gitgud.tv/wg-api/v1/stream?key=4452";
  if(argc == 2)
  {
    url = argv[1];
  }
  
  wt_info("WASMTorrent starting up");
  wt_info(url);
  wt_client * wt = nullptr;
  wt_client_init(&wt);

  gitgud::Player * player = new gitgud::Player(url);
  
  wt_run(wt, &on_started, player);
  wt_client_free(&wt);
  wt_info("WASMTorrent done");

  delete player;
  
  return 0;
}
