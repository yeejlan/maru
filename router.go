package maru

import (
	"net/http"
	"fmt"
	"log"
	"regexp"
	"os"
	"path"
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
	_ = controller
	_ = action
	//serve static files, there is no security check, DISABLED on production server
	if(this.App.Env() > PRODUCTION) {
		staticFile := path.Join("public", requestPath)
		if isRegularFile(staticFile) && (req.Method == "GET" || req.Method == "HEAD") {
			http.ServeFile(w, req, staticFile)
			return
		}
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

	}

	fmt.Fprintf(w, "Hello, you've requested: %s\n", req.URL.Path)
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