package logging

import (
	stdlog "log"
	"os"

	golog "github.com/op/go-logging"
)

// go-logging 没有 verbose 级别,将日志级别降级一档
// Debug -> Verbose
// Info -> Debug
// Notice -> Info
// Warn, Error, and Critical 保持一致.
var defaultLevel golog.Level = golog.DEBUG

// 回头试试 https://github.com/Sirupsen/logrus docker 在用
//参考https://sourcegraph.com/github.com/sputnik-maps/gopnik/-/blob/src/loghelper/loghelper.go#L0
// setlevel 时,变量指定 logger name
func Configure() {
	golog.SetFormatter(golog.MustStringFormatter("[0x%{id:x}] [%{level}] [%{module}] %{message}"))
	stdoutLogBackend := golog.NewLogBackend(os.Stdout, "", stdlog.LstdFlags|stdlog.Lshortfile)
	stdoutLogBackend.Color = true
	golog.SetLevel(defaultLevel, "")
	//logging.SetLevel(logging.INFO, "threadpool") 可区分 MustGetLogger 的模块以区分文件

	// NOTE these file permissions are restricted by umask, so they probably won't work right.
	err := os.MkdirAll("./log", 0775)
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile("./log/cola.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}

	fileLogBackend := golog.NewLogBackend(logFile, "", stdlog.LstdFlags|stdlog.Lshortfile)
	fileLogBackend.Color = false

	golog.SetBackend(stdoutLogBackend, fileLogBackend)

}
