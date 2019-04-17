package maru

import (
	"os"
	"time"
	"fmt"
	"path"
)

//rotate log by date
type DailyLogRotate struct {
	logname string
	logdir string
	filename string
	fd *os.File
}

//create new log rotate
func NewDailyLogRotate(logDir string, logName string) *DailyLogRotate {
	fname := getLogFileName(logName)
	fullpath := path.Join(logDir, fname)
	fd, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return &DailyLogRotate{
		logname: logName,
		logdir: logDir,
		filename: fname,
		fd: fd,
	}
}

//implement io.Writer interface for log module
func (this *DailyLogRotate) Write(p []byte) (n int, err error) {
	fname := getLogFileName(this.logname)
	if(fname != this.filename || nil == this.fd){
		this.Close()
		this.filename = fname
		fullpath := path.Join(this.logdir, fname)
		//ignore openfile error
		fd, _ := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		this.fd = fd
	}
	if(nil == this.fd){ //ignore any error
		return 0, nil
	}
	n, err = this.fd.Write(p)
	if err != nil {
		this.Close()
	}
	return
}

//close log
func (this *DailyLogRotate) Close() {
	if(this.fd != nil) {
		this.fd.Sync()
		this.fd.Close()
		this.fd = nil
	}
}

func getLogFileName(logName string) string {
	dt := time.Now()
	date := dt.Format("2006_01_02")
	name := fmt.Sprintf("%s_%s.log", logName, date)
	return name
}
