package maru

import (
	"encoding/json"
)

var(
	//Set with go build -ldflags "-X github.com/yeejlan/maru.BuildDir=xxx"
	BuildDir string

	jsonEncode = json.Marshal
	jsonDecode = json.Unmarshal
)

func init() {
	if BuildDir == "" {
		println("empty")
	}else{
		println("not empty")
	}
}