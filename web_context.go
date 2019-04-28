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
type WebContext struct {
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
	Session Session
	//jet template vars
	View jet.VarMap

	sessionName string
	cookieDomain string
}

func newWebContext(app *App, w http.ResponseWriter, req *http.Request) *WebContext {
	sessionName := app.Config().Get("session.name")
	cookieDomain := app.Config().Get("cookie.domain")
	return &WebContext{
		App: app,
		Req: req,
		W: w,
		sessionName: sessionName,
		cookieDomain: cookieDomain,
		View: make(jet.VarMap),
	}
}

//set cookie
func (this *WebContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(this.W, cookie)
}

//new session
func (this *WebContext) NewSession(){
	if SessionStorage == nil{
		return
	}
	this.Session.Destroy()
	this.Session.SetId(GetUniqueId())

	cookie := &http.Cookie{
		Name: this.sessionName,
		Value: this.Session.Id(),
		Domain: this.cookieDomain,
	}
	this.SetCookie(cookie)
}

//load session
func (this *WebContext) LoadSession(){
	if SessionStorage == nil{
		return
	}
	sessionId := this.Cookie.Get(this.sessionName)
	if sessionId == ""{
		this.NewSession()
	}else{
		this.Session.SetId(sessionId)
		this.Session.Load()
	}
}

//abort a request
func (this *WebContext) Abort(status int, body string) {
	this.W.WriteHeader(status)
	this.W.Write([]byte(body))
}

//redirect a request
func (this *WebContext) Redirect(url string) {
	this.W.Header().Set("Location", url)
	this.W.WriteHeader(302)
}

//exit current request
func (this *WebContext) Exit() {
	panic(internalRequestExit{})
}

type internalRequestExit struct{}

//render template
func (this *WebContext) Render(templateName string) {
	t, err := JetSet.GetTemplate(templateName)
	if err != nil {
		panic(err)
	}
	if err = t.Execute(this.W, this.View, nil); err != nil {
		panic(err)
	}
}

//render to string
func (this *WebContext) RenderToString(templateName string) (string, error) {
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
