package maru

import (
	"net/http"
	"bytes"
	"github.com/CloudyKit/jet"
)

var(
	//jet template set
	JetSet = jet.NewHTMLSet("./templates")
)

//web context
type Ctx struct {
	App *App
	Req *http.Request
	W http.ResponseWriter
	//current controller
	Controller string
	//current action
	Action string
	//internal server error
	Error interface{}
	//param map
	Param StringMap
	//cookie map
	Cookie StringMap
	//session instance
	Session *Session
	//jet template vars
	View jet.VarMap
}

func newCtx(app *App, w http.ResponseWriter, req *http.Request) *Ctx {
	return &Ctx{
		App: app,
		Req: req,
		W: w,
		View: make(jet.VarMap),
		Session: NewSession(),
	}
}

//set cookie
func (this *Ctx) SetCookie(cookie *http.Cookie) {
	http.SetCookie(this.W, cookie)
}

//new session
func (this *Ctx) NewSession(){
	if SessionStorage == nil{
		return
	}
	this.Session.Destroy()
	this.Session.SetId(GetUniqueId())

	cookie := &http.Cookie{
		Name: this.App.SessionName(),
		Value: this.Session.Id(),
		Domain: this.App.CookieDomain(),
	}
	this.SetCookie(cookie)
}

//load session
func (this *Ctx) LoadSession(){
	if SessionStorage == nil{
		return
	}
	sessionId := this.Cookie.Get(this.App.SessionName())
	if sessionId == ""{
		this.NewSession()
	}else{
		this.Session.SetId(sessionId)
		this.Session.Load()
	}
}

//abort a request
func (this *Ctx) Abort(status int, body string) {
	this.W.WriteHeader(status)
	this.W.Write([]byte(body))
}

//redirect a request
func (this *Ctx) Redirect(url string) {
	this.W.Header().Set("Location", url)
	this.W.WriteHeader(302)
}

//exit current request
func (this *Ctx) Exit() {
	panic(internalRequestExit{})
}

type internalRequestExit struct{}

//render template
func (this *Ctx) Render(templateName string) {
	t, err := JetSet.GetTemplate(templateName)
	if err != nil {
		panic(err)
	}
	if err = t.Execute(this.W, this.View, nil); err != nil {
		panic(err)
	}
}

//render to string
func (this *Ctx) RenderToString(templateName string) (string, error) {
	t, err := JetSet.GetTemplate(templateName)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	if err = t.Execute(&buffer, this.View, nil); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
