#ifndef WASM_TIME_HPP
#define WASM_TIME_HPP

#include <chrono>

namespace wasmtorrent
{
  uint32_t NowSeconds()
  {
    return std::chrono::duration_cast<std::chrono::seconds>(
      std::chrono::steady_clock::now().time_since_epoch()).count();
  }
}

#endif
