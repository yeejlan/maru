package maru

import (
	"fmt"
	"log"
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
	config *Config
	appName string
	cookieDomain string
	sessionName string
}

//create new App
func NewApp(envString string, appName string) *App {
	env := strToEnv(envString)
	envStr := envStrMapping[env]
	return &App{
		isInit: false,
		env: env,
		envString: envStr,
		appName: appName,
	}
}

//implement string interface
func (this *App) String() string {
	return fmt.Sprintf("App[name=%s, env=%s]", this.appName, this.envString)
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
	
	//set log flag
	//log.SetFlags(log.LstdFlags | log.Llongfile)

	//load config file
	configFile := fmt.Sprintf("config/%s/%s.ini", this.envString, this.appName)
	config := NewConfig(configFile)
	this.config = config

	//set log path
	logpath := config.Get("log.path", "logs")
	InitLog(logpath)

	if(this.env < DEVELOPMENT) {
		log.SetOutput(NewLogger(logpath, this.appName))
	}

	this.cookieDomain = config.Get("cookie.domain")
	this.sessionName = config.Get("session.name")

	//set Jet template to dev mode
	if(this.env == DEVELOPMENT) {
		JetSet.SetDevelopmentMode(true)
	}

	log.Printf("%s starting with %s", this, this.config.ConfigFile)
}

//get app config
func (this *App) Config() *Config {
	this.checkInit()
	return this.config
}

//get app env
func (this *App) Env() int {
	return this.env
}

//get app env string
func (this *App) EnvString() string {
	return this.envString
}

//get app name
func (this *App) Name() string {
	return this.appName
}

//get cookie domain
func (this *App) CookieDomain() string {
	return this.cookieDomain
}

//get session name
func (this *App) SessionName() string {
	return this.sessionName
}

func (this *App) checkInit() {
	if(!this.isInit) {
		log.Fatal(`Please call "App.Init()" first`)
	}
}