package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	LAYOUT = "2006-01-02"
)

func getTodayLogName(prefix string) string {
	return prefix + time.Now().Format(LAYOUT) + ".log"
}

type Writer struct {
	name   string
	prefix string
	fp     *os.File
	mutex  sync.Mutex
}

func (w *Writer) Write(data []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	today := getTodayLogName(w.prefix)
	if !strings.EqualFold(today, w.name) {
		var fp *os.File
		fp, err = os.OpenFile(today, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return
		}
		w.fp.Close()
		w.name = today
		w.fp = fp
	}

	n, err = w.fp.Write(data)
	return
}

func (w *Writer) Close() (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err = w.fp.Close()
	return
}

func NewWriter(name string) (io.WriteCloser, error) {
	prefix := filepath.Dir(name) + "/"
	today := getTodayLogName(prefix)
	fp, err := os.OpenFile(today, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Writer{name: today, prefix: prefix, fp: fp}, nil
}

type Logger struct {
	logger *log.Logger
	writer io.WriteCloser
}

func NewLogger(dir, prefix string, flag int) (*Logger, error) {
	writer, err := NewWriter(dir)
	if err != nil {
		return nil, err
	}
	logger := log.New(writer, prefix, flag)
	return &Logger{logger: logger, writer: writer}, nil
}

func (this *Logger) Print(v ...interface{}) {
	this.logger.Print(v...)
}

func (this *Logger) Printf(str string, v ...interface{}) {
	this.logger.Printf(str, v...)
}

func (this *Logger) Error(v ...interface{}) {
	tmp := make([]interface{}, len(v)+1)
	tmp[0] = "[ERROR]"
	copy(tmp[1:], v)
	this.logger.Print(tmp...)
}

func (this *Logger) Errorf(str string, v ...interface{}) {
	str = "[ERROR]" + str
	this.logger.Printf(str, v...)
}

func (this *Logger) Fatal(v ...interface{}) {
	tmp := make([]interface{}, len(v)+1)
	tmp[0] = "[ERROR]"
	copy(tmp[1:], v)
	this.logger.Print(tmp...)
	os.Exit(1)
}

func (this *Logger) Fatalf(str string, v ...interface{}) {
	str = "[ERROR]" + str
	this.logger.Printf(str, v...)
	os.Exit(1)
}

func (this *Logger) Close() {
	if this.writer != nil {
		this.writer.Close()
	}
}

var DefaultLogger *Logger = &Logger{log.New(os.Stdout, "", log.LstdFlags), os.Stdout}

func Print(v ...interface{}) {
	DefaultLogger.Print(v...)
}

func Printf(str string, v ...interface{}) {
	DefaultLogger.Printf(str, v...)
}

func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

func Errorf(str string, v ...interface{}) {
	DefaultLogger.Errorf(str, v...)
}

func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}

func Fatalf(str string, v ...interface{}) {
	DefaultLogger.Fatalf(str, v...)
}
