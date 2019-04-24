package maru

import (
	"fmt"
	"regexp"
	"github.com/go-redis/redis"
	"log"
)

type ResourceLoader struct {
	app *App
	env string
}

func NewResourceLoader(app *App) *ResourceLoader {
	envString := app.EnvString()
	return &ResourceLoader{
		app: app,
		env: envString,
	}
}
//autoload redis and databse resources based on config
func (this *ResourceLoader) Autoload() {
	this.AutoloadRedis()
	this.LoadSessionStorage()
}

//autoload redis instances based on config
func (this *ResourceLoader) AutoloadRedis() {
	configFile := fmt.Sprintf("config/%s/redis.ini", this.env)
	config := NewConfig(configFile)
	configMatcher := regexp.MustCompile("^redis\\.([_a-zA-Z0-9]+)\\.host")
	for key, _ := range config.GetMap() {
		if(!configMatcher.MatchString(key)){
			continue;
		}
		//found one
		redisName := key[0 : len(key) - len(".host")]
		if(config.GetBool(redisName + ".autoload")) {
			//autoload redis
			this.LoadRedis(config, redisName)
		}
	}
}

//load one redis instance based on config
func (this *ResourceLoader) LoadRedis(config *Config, redisName string) *redis.Client {
	host := config.Get(redisName + ".host", "localhost")
	port := config.GetInt(redisName + ".port", 6379)
	database := config.GetInt(redisName + ".database", 0)
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		Password: "", // no password set
		DB: database,
	})

	Registry.set(redisName, client)
	return client
}

func (this *ResourceLoader) LoadSessionStorage() {
	sessionEnable := this.app.Config().GetBool("session.enable")
	storageName := this.app.Config().Get("session.storage")
	if(!sessionEnable) {
		log.Print("Session is NOT enabled.")
		return
	}

	//set session storage
	switch storageName {
		case "redis":
			SessionStorage = NewSessionStorageRedis(this.app)
		default:
			panic("Session storage not supported: " + storageName)
	}
}