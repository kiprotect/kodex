package web

import (
	"embed"
	"io/fs"
)

//go:embed static/build
var StaticBuild embed.FS

//go:embed static/assets
var StaticAssets embed.FS

var StaticBuildPrefix, _ = fs.Sub(StaticBuild, "static/build")
var StaticAssetsPrefix, _ = fs.Sub(StaticAssets, "static/assets")

var StaticFiles = []fs.FS{StaticBuildPrefix, StaticAssetsPrefix}
