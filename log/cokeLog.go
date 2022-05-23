package log

import (
	"github.com/chuccp/utils/queue"
	log2 "log"
	"os"
	"time"
)

type Level uint32

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func (l Level) Level() string {
	if l == PanicLevel {
		return "panic"
	}
	if l == FatalLevel {
		return "fatal"
	}
	if l == ErrorLevel {
		return "error"
	}
	if l == WarnLevel {
		return "warn"
	}
	if l == InfoLevel {
		return "info"
	}
	if l == DebugLevel {
		return "debug"
	}
	if l == TraceLevel {
		return "trace"
	}
	return ""
}

type Logger struct {
	config    *Config
	queue     *queue.Queue
	fileQueue *queue.Queue
}

func New() *Logger {
	log := &Logger{config: defaultConfig, queue: queue.NewQueue(), fileQueue: queue.NewQueue()}
	go log.printLog()
	go func() {
		err := log.printFileLog()
		if err != nil {
			log2.Panic(err)
		}
	}()
	return log
}
func (logger *Logger) printLog() {
	out := os.Stdout
	for {
		v, _ := logger.queue.Poll()
		p := v.(*Entry)
		p.WriteTo(out)
		freeEntry(p)
	}
}
func (logger *Logger) printLevelMapFileLog(fileCut *cut) (err error) {
	var writeFileMap = make(map[Level]*WriteFile)
	for {
		v, _ := logger.fileQueue.Poll()
		p := v.(*Entry)
		if writeFileMap[p.Level] == nil {
			writeFileMap[p.Level] = NewWriteFile(fileCut)
		}
		outFile := writeFileMap[p.Level]
		err := outFile.fileTo(p.now, p.Level)
		if err != nil {
			return err
		}
		outFile.WriteLog(p)
		if err != nil {
			return err
		}
		freeEntry(p)
	}
}
func (logger *Logger) printLevelSingleFileLog(fileCut *cut) (err error) {

	outFile := NewWriteFile(fileCut)

	for {
		v, _ := logger.fileQueue.Poll()
		p := v.(*Entry)

		err := outFile.fileTo(p.now, p.Level)
		if err != nil {
			return err
		}
		outFile.WriteLog(p)
		if err != nil {
			return err
		}
		freeEntry(p)
	}
}

func (logger *Logger) printFileLog() (err error) {
	cut, err := parse(logger.config.filePattern)
	if err != nil {
		return err
	}
	if cut.hasLevel {
		return logger.printLevelMapFileLog(cut)
	} else {
		return logger.printLevelSingleFileLog(cut)
	}
}
func (logger *Logger) Info(format string, args ...interface{}) {
	logger.log(InfoLevel, format, args...)
}
func (logger *Logger) Debug(format string, args ...interface{}) {
	logger.log(DebugLevel, format, args...)
}
func (logger *Logger) Trace(format string, args ...interface{}) {
	logger.log(TraceLevel, format, args...)
}
func (logger *Logger) Fatal(format string, args ...interface{}) {
	logger.log(FatalLevel, format, args...)
}
func (logger *Logger) Panic(format string, args ...interface{}) {
	logger.log(PanicLevel, format, args...)
}

func (logger *Logger) Error(format string, args ...interface{}) {
	logger.log(ErrorLevel, format, args...)
}
func (logger *Logger) log(level Level, format string, args ...interface{}) {
	if logger.config.level >= level || logger.config.fileLevel >= level {
		now := time.Now()
		if logger.config.level >= level {
			p := newEntry(logger.config, level, &now)
			p.Log(format, args...)
			logger.queue.Offer(p)
		}
		if logger.config.fileLevel >= level {
			p := newEntry(logger.config, level, &now)
			p.Log(format, args...)
			logger.fileQueue.Offer(p)
		}
	}
}

var defaultLogger = New()

func InfoF(format string, value ...interface{}) {
	defaultLogger.Info(format, value...)
}
func DebugF(format string, value ...interface{}) {
	defaultLogger.Debug(format, value...)
}
func FatalF(format string, value ...interface{}) {
	defaultLogger.Fatal(format, value...)
}
func ErrorF(format string, value ...interface{}) {
	defaultLogger.Error(format, value...)
}
func TraceF(format string, value ...interface{}) {
	defaultLogger.Trace(format, value...)
}
func PanicF(format string, value ...interface{}) {
	defaultLogger.Panic(format, value...)
}

func Info(value ...interface{}) {
	defaultLogger.Info("", value...)
}
func Debug(value ...interface{}) {
	defaultLogger.Debug("", value...)
}
func Fatal(value ...interface{}) {
	defaultLogger.Fatal("", value...)
}
func Error(value ...interface{}) {
	defaultLogger.Error("", value...)
}
func Trace(value ...interface{}) {
	defaultLogger.Trace("", value...)
}
func Panic(value ...interface{}) {
	defaultLogger.Panic("", value...)
}
