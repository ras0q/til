package playground

import (
	"embed"
	"io/fs"
)

//go:embed client/web/dist/*
var publicFS embed.FS

var PublicFS fs.FS

func init() {
	_publicFS, err := fs.Sub(publicFS, "client/web/dist")
	if err != nil {
		panic(err)
	}
	PublicFS = _publicFS
}
