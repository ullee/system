package logger

import (
    "github.com/op/go-logging"
    "go/build"
    "os"
    "time"
)

var Log logging.Logger

func init() {

    date := time.Now().Format("2006-01-02")

    var logPath = build.Default.GOPATH + "/logs/log-" + date + ".log"

    file, err := os.OpenFile(logPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil {
        panic(err)
    }

    //logging.MustGetLogger("system")
    var format = logging.MustStringFormatter(
        `[%{time:2006-01-02 15:04:05.000}][%{program}][%{shortfile}][%{shortfunc}][%{level:.5s}] [%{pid}] %{message}`,
    )
    //var format = logging.MustStringFormatter(
    //    `%{color}[%{time:2006-01-02 15:04:05.000}][%{shortfunc}][%{shortfile}][%{level:.4s}] [%{pid}]%{color:reset} %{message}`,
    //)

    // For demo purposes, create two backend for os.Stderr.
    //backend1 := logging.NewLogBackend(file, "", 0)
    backend := logging.NewLogBackend(file, "", 0)

    // For messages written to backend2 we want to add some additional
    // information to the output, including the used log level and the name of
    // the function.
    backendFormatter := logging.NewBackendFormatter(backend, format)

    // Only errors and more severe messages should be sent to backend1
    //backend1Leveled := logging.AddModuleLevel(backend1)
    //backend1Leveled.SetLevel(logging.ERROR, "")

    // Set the backends to be used.
    //logging.SetBackend(backend1Leveled, backendFormatter)
    logging.SetBackend(backendFormatter)
}

