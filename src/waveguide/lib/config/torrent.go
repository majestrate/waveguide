package config

import (
	"math"
	"strconv"
	"waveguide/lib/config/parser"
)

type TorrentConfig struct {
	PieceSize  uint32
	TrackerURL string
}

func (c *TorrentConfig) Load(s *parser.Section) error {
	sz, err := strconv.ParseInt(s.ValueOf("torrent_piece_size"), 10, 64)
	if err != nil {
		return err
	}
	c.PieceSize = uint32(math.Pow(2, float64(sz)))
	c.TrackerURL = s.ValueOf("torrent_tracker")
	return nil
}
