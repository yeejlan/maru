package maru

import (
	"net/http"
	"fmt"
	"log"
	"regexp"
	"os"
	"path"
	"strings"
	"reflect"
)

//the router for request dispatch
type Router struct {
	App *App
	routes []rewriteRoute
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
		routes: make([]rewriteRoute, 0, 10),
	}
}

//add a regex route
//AddRoute("shop/product/(\d+)", "shop/showprod", array(1 => "prod_id"))
//will match uri "/shop/product/1001" to "shop" controller and "showprod" action, with ctx.params["prod_id"] = 1001
func (this *Router) AddRoute(regex string, rewriteTo string, paramMapping map[int]string){
	cr := regexp.MustCompile(regex)
	route := rewriteRoute{
		r: regex,
		compiledR: cr,
		rewriteTo: rewriteTo,
		paramMapping: paramMapping,
	}
	this.routes = append(this.routes, route)
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

	//serve static files, there is no security check, must DISABLE on production server
	if(this.App.Env() > PRODUCTION) {
		staticFile := path.Join("public", requestPath)
		if isRegularFile(staticFile) && (req.Method == "GET" || req.Method == "HEAD") {
			http.ServeFile(w, req, staticFile)
			return
		}
	}

	//parse request params
	if req.Method == "POST" && req.Header.Get("Content-Type") == "multipart/form-data" {
		req.ParseMultipartForm(5*1024*1024)
	}else{
		req.ParseForm()
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
	for _, route := range this.routes {
		cr := route.compiledR
		if !cr.MatchString(requestPath) {
			continue
		}
		match := cr.FindStringSubmatch(requestPath)
		if len(match[0]) != len(requestPath) {
			continue
		}
		//found one
		routeMatched = true
		pathArr := strings.SplitN(route.rewriteTo, "/", 2)
		controller = pathArr[0]
		action = pathArr[1]
		for k, v := range route.paramMapping {
			if(k > (len(match) -1) ){
				continue
			}
			val := match[k]
			paramMap[v] = val
		}
		break
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
	defer func(){
		if e := recover(); e != nil {
			ctx.Error = e
			internalServerError(ctx)
		}
	}()

	actionKey := strings.ToLower(fmt.Sprintf("%s/%s", controller, action))
	ap, ok := actionMap[actionKey]
	if !ok {
		pageNotFound(ctx)
		return
	}
	callMethod(ctx, &ap)
}

func callMethod(ctx *WebContext, ap *actionPair) {
	t := reflect.TypeOf(ap.I)
	instancePtr := reflect.New(t)
	instance := reflect.Indirect(instancePtr)

	//bind ctx
	ctxVal := instance.FieldByName("WebContext")
	if ctxVal.IsValid() {
		ctxVal.Set(reflect.ValueOf(ctx))
	}

	//call "Before()"
	beforeFunc := instancePtr.MethodByName("Before")
	if beforeFunc.IsValid() {
		beforeFunc.Call([]reflect.Value{})
	}

	//call "SomeAction()"
	actionFunc := instancePtr.MethodByName(ap.A + "Action")
	if !actionFunc.IsValid() {
		pageNotFound(ctx)
		return
	}
	retVal := actionFunc.Call([]reflect.Value{})

	//output
	if(len(retVal) == 0) {
		return
	}
	fmt.Fprintf(ctx.W, "%s", retVal[0])
}

func pageNotFound(ctx *WebContext) {
	actionKey := "error/page404"
	ap, ok := actionMap[actionKey]
	if !ok {
		ctx.Abort(404, "Page Not Found!")
		return
	}
	callMethod(ctx, &ap)
}

func internalServerError(ctx *WebContext) {
	actionKey := "error/page500"
	ap, ok := actionMap[actionKey]
	if !ok {
		ctx.Abort(500, "Internal Server Error!")
		return
	}
	callMethod(ctx, &ap)
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