#include "log.h"


struct wt_logger
{
  FILE * out;
};

static void _wt_log(struct wt_logger * log, const char * tag, const char * msg)
{
  fprintf(log->out, "[%s] %s\n", tag, msg);
}

struct wt_logger wt_log;

bool wt_log_init()
{
  wt_log.out = stdout;
  return true;
}

void wt_info(const char * msg)
{
  _wt_log(&wt_log, "NFO", msg);
}

void wt_warn(const char * msg)
{
  _wt_log(&wt_log, "WRN", msg);
}

void wt_error(const char * msg)
{
  _wt_log(&wt_log, "ERR", msg);
}
