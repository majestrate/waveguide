package worker

import (
	"os"
	"waveguide/lib/util"
)

func (w *Worker) TempFileName(ext string) string {
	return util.TempFileName(w.TempDir, ext)
}

func (w *Worker) AcquireTempFile(ext string) (f *os.File, err error) {
	return os.Create(w.TempFileName(ext))
}
