package torrent

import (
	"io"
	"os"
	"path/filepath"
	"time"
	"waveguide/lib/torrent/metainfo"
)

func NewFactory() (*Factory, error) {
	return &Factory{
		AnnounceURL: "wss://tracker.btorrent.xyz",
		/* 1 MB */
		PieceLength: 1048576,
	}, nil
}

type Factory struct {
	AnnounceURL string
	PieceLength uint32
}

func (f *Factory) MakeSingle(fname string, body io.Reader, out io.Writer) error {
	t := metainfo.TorrentFile{
		Announce: f.AnnounceURL,
	}
	t.Info.PieceLength = f.PieceLength
	t.Info.Path = filepath.Base(fname)
	err := t.Info.BuildSingle(body)
	if err == nil {
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
