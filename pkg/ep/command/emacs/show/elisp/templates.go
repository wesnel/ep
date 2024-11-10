package elisp

import (
	"embed"
	"text/template"
)

//go:embed *.el
var fs embed.FS

var Show = template.Must(template.ParseFS(fs, "show.el"))
