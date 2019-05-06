package maru

import( 
	"fmt"
	"path/filepath"
	"os"
	"strings"
	"log"
	"bufio"
	"regexp"
	"bytes"
	"sort"
)

type actionInfo struct {
	Controller string
	Action string
}

var tpl = `
//auto generated file, please do not modify.

import "github.com/yeejlan/maru"
import "reflect"

func LoadActions() {
`

type genController struct {
	packageName string
	controllerDir string
	outFile string
	actionMap map[string]actionInfo
	controllerSuffix string
}

//generate controller/action mapping
func NewGenController(packageName string, controllerDir string) *genController{
	return &genController{
		packageName: packageName,
		controllerDir: controllerDir,
		outFile: controllerDir + "/controller.go",
		actionMap: make(map[string]actionInfo),
		controllerSuffix: "Controller.go",
	}
}

//generate file
func (this *genController) Generate() {
	this.getControllerList()

	var sortedKeys []string
	for idx := range this.actionMap {
		sortedKeys = append(sortedKeys, idx)
	}
	sort.Strings(sortedKeys)

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("package %s", this.packageName))
	buffer.WriteString(tpl)	
	for _, key := range sortedKeys {
		v := this.actionMap[key]
		//maru.AddAction("home/index", reflect.TypeOf(HomeController{}), "index") 
		action := fmt.Sprintf("\tmaru.AddAction(\"%s\", reflect.TypeOf(%sController{}), \"%s\")\n",
			key, v.Controller, v.Action)

		buffer.WriteString(action)
	}
	buffer.WriteString("\n}")

	//write file
	f, err := os.OpenFile(this.outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s generated", this.outFile)	
}

func (this *genController) getControllerList() {
	cpath := this.controllerDir
	err := filepath.Walk(cpath, func(cpath string, f os.FileInfo, err error) error {
			if (f == nil) {return err}
			if f.IsDir() {return nil}
			if(strings.HasSuffix(cpath, this.controllerSuffix)){
				this.getActionList(cpath)
			}
			return nil
		})

	if err != nil {
		log.Fatalf("filepath.Walk() error %v\n", err)
	}
}

func (this *genController) getActionList(cpath string) {
	controller := strings.TrimSuffix(filepath.Base(cpath), this.controllerSuffix)
	//func (this *TestController) indexAction() string {
	validAction := regexp.MustCompile(`func[[:space:]]+\(.*` + controller+ `.*\)[[:space:]]+([A-Z][a-zA-Z0-9]*)Action\([[:space:]]*\).*\{`)
	file, err := os.Open(cpath)
		if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !validAction.MatchString(line) {
			continue
		}
		match := validAction.FindStringSubmatch(line)
		action := match[1]
		actionKey := strings.ToLower(fmt.Sprintf("%s/%s", controller, action))
		this.actionMap[actionKey] = actionInfo{Controller: controller, Action: action,}
	}
	
}
