package torrent

import (
	"io"
	"os"
	"path/filepath"
	"time"
	"waveguide/lib/config"
	"waveguide/lib/log"
	"waveguide/lib/torrent/metainfo"
)

func NewFactory(c *config.TorrentConfig) (*Factory, error) {
	return &Factory{
		AnnounceURL: c.TrackerURL,
		PieceLength: c.PieceSize,
	}, nil
}

type Factory struct {
	AnnounceURL string
	PieceLength uint32
}

func (f *Factory) MakeSingleWithWebseed(filename, webseed string, body io.Reader, out io.Writer) (err error) {
	t := metainfo.TorrentFile{
		Announce: f.AnnounceURL,
		Webseed:  webseed,
	}
	t.Info.PieceLength = f.PieceLength
	t.Info.Path = filepath.Base(filename)
	err = t.Info.BuildSingle(body)
	if err == nil {
		log.Debugf("created torrent for %s", filename)
		t.Created = time.Now().Unix()
		err = t.BEncode(out)
	}
	return
}

func (f *Factory) MakeSingle(filename string, body io.Reader, out io.Writer) error {
	t := metainfo.TorrentFile{
		Announce: f.AnnounceURL,
	}
	t.Info.PieceLength = f.PieceLength
	t.Info.Path = filepath.Base(filename)
	err := t.Info.BuildSingle(body)
	if err == nil {
		log.Debugf("created torrent for %s", filename)
		t.Created = time.Now().Unix()
		err = t.BEncode(out)
	}
	return err
}

func (f *Factory) Make(root, outfile string) error {
	t := metainfo.TorrentFile{
		Announce: f.AnnounceURL,
	}
	t.Info.PieceLength = f.PieceLength
	err := t.Info.BuildFromFilePath(root)
	if err == nil {
		t.Created = time.Now().Unix()
		f, e := os.Create(outfile)
		if e == nil {
			err = t.BEncode(f)
			f.Close()
		} else {
			err = e
		}
	}
	return err
}
