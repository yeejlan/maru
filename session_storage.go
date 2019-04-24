package maru

type ISessionStorage interface {
	Load(sessionId string) (string, error)
	Save(sessionId string, data string) error
}

var (
	//SessionStorage singleton
	SessionStorage ISessionStorage
)