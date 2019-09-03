package testhelper

import (
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

//LogCapturer logger interface
type LogCapturer interface {
	//Release release log
	Release()
	//SetJSONFormatLog set formmater
	SetJSONFormatLog(destinationFile string) (err error)
}

type logCapturer struct {
	*testing.T
	origOut       io.Writer
	jsonFormatter *logrus.TextFormatter
}

func (tl logCapturer) Format(entry *logrus.Entry) (rslt []byte, err error) {
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.PanicLevel {
		tl.Error(entry.Message)
	} else if entry.Level == logrus.FatalLevel {
		tl.Fatal(entry.Message)
	} else {
		tl.Logf(entry.Message)
	}
	return tl.jsonFormatter.Format(entry)
}
func (tl logCapturer) Write(p []byte) (n int, err error) {
	tl.Logf((string)(p))
	return len(p), nil
}

func (tl logCapturer) Release() {
	logrus.SetOutput(tl.origOut)
}

//SetJSONFormatLog set logger to json. and log to fle
func (tl logCapturer) SetJSONFormatLog(destinationFile string) (err error) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	f, err1 := os.OpenFile(destinationFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	err = err1
	logrus.SetOutput(f)
	return

}

// CaptureLog redirects logrus output to testing.Log
func CaptureLog(t *testing.T) LogCapturer {
	fmt := logrus.TextFormatter{}
	lc := logCapturer{T: t, origOut: logrus.StandardLogger().Out, jsonFormatter: &fmt}
	logrus.SetOutput(lc)

	if !testing.Verbose() {
		logrus.SetOutput(lc)
	}
	return &lc
}
