package maru

import (
	"os"
	"bufio"
	"log"
	"strings"
	"fmt"
	"strconv"
)

type Config struct {
	FilePath string
	data map[string]string
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
		kv := strings.Split(line, "=")
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
		FilePath: filePath,
		data: data,
	}
}

//implement string interface
func (this *Config) String() string {
	return fmt.Sprintf("%#v", *this)
}

//get config map
func (this *Config) GetMap() map[string]string {
	return this.data
}

func (this *Config) Get(pathStr string, defaultVal ...string) string {
	if val, ok := this.data[pathStr]; ok {
		return val
	}
	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return ""
}

func (this *Config) GetString(pathStr string, defaultVal ...string) string {
	return this.Get(pathStr, defaultVal...)
}

func (this *Config) GetInt(pathStr string, defaultVal ...int) int {
	if val, ok := this.data[pathStr]; ok {
		if intval, err := strconv.Atoi(val); err == nil {
			return intval
		}
	}

	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return 0
}

func (this *Config) GetBool(pathStr string, defaultVal ...bool) bool {
	if val, ok := this.data[pathStr]; ok {
		if boolval, err := strconv.ParseBool(val); err == nil {
			return boolval
		}
	}

	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return false
}