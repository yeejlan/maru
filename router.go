package maru

import (
	"net/http"
	"fmt"
	"log"
	"regexp"
	"os"
	"path"
	"strings"
)

//action info
type actionPair struct {
	//controller instance, ex. HomeController{}
	I interface{}
	//action, ex. "Index"
	A string 
}

//implement string interface
func (this actionPair) String() string {
	return fmt.Sprintf("ActionPair{I: %T, A: %s}", this.I, this.A)
}

//action storage, for example: 
//ActionMap["home/index"] = ActionPair{I: HomeController{}, A: "Index",}
var actionMap = make(map[string]actionPair)

//url rewite route storage
var routes = make([]rewriteRoute, 0, 10)

//add a action
func AddAction(idx string, I interface{}, A string) {
	action := actionPair{
		I: I,
		A: A,
	}
	actionMap[idx] = action
}

//the router for request dispatch
type Router struct {
	App *App
}

//url rewrite route
type rewriteRoute struct {
	r string
	compiledR *regexp.Regexp
	rewriteTo string
	paramMapping map[int]string
}

func NewRouter(app *App) *Router {
	return &Router{
		App: app,
	}
}


//add a regex route
//AddRoute("shop/product/(\d+)", "shop/showprod", array(1 => "prod_id"))
//will match uri "/shop/product/1001" to "shop" controller and "showprod" action, with ctx.params["prod_id"] = 1001
func AddRoute(regex string, rewriteTo string, paramMapping map[int]string){
	cr := regexp.MustCompile(regex)
	route := rewriteRoute{
		r: regex,
		compiledR: cr,
		rewriteTo: rewriteTo,
		paramMapping: paramMapping,
	}
	routes = append(routes, route)
}

//router dispatcher
func (this *Router) ServeHTTP(w http.ResponseWriter, req *http.Request){
	defer func(){
		if e := recover(); e != nil {
			log.Printf("Dispatch error: %v", e)
		}
	}()

	requestPath := req.URL.Path
	routeMatched := false
	controller := ""
	action := ""
	ctx := newWebContext(w, req)

	//serve static files, there is no security check, DISABLED on production server
	if(this.App.Env() > PRODUCTION) {
		staticFile := path.Join("public", requestPath)
		if isRegularFile(staticFile) && (req.Method == "GET" || req.Method == "HEAD") {
			http.ServeFile(w, req, staticFile)
			return
		}
	}

	//parse request params
	if(req.Method == "POST"){
		if(req.Header.Get("Content-Type") == "multipart/form-data") {
			req.ParseMultipartForm(5*1024*1024)
		}else{
			req.ParseForm()
		}
	}

	//create param map
	paramMap := make(map[string]string)
	for key, val := range req.Form {
		paramMap[key] = val[0]
	}

	//create cookie map
	cookieMap := make(map[string]string)
	for _, cookie := range req.Cookies() {
		cookieMap[cookie.Name] = cookie.Value
	}

	//handle regex route
	for _, route := range routes {
		cr := route.compiledR
		if !cr.MatchString(requestPath) {
			continue
		}
		match := cr.FindStringSubmatch(requestPath)
		if len(match[0]) != len(requestPath) {
			continue
		}
		println(match)
	}

	//handle normal controller/action
	if !routeMatched {
        if requestPath == "/" { //home page
            controller = "home"
            action = "index"
        }else{
            pathStr := requestPath[1:]
            if strings.HasSuffix(requestPath, "/"){
                pathStr = requestPath[1:len(requestPath)-1]
            }
            pathArr := strings.SplitN(pathStr, "/", 2)
            controller = pathArr[0]
            if len(pathArr) == 1 { //only have section
                action = "index"
            }else { //have section and action
                action = pathArr[1]      
            }
        }
	}

	ctx.Param = paramMap
	ctx.Cookie = cookieMap
	this.callAction(ctx, controller, action)
}

func (this *Router) callAction(ctx *WebContext, controller string, action string) {
	fmt.Fprintf(ctx.W, "Hello, you've requested: %s\n", ctx.Req.URL.Path)
}

func isRegularFile(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	mode := fileInfo.Mode()
	if mode.IsRegular() {
		return true
	}
	return false
}