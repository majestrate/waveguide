package templates

import (
	"time"
)

func FormatDate(t time.Time) string {
	return t.Format(time.ANSIC)
}
