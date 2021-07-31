package output

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

const (
	AppName = "[kf2-antiddos] "
)

var (
	endOfLine string      = "\n"
	devNull   *log.Logger = log.New(ioutil.Discard, "", 0)
	stdout    *log.Logger = log.New(os.Stdout, "", 0)
	stderr    *log.Logger = log.New(os.Stderr, "", 0)
	proxy     *log.Logger = log.New(os.Stdout, "", 0)
)

func ProxyMode() {
	stdout = devNull
	stderr = devNull
	proxy = log.New(os.Stdout, "", 0)
}

func SelfMode() {
	proxy = devNull
	stdout = log.New(os.Stdout, "", 0)
	stderr = log.New(os.Stderr, "", 0)
}

func AllMode() {
	stdout = log.New(os.Stdout, AppName, 0)
	stderr = log.New(os.Stderr, AppName, 0)
	proxy = log.New(os.Stdout, "", 0)
}

func StdoutWriter() io.Writer {
	return stdout.Writer()
}

func StderrWriter() io.Writer {
	return stderr.Writer()
}

func QuietMode() {
	stdout = devNull
	stderr = devNull
	proxy = devNull
}

func SetEndOfLineNative() {
	switch os := runtime.GOOS; os {
	case "windows":
		setEndOfLineWindows()
	default:
		setEndOfLineUnix()
	}
}

func EOL() string {
	return endOfLine
}

func setEndOfLineUnix() {
	endOfLine = "\n"
}

func setEndOfLineWindows() {
	endOfLine = "\r\n"
}

func Print(v ...interface{}) {
	stdout.Print(v...)
}

func Printf(format string, v ...interface{}) {
	stdout.Printf(format, v...)
}

func Println(v ...interface{}) {
	stdout.Print(fmt.Sprint(v...) + endOfLine)
}

func Error(v ...interface{}) {
	stderr.Print(v...)
}

func Errorf(format string, v ...interface{}) {
	stderr.Printf(format, v...)
}

func Errorln(v ...interface{}) {
	stderr.Print(fmt.Sprint(v...) + endOfLine)
}

func Proxy(v ...interface{}) {
	proxy.Print(v...)
}

func Proxyf(format string, v ...interface{}) {
	proxy.Printf(format, v...)
}

func Proxyln(v ...interface{}) {
	proxy.Print(fmt.Sprint(v...) + endOfLine)
}
