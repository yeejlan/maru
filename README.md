# maru
A golang web framework

Here is the demo project [release-manager](https://github.com/yeejlan/go_release_manager)

```golang
package main

import(
	"github.com/yeejlan/maru"
	"release_manager/controller"
	"release_manager/dal"
)

func main() {

	app := maru.NewApp("development", "release_manager")
	app.Init()

	//load db and redis
	loader := maru.NewResourceLoader(app)
	loader.Autoload()

	dal.InitSharedVars()

	//load actions
	controller.LoadActions()

	//add rewrite route
	router := maru.NewRouter(app)
	router.AddRoute("/hello/(.*)", "home/hi", map[int]string{1 : "username"});

	maru.StartHttpServer(router, "0.0.0.0", 8080)
}
```	
