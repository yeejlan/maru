package maru

import (
	"os"
	"bufio"
	"log"
	"strings"
	"fmt"
)

type Config struct {
	ConfigFile string
	StringMap
}

//load new config
func NewConfig(filePath string) *Config {
	file, err := os.Open(filePath)
		if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if(len(kv) == 2) {
			k := strings.TrimSpace(kv[0])
			v := strings.TrimSpace(kv[1])
			data[k] = v
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &Config {
		ConfigFile: filePath,
		StringMap: data,
	}
}

//implement string interface
func (this *Config) String() string {
	return fmt.Sprintf("Config{ConfigFile: %s, StringMap: %#v}", this.ConfigFile, this.StringMap)
}

//get config map
func (this *Config) GetMap() map[string]string {
	return this.StringMap
}