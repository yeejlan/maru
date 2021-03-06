package maru

import (
	"log"
	"strconv"
	"encoding/json"
)

type Session struct {
	data map[string]interface{}
	changed bool
	sessionId string
}

//new session
func NewSession() *Session {
	return &Session{
		data: make(map[string]interface{}),
		changed: false,
		sessionId: "",
	}
}

//set
func (this *Session) Set(key string, val interface{}) {
	this.change()
	this.data[key] = val
}

//get value as interface{}
func (this *Session) Get(key string) interface{} {
	return this.data[key]
}

//get value as string
func (this *Session) GetString(key string, defaultVal ...string) string {
	if val, ok := this.data[key]; ok {
		return val.(string)
	}
	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return ""
}

//get int value
func (this *Session) GetInt(key string, defaultVal ...int) int {
	if val, ok := this.data[key]; ok {
		valStr := val.(json.Number).String()
		if intval, err := strconv.Atoi(valStr); err == nil {
			return intval
		}
	}

	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return 0
}

//get bool value
func (this *Session) GetBool(key string, defaultVal ...bool) bool {
	if val, ok := this.data[key]; ok {
		return val.(bool)
	}

	if(len(defaultVal) > 0) {
		return defaultVal[0]
	}
	return false
}

//delete
func (this *Session) Delete(key string) {
	this.change()
	delete(this.data, key)
}

//destroy this session
func (this *Session) Destroy() {
	this.change()
	this.data = make(map[string]interface{})
	this.Save()
}

//get session id
func (this *Session) Id() string {
	return this.sessionId
}

//set session id
func (this *Session) SetId(sid string) {
	this.sessionId = sid
}

//touch a session, make a session reset it's lifetime
func (this *Session) Touch() {
	this.change()
}

func (this *Session) change() {
	this.changed = true
}

//sava session
func (this *Session) Save() {
	if(!this.changed){
		return
	}
	if(this.sessionId == ""){
		return
	}
	if(SessionStorage != nil){
		val, err := JsonEncode(this.data)
		if err!= nil {
			log.Print("session save:json encode error:", err)
			return
		}
		err = SessionStorage.Save(this.sessionId, string(val[:]))
		if err!= nil {
			log.Print("session save error:", err)
		}
	}
}

func (this *Session) Load() {
	if(this.sessionId == ""){
		return
	}
	if(SessionStorage != nil){
		val, err := SessionStorage.Load(this.sessionId)
		if err != nil {
			log.Print("session load error:", err)
			return
		}
		this.data, err = JsonDecodeToMap(val)
		if err != nil {
			log.Print("session load: json decode error:", err)
			return
		}
	}
}