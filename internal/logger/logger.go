package logger

import (
	consts "crontab/internal/const"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sync"
	"time"
)

var sysLog *log.Logger
var runLog *wyLogger
var dirPath *string

func SysWriteLn(v ...interface{}) {
	sysLog.Println(v...)
}

func SysPrintf(fmt string, vals ...interface{}) {
	sysLog.Printf(fmt, vals...)
}

func Printf(fmt string, vals ...interface{}) {
	runLog.lg.Printf(fmt, vals...)
}

func Configure(path *string) {
	dirPath = path
	slogs := *dirPath
	if slogs[len(slogs)-1] != '/' {
		*dirPath = *dirPath + "/"
	}

	if _, err := os.Stat(*dirPath); os.IsNotExist(err) {
		err = os.Mkdir(*dirPath, fs.ModePerm)
		if err != nil {
			fmt.Printf("Mkdir err: %s", err)
			os.Exit(1)
		}
	}

	sysLogFile, err := os.OpenFile(*dirPath+consts.SVR_LOG, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("%s \nStart failed!\n", err)
		os.Exit(1)
	}
	sysLog = log.New(sysLogFile, "", log.LstdFlags)
	runLog = newWyLogger(*dirPath, consts.RUN_LOG_POSTFIX)
}

func FilePath() string {
	return *dirPath
}

type wyLogger struct {
	dir      string
	filename string
	_date    *time.Time
	mu       *sync.RWMutex
	logfile  *os.File
	lg       *log.Logger
}

func newWyLogger(dir string, filename string) *wyLogger {

	logger := &wyLogger{dir: dir, filename: filename}
	logger._date = new(time.Time)
	logger.mu = new(sync.RWMutex)
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.rename()

	go fileMonitor(logger)
	return logger
}

func (l *wyLogger) isMustRename() bool {
	t, _ := time.Parse(consts.DATEFORMAT, time.Now().Format(consts.DATEFORMAT))
	if t.After(*l._date) {
		return true
	}
	return false
}

func (l *wyLogger) rename() {

	if l.isMustRename() {
		if l.logfile != nil {
			l.logfile.Close()
		}
		tf := time.Now().Format(consts.DATEFORMAT)
		t, _ := time.Parse(consts.DATEFORMAT, tf)
		l._date = &t
		fn := l.dir + tf + "_" + l.filename
		l.logfile, _ = os.OpenFile(fn, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		l.lg = log.New(l.logfile, "", log.LstdFlags)
	}
}

func fileMonitor(l *wyLogger) {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			go fileCheck(l)
		}
	}
}

func fileCheck(l *wyLogger) {
	if l != nil && l.isMustRename() {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.rename()
	}
}
