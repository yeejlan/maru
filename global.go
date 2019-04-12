package maru

import(
	"fmt"
)

var(
	println = fmt.Println
	
	//Set with go build -ldflags "-X github.com/yeejlan/maru.BuildDir=xxx"
	BuildDir string
)