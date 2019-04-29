package maru

import(
	"fmt"
	"runtime"
	"strings"
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
func NewError(message string, cause ...string) *ErrorChain {

	chained := make([]string, 0, 3)
	detail := getLocation(2)
	chained = append(chained, fmt.Sprintf("%s {%s}", message, detail))
	chained = append(chained, cause...)
	return &ErrorChain {
		Message: message,
		Cause: chained,
	}
}

//Wrap existing error to ErrorChain
func WrapError(err error, location ...int) *ErrorChain {
	chainedErr, ok := err.(*ErrorChain)
	if ok {
		return chainedErr
	}
	chained := make([]string, 0, 3)
	var detail string
	if len(location) > 0 {
		detail = getLocation(location[0])
	}else{
		detail = getLocation(2)
	}
	message := fmt.Sprintf("%s", err)
	chained = append(chained, fmt.Sprintf("%s {%s}", message, detail))
	return &ErrorChain {
		Message: message,
		Cause: chained,
	}
}

//Create a new error chain from existing error
func FromError(message string, err error, location ...int) *ErrorChain {

	chained := make([]string, 0, 5)
	var detail string
	if len(location) > 0 {
		detail = getLocation(location[0])
	}else{
		detail = getLocation(2)
	}
	chained = append(chained, fmt.Sprintf("%s {%s}", message, detail))
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
	err := strings.Join(this.Cause, ", ")
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
	return fmt.Sprintf("%s:%s:%d", function, file, line)
}

