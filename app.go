package maru

type appEnv int

const (
	PRODUCTION appEnv = iota + 1
	STAGING
	TESTING
	DEVELOPMENT
)

type App struct {
	isInit bool
	env appEnv
	envString string
	configFile string
	config string
	appName string
}

func NewApp(env string, appName string) *App {
	return &App{
		isInit: false,
		env: PRODUCTION,
		envString: "production",
	}
}

func (this *App) Init() {
	this.isInit = true
	println("app is init...")
}