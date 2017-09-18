package templates

import (
	"html/template"
)

func Funcs() template.FuncMap {
	return template.FuncMap{
		"FormatDate": FormatDate,
	}
}
