package maru

import (
	"github.com/go-redis/redis"
	"time"
)

type SessionStorageRedis struct {
	client *redis.Client
	expire time.Duration
}

func NewSessionStorageRedis(app *App) *SessionStorageRedis {
	sessionExpire := app.Config().GetInt("session.expire.seconds", 3600)
	storageProvider := app.Config().Get("session.storage.provider")
	_=storageProvider
	return &SessionStorageRedis{
		expire: time.Duration(sessionExpire) * time.Second,
	}
}

func (this *SessionStorageRedis) load(sessionId string) string {
	val, _ := this.client.Get("key").Result()
	return val
}

func (this *SessionStorageRedis) save(sessionId string, data string) {
	this.client.Set("key", "value", this.expire)
}
