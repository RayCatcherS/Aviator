package web

import (
	"embed"
	"io/fs"
)

//go:embed static/* docs/*
var content embed.FS

// GetFS returns the filesystem containing the static web assets.
// It returns a subtree rooted at "static" folder if we want, or just the whole thing.
// Our served files are in "static" folder, so we might want to strip that prefix?
// The original python served "static" folder at root "/".
// Let's verify the structure.
// c:\Github\Aviator\aviator-go\web\static contains index.html
// So we want to serve "static" folder content at root.
func GetFS() (fs.FS, error) {
	return fs.Sub(content, "static")
}

// Docs?
func GetDocsFS() (fs.FS, error) {
	return fs.Sub(content, "docs")
}
