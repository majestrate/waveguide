package util

import (
	"fmt"
	"path/filepath"
	"time"
)

func TempFileName(tmpdir, ext string) string {
	f, _ := filepath.Abs(filepath.Join(tmpdir, fmt.Sprintf("file-%d%s", time.Now().UnixNano(), ext)))
	return f
}
