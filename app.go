package maru

import (
	"fmt"
)

const (
	PRODUCTION int = iota + 1
	STAGING
	TESTING
	DEVELOPMENT
)

var (
	envStrMapping = map[int]string {
		PRODUCTION : "production",
		STAGING : "staging",
		TESTING : "testing",
		DEVELOPMENT: "development",
	}
)

type App struct {
	isInit bool
	env int
	envString string
	configFile string
	config string
	name string
}

//create new App
func NewApp(envString string, appName string) *App {
	env := strToEnv(envString)
	envStr := envStrMapping[env]
	return &App{
		isInit: false,
		env: env,
		envString: envStr,
		name: appName,
	}
}

//implement string interface
func (this *App) String() string {
	return fmt.Sprintf("App[name=%s, env=%s]", this.name, this.envString)
}

//get App.env
func (this *App) Env() int {
	return this.env
}

//get App.envString
func (this *App) EnvString() string {
	return this.envString
}

func strToEnv(envString string) int {
	env := PRODUCTION
	for k, v := range envStrMapping {
		if(v == envString){
			env = k
			break
		}
	}
	return env
}

//App initialize
func (this *App) Init() {
	this.isInit = true
	println("app is init... %s", this)
}