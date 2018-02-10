namespace wasmtorrent
{
  template<typename CharA_t, typename CharB_t>
  bool StrEq(CharA_t * a, CharB_t * b, size_t sz)
  {
    size_t idx = 0;
    while(idx < sz)
    {
      if (a[idx] == 0 && b[idx] == 0) return true;
      if (a[idx] != b[idx]) return false;
      ++idx;
    }
    return false;
  }
}
