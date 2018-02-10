#ifndef _WASM_H
#define _WASM_H
#ifdef WASM
#include <emscripten/emscripten.h>
#define WT_EXPORT EMSCRIPTEN_KEEPALIVE
#else
#define WT_EXPORT extern
#endif
#endif
