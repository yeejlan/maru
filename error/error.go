//Package error provide a chained error to track the root cause of an error
package error

import(
	"fmt"
	"runtime"
	"strings"
)

var(
	println = fmt.Println
	//Setting with go build -ldflags "-X github.com/yeejlan/maru/error.BuildDir=xxx"
	BuildDir string
)

//ErrorChain struct to track the root error
type ErrorChain struct {
	Message string
	Cause []string
}

func init() {
	BuildDir = strings.ReplaceAll(BuildDir, "\\", "/")
}

//Create a new error chain
func New(message string, cause ...string) *ErrorChain {

	chained := make([]string, 0)
	detail := getLocation(2)
	chained = append(chained, fmt.Sprintf("%s:%s", message, detail))
	chained = append(chained, cause...)
	return &ErrorChain {
		Message: message,
		Cause: chained,
	}
}

//Create a new error chain from existing error
func From(message string, err error) *ErrorChain {

	chained := make([]string, 0)
	detail := getLocation(2)
	chained = append(chained, fmt.Sprintf("%s:%s", message, detail))
	chainedErr, ok := err.(*ErrorChain)
	if ok {
		chained = append(chained, chainedErr.Cause...)
	}else {
		chained = append(chained, fmt.Sprintf("%s", err))
	}
	return &ErrorChain {
		Message: message,
		Cause: chained,
	}
}

//implement error interface
func (this *ErrorChain) Error() string {
	err := strings.Join(this.Cause, ",")
	return err
}

//get error location
func getLocationDetail(callDepth int) (file string, line int, function string) {
	pc, file, line, _ := runtime.Caller(callDepth + 1)
	f := runtime.FuncForPC(pc)
	function = f.Name()

	file = strings.TrimPrefix(file, BuildDir)

	return
}

//get error location as string
func getLocation(callDepth int) string {
	file, line, function := getLocationDetail(callDepth)
	return fmt.Sprintf("%s(%s:%d)", function, file, line)
}

