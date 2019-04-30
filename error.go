package maru

import(
	"fmt"
	"runtime"
	"strings"
)

//Error struct to track the error
type Error struct {
	Message string
	Cause []string
}

//Create a new error chain
func NewError(message string, cause ...string) *Error {

	chained := make([]string, 0, 3)
	detail := getLocation(2)
	chained = append(chained, fmt.Sprintf("%s {%s}", message, detail))
	chained = append(chained, cause...)
	return &Error {
		Message: message,
		Cause: chained,
	}
}

//Wrap existing error to ErrorChain
func WrapError(err error, location ...int) *Error {
	chainedErr, ok := err.(*Error)
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
	return &Error {
		Message: message,
		Cause: chained,
	}
}

//Create a new error chain from existing error
func FromError(message string, err error, location ...int) *Error {

	chained := make([]string, 0, 5)
	var detail string
	if len(location) > 0 {
		detail = getLocation(location[0])
	}else{
		detail = getLocation(2)
	}
	chained = append(chained, fmt.Sprintf("%s {%s}", message, detail))
	chainedErr, ok := err.(*Error)
	if ok {
		chained = append(chained, chainedErr.Cause...)
	}else {
		chained = append(chained, fmt.Sprintf("%s", err))
	}
	return &Error {
		Message: message,
		Cause: chained,
	}
}

//implement error interface
func (this *Error) Error() string {
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

