package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/nfv-aws/wcafe-api-controller/config"
)

type logFormat struct {
	TimestampFormat string
}

//Format ログの形式を設定
func (f *logFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteByte('[')
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("]:")
	b.WriteString(entry.Time.Format(f.TimestampFormat))

	b.WriteString(" [")
	b.WriteString(formatFilePath(entry.Caller.File))
	b.WriteString(":")
	fmt.Fprint(b, entry.Caller.Line)
	b.WriteString("] ")

	if entry.Message != "" {
		b.WriteString(" - ")
		b.WriteString(entry.Message)
	}

	if len(entry.Data) > 0 {
		b.WriteString(" || ")
	}
	for key, value := range entry.Data {
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteByte('{')
		fmt.Fprint(b, value)
		b.WriteString("}, ")
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

//init パッケージ読み込み時に実行される。
func init() {
	logrus.SetReportCaller(true) //Caller(実行ファイル(ex. main.go)を扱うため)
	formatter := logFormat{}
	formatter.TimestampFormat = "2006-01-02 15:04:05" //時刻設定

	logrus.SetFormatter(&formatter)

	//ログ出力ファイルの設定
	config.Configure()
	f, err := os.OpenFile(config.C.LOG.File_path+"/conductor.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, f))

	//ログレベルの設定
	logrus.SetLevel(logrus.InfoLevel)

}

//SetLevelDebug Traceレベルに設定
func SetLevelTrace() {
	logrus.SetLevel(logrus.TraceLevel)
}

//SetLevelDebug Debugレベルに設定
func SetLevelDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}

//SetLevelInfo Set Infoレベルに設定
func SetLevelInfo() {
	logrus.SetLevel(logrus.InfoLevel)
}

//SetLevelDebug Warnレベルに設定
func SetLevelWarn() {
	logrus.SetLevel(logrus.WarnLevel)
}

//SetLevelInfo Set Errorレベルに設定
func SetLevelError() {
	logrus.SetLevel(logrus.ErrorLevel)
}

//SetLevelDebug Fatalレベルに設定
func SetLevelFatal() {
	logrus.SetLevel(logrus.FatalLevel)
}

//SetLevelInfo Set Panicレベルに設定
func SetLevelPanic() {
	logrus.SetLevel(logrus.PanicLevel)
}
