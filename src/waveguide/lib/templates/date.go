package templates

import (
	"time"
)

func FormatDate(t int64) string {
	return time.Unix(t, 0).Format(time.ANSIC)
}
