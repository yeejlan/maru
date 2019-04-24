package maru

import (
	"github.com/json-iterator/go"
)

var(
	//Set with go build -ldflags "-X github.com/yeejlan/maru.BuildDir=xxx"
	BuildDir string

	json = jsoniter.ConfigCompatibleWithStandardLibrary
	JsonEncode = json.Marshal
	JsonDecode = json.Unmarshal
)

func init() {
	if BuildDir == "" {
		println("empty")
	}else{
		println("not empty")
	}
}