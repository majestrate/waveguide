#include "player.hpp"

namespace gitgud
{


  struct PlayerImpl
  {
    PlayerImpl(const char * url) :
      URL(url)
    {
    }

    ~PlayerImpl()
    {
    }

    void StopPlayer()
    {
    }

    void EOS()
    {
    }
    
    const char * URL;
  };
  
  Player::Player(const char * url) :
    m_Impl(new PlayerImpl(url))
  {
  }
    
  Player::~Player()
  {
    delete m_Impl;
  }

  const char * Player::ApiURL() const
  {
    return m_Impl->URL;
  }
  void Player::EOS()
  {
    m_Impl->StopPlayer();
    m_Impl->EOS();
  }

  void Player::StopPlayback()
  {
    m_Impl->StopPlayer();
  }
}
