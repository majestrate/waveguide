package metainfo

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/zeebo/bencode"
	"io"
	"os"
	"path/filepath"
)

type FilePath []string

// get filepath
func (f FilePath) FilePath() string {
	return filepath.Join(f...)
}

/** open file using base path */
func (f FilePath) Open(base string) (*os.File, error) {
	return os.OpenFile(filepath.Join(base, f.FilePath()), os.O_RDWR|os.O_CREATE, 0600)
}

type FileInfo struct {
	// length of file
	Length uint64 `bencode:"length"`
	// relative path of file
	Path FilePath `bencode:"path"`
	// md5sum
	Sum []byte `bencode:"md5sum,omitempty"`
}

// info section of torrent file
type Info struct {
	// length of pices in bytes
	PieceLength uint32 `bencode:"piece length"`
	// piece data
	Pieces []byte `bencode:"pieces"`
	// name of root file
	Path string `bencode:"name"`
	// file metadata
	Files []FileInfo `bencode:"files,omitempty"`
	// private torrent
	Private *uint64 `bencode:"private,omitempty"`
	// length of file in signle file mode
	Length uint64 `bencode:"length,omitempty"`
	// md5sum
	Sum []byte `bencode:"md5sum,omitempty"`
}

// get fileinfos from this info section
func (i Info) GetFiles() (infos []FileInfo) {
	if i.Length > 0 {
		infos = append(infos, FileInfo{
			Length: i.Length,
			Path:   FilePath([]string{i.Path}),
			Sum:    i.Sum,
		})
	} else {
		infos = append(infos, i.Files...)
	}
	return
}

func (i Info) NumPieces() uint32 {
	return uint32(len(i.Pieces) / 20)
}

var ErrNoMultifile = errors.New("no multifile support")

func (i *Info) BuildSingle(r io.Reader) error {
	piece := make([]byte, i.PieceLength)
	for {
		n, err := io.ReadFull(r, piece)
		if n > 0 {
			i.Length += uint64(n)
		}
		if err == nil {
			h := sha1.Sum(piece)
			i.Pieces = append(i.Pieces, h[:]...)
		} else if err == io.EOF {
			h := sha1.Sum(piece[:n])
			i.Pieces = append(i.Pieces, h[:]...)
			err = nil
			break
		} else {
			return err
		}
	}
	return nil
}

func (i *Info) BuildFromFilePath(fpath string) error {
	stat, err := os.Stat(fpath)
	if err == nil {
		i.Path = fpath
		if stat.IsDir() {
			err = ErrNoMultifile
		} else {
			i.Length = uint64(stat.Size())
			f, err := os.Open(fpath)
			if err != nil {
				return err
			}
			defer f.Close()
			piece := make([]byte, i.PieceLength)
			l := i.Length
			for l > uint64(i.PieceLength) {
				_, err = io.ReadFull(f, piece)
				if err != nil {
					return err
				}
				h := sha1.Sum(piece)
				i.Pieces = append(i.Pieces, h[:]...)
				l -= uint64(i.PieceLength)
			}
			piece = make([]byte, l)
			_, err = io.ReadFull(f, piece)
			if err != nil {
				return err
			}
			h := sha1.Sum(piece)
			i.Pieces = append(i.Pieces, h[:]...)
			return nil
		}
	}
	return err
}

// a torrent file
type TorrentFile struct {
	Info         Info       `bencode:"info"`
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list"`
	Created      int64      `bencode:"created"`
	Comment      []byte     `bencode:"comment"`
	CreatedBy    []byte     `bencode:"created by"`
	Encoding     []byte     `bencode:"encoding"`
}

func (tf *TorrentFile) LengthOfPiece(idx uint32) (l uint32) {
	i := tf.Info
	np := i.NumPieces()
	if np == idx+1 {
		sz := tf.TotalSize()
		l64 := uint64(i.PieceLength) - ((uint64(np) * uint64(i.PieceLength)) - sz)
		l = uint32(l64)
	} else {
		l = i.PieceLength
	}
	return
}

// get total size of files from torrent info section
func (tf *TorrentFile) TotalSize() uint64 {
	if tf.IsSingleFile() {
		return tf.Info.Length
	}
	total := uint64(0)
	for _, f := range tf.Info.Files {
		total += f.Length
	}
	return total
}

func (tf *TorrentFile) GetAllAnnounceURLS() (l []string) {
	if len(tf.Announce) > 0 {
		l = append(l, tf.Announce)
	}
	for _, al := range tf.AnnounceList {
		for _, a := range al {
			if len(a) > 0 {
				l = append(l, a)
			}
		}
	}
	return
}

func (tf *TorrentFile) TorrentName() string {
	return tf.Info.Path
}

type Infohash [20]byte

func (ih Infohash) Hex() string {
	return hex.EncodeToString(ih[:])
}

// calculate infohash
func (tf *TorrentFile) Infohash() (ih Infohash) {
	s := sha1.New()
	enc := bencode.NewEncoder(s)
	enc.Encode(&tf.Info)
	d := s.Sum(nil)
	copy(ih[:], d[:])
	return
}

// return true if this torrent is for a single file
func (tf *TorrentFile) IsSingleFile() bool {
	return tf.Info.Length > 0
}

// bencode this file via an io.Writer
func (tf *TorrentFile) BEncode(w io.Writer) (err error) {
	enc := bencode.NewEncoder(w)
	err = enc.Encode(tf)
	return
}

// load from an io.Reader
func (tf *TorrentFile) BDecode(r io.Reader) (err error) {
	dec := bencode.NewDecoder(r)
	err = dec.Decode(tf)
	return
}

// IsPrivate returns true if this torrent is a private torrent
func (tf *TorrentFile) IsPrivate() bool {
	return tf.Info.Private != nil && *tf.Info.Private > 0
}
