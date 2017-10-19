package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (w *Worker) TempFileName(ext string) string {
	return filepath.Join(w.TempDir, fmt.Sprintf("file-%d%s", time.Now().UnixNano(), ext))
}

func (w *Worker) AcquireTempFile(ext string) (f *os.File, err error) {
	return os.Create(w.TempFileName(ext))
}
