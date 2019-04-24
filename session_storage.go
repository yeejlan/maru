package maru

type ISessionStorage interface {
	load(sessionId string) string
	save(sessionId string, data string)
}

var (
	//SessionStorage singleton
	SessionStorage ISessionStorage
)