package maru

import(
	"strings"
)

var(
	//Set with go build -ldflags "-X github.com/yeejlan/maru.BuildDir=xxx"
	BuildDir string
)

func init() {
	BuildDir = strings.ReplaceAll(BuildDir, "\\", "/")
}