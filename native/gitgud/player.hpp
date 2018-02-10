#ifndef GITGUD_PLAYER_H
#define GITGUD_PLAYER_H

namespace gitgud
{
  struct PlayerImpl;
  struct Player
  {
    Player(const char * url);
    ~Player();

    const char * ApiURL() const;
    void EOS();
    void StopPlayback();
  private:
    PlayerImpl * m_Impl;
  };
}

#endif
