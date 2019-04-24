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
	
	redisClient := Registry.get(storageProvider)
	if redisClient == nil {
		panic("NewSessionStorageRedis() "+storageProvider + " not found")
	}
	return &SessionStorageRedis{
		expire: time.Duration(sessionExpire) * time.Second,
		client: redisClient.(*redis.Client),
	}
}

func (this *SessionStorageRedis) Load(sessionId string) (string, error) {
	val, err := this.client.Get(sessionId).Result()
	if err == redis.Nil { //key not exist
		return "", nil
	}
	return val, err
}

func (this *SessionStorageRedis) Save(sessionId string, data string) error {
	err := this.client.Set(sessionId, data, this.expire).Err()
	return err
}
