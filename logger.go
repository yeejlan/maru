package maru

import (
	"os"
	"time"
	"fmt"
	"path"
	"sync"
	"container/list"
)


const(
	//max opened log files
	maxFileOpened = 16
	logFilePerm = 0666
	logDirPerm = 0777
)


//loggerHolder singleton
var lholder = newLoggerHolder()

type loggerHolder struct {
	basedir string
	cache *writerCache
}

//create loggerHolder
func newLoggerHolder() *loggerHolder {
	return &loggerHolder{
		basedir: "logs",
		cache: newWriterCache(maxFileOpened),
	}
}

//init log setting
func InitLog(basedir string) {
	lholder.basedir = basedir
}

//log message
func Log(prefix string, message string) {
	msg := fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), message)
	logger := getLogger(prefix)
	logger.log([]byte(msg))
}

//log message with format
func Logf(prefix string, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	msg := fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), s)
	logger := getLogger(prefix)
	logger.log([]byte(msg))
}

//logger rotate by date
type Logger struct {
	prefix string
	basedir string
	logfile string
	fd *os.File
}

//create new logger without cache
func NewLogger(basedir string, prefix string) *Logger {
	return newLogger(basedir, prefix, true)
}

func newLogger(basedir string, prefix string, isPanic bool) *Logger {
	p := getLogPath(prefix)
	fd := openLogFile(basedir, p, isPanic)
	return &Logger{
		prefix: prefix,
		basedir: basedir,
		logfile: p.logfile,
		fd: fd,
	}
}

//get logger from cache
func getLogger(prefix string) *Logger {
	return lholder.cache.getOrNew(prefix)
}

//implement io.Writer interface for log module
func (this *Logger) Write(payload []byte) (n int, err error) {
	p := getLogPath(this.prefix)
	if(p.logfile != this.logfile || this.fd == nil){
		this.Close()
		fd := openLogFile(this.basedir, p, false)
		this.fd = fd
		this.logfile = p.logfile
	}
	return this.log(payload)
}

//write log
func (this *Logger) log(p []byte) (n int, err error) {
	if(this.fd == nil){ //ignore error
		return 0, nil
	}
	n, err = this.fd.Write(p)
	if err != nil {
		this.Close()
	}
	return
}

//sync log
func (this *Logger) Sync() {
	if(this.fd != nil) {
		this.fd.Sync()
	}
}

//close log
func (this *Logger) Close() {
	if(this.fd != nil) {
		this.fd.Close()
		this.fd = nil
	}
}

func openLogFile(baseDir string, p *logPath, isPanic bool) *os.File {
	fullpath := path.Join(baseDir, p.logfile)
	fd, err := os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, logFilePerm)
	if err != nil {
		//try again, mkdir first
		os.Mkdir(path.Join(baseDir, p.y), logDirPerm)
		os.Mkdir(path.Join(baseDir, p.y, p.m), logDirPerm)
		fd, err = os.OpenFile(fullpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, logFilePerm)
		if err != nil && isPanic == true {
			panic(err)
		}
	}
	return fd
}

type logPath struct {
	//year
	y string
	//month
	m string
	logfile string
}

func getLogPath(prefix string) *logPath {
	dt := time.Now()
	year := dt.Format("2006")
	month := dt.Format("01")
	day := dt.Format("02")
	//2006/01/prefix_02.log
	logpath := fmt.Sprintf("%s/%s/%s_%s.log", year, month, prefix, day)
	return &logPath{y:year, m: month, logfile: logpath}
}

type writerCache struct {
	mu sync.RWMutex
	cache map[string]*Logger
	lst *list.List
	maxItem int
}

func newWriterCache(maxItem int) *writerCache {
	return &writerCache{
		cache: make(map[string]*Logger),
		lst: list.New(),
		maxItem: maxItem,
	}
}

func (this *writerCache) put(logfile string, writer *Logger) {
	this.cache[logfile] = writer
	this.lst.PushBack(logfile)

	for(this.lst.Len() > this.maxItem) {
		e := this.lst.Front()
		logfile := e.Value.(string)
		logger := this.cache[logfile]
		if logger != nil {
			logger.Close()
		}
		delete(this.cache, logfile)
		this.lst.Remove(e)
	}
}

func (this *writerCache) get(logfile string) (l *Logger, ok bool) {
	l, ok = this.cache[logfile]
	return
}

func (this *writerCache) getOrNew(prefix string) *Logger {
	p := getLogPath(prefix)

	this.mu.RLock()
	//from cache
	logger, ok := this.get(p.logfile)
	if ok {
		return logger
	}
	this.mu.RUnlock()

	this.mu.Lock()
	defer this.mu.Unlock()

	logger, ok = this.get(p.logfile)
	if ok {
		return logger
	}
	//not found in cache, create new one
	logger = newLogger(lholder.basedir, prefix, false)
	if logger.fd != nil {
		this.put(logger.logfile, logger)
	}
	return logger
}