package std

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	patronLog "github.com/beatlabs/patron/log"
)

var levelMap = map[patronLog.Level]string{
	patronLog.DebugLevel: "DBG",
	patronLog.InfoLevel:  "INF",
	patronLog.WarnLevel:  "WRN",
	patronLog.ErrorLevel: "ERR",
	patronLog.FatalLevel: "FTL",
	patronLog.PanicLevel: "PNC",
}

type Logger struct {
	level      patronLog.Level
	fields     map[string]interface{}
	fieldsLine string
	out        io.Writer
	debug      *log.Logger
	info       *log.Logger
	warn       *log.Logger
	error      *log.Logger
	panic      *log.Logger
	fatal      *log.Logger
}

func NewLogger(out io.Writer, lvl patronLog.Level, fields map[string]interface{}) *Logger {

	fieldsLine := createFieldsLine(fields)

	return &Logger{
		out:        out,
		debug:      createLogger(out, patronLog.DebugLevel, fieldsLine),
		info:       createLogger(out, patronLog.InfoLevel, fieldsLine),
		warn:       createLogger(out, patronLog.WarnLevel, fieldsLine),
		error:      createLogger(out, patronLog.ErrorLevel, fieldsLine),
		panic:      createLogger(out, patronLog.PanicLevel, fieldsLine),
		fatal:      createLogger(out, patronLog.FatalLevel, fieldsLine),
		level:      lvl,
		fields:     fields,
		fieldsLine: fieldsLine,
	}
}

func createFieldsLine(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	// always return the fields in the same order
	keys := make([]string, 0, len(fields))
	for key, _ := range fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	sb := strings.Builder{}

	for _, key := range keys {
		sb.WriteString(key)
		sb.WriteRune('=')
		sb.WriteString(fmt.Sprintf("%v", fields[key]))
		sb.WriteRune(' ')
	}

	return sb.String()
}

func createLogger(out io.Writer, lvl patronLog.Level, fieldLine string) *log.Logger {
	logger := log.New(out, levelMap[lvl]+" "+fieldLine, log.LstdFlags|log.Lmicroseconds|log.LUTC|log.Lmsgprefix|log.Lshortfile)
	return logger
}

func (l *Logger) Sub(fields map[string]interface{}) patronLog.Logger {

	for key, value := range l.fields {
		fields[key] = value
	}

	return NewLogger(l.out, l.level, fields)
}

func (l *Logger) Fatal(args ...interface{}) {
	output(l.fatal, args...)
	os.Exit(1)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	outputf(l.fatal, msg, args...)
	os.Exit(1)
}

func (l *Logger) Panic(args ...interface{}) {
	panic(output(l.panic, args...))
}

func (l *Logger) Panicf(msg string, args ...interface{}) {
	panic(outputf(l.panic, msg, args...))
}

func (l *Logger) Error(args ...interface{}) {
	output(l.error, args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	outputf(l.error, msg, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	output(l.warn, args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	outputf(l.warn, msg, args...)
}

func (l *Logger) Info(args ...interface{}) {
	output(l.info, args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	outputf(l.info, msg, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	output(l.debug, args...)
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	outputf(l.debug, msg, args...)
}

func (l *Logger) Level() patronLog.Level {
	return l.level
}

func output(logger *log.Logger, args ...interface{}) string {
	msg := fmt.Sprint(args...)
	_ = logger.Output(3, msg)
	return msg
}

func outputf(logger *log.Logger, msg string, args ...interface{}) string {
	fmtMsg := fmt.Sprintf(msg, args...)
	_ = logger.Output(3, fmtMsg)
	return fmtMsg
}
