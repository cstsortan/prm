package web

import (
	"embed"
	"io/fs"
)

//go:embed all:static
var staticFS embed.FS

// staticFiles returns the embedded filesystem rooted at "static/".
func staticFiles() (fs.FS, error) {
	return fs.Sub(staticFS, "static")
}
