package std

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/beatlabs/patron/log"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	var b bytes.Buffer
	logger := NewLogger(&b, log.InfoLevel, map[string]interface{}{"name": "john doe", "age": 18})
	assert.NotNil(t, logger.debug)
	assert.NotNil(t, logger.info)
	assert.NotNil(t, logger.warn)
	assert.NotNil(t, logger.error)
	assert.NotNil(t, logger.fatal)
	assert.NotNil(t, logger.panic)
	assert.Equal(t, log.InfoLevel, logger.Level())
	assert.Equal(t, logger.fields, map[string]interface{}{"name": "john doe", "age": 18})
	assert.Contains(t, logger.fieldsLine, "age=18")
	assert.Contains(t, logger.fieldsLine, "name=john doe")
}

func TestNewSub(t *testing.T) {
	var b bytes.Buffer
	logger := NewLogger(&b, log.InfoLevel, map[string]interface{}{"name": "john doe"})
	assert.NotNil(t, logger)
	subLogger := logger.Sub(map[string]interface{}{"age": 18}).(*Logger)
	assert.NotNil(t, subLogger.debug)
	assert.NotNil(t, subLogger.info)
	assert.NotNil(t, subLogger.warn)
	assert.NotNil(t, subLogger.error)
	assert.NotNil(t, subLogger.fatal)
	assert.NotNil(t, subLogger.panic)
	assert.Equal(t, log.InfoLevel, subLogger.Level())
	assert.Equal(t, subLogger.fields, map[string]interface{}{"name": "john doe", "age": 18})
	assert.Contains(t, subLogger.fieldsLine, "age=18")
	assert.Contains(t, subLogger.fieldsLine, "name=john doe")
}

func TestLogger(t *testing.T) {
	// BEWARE: Since we are testing the log output change in line number of statements affect the test outcome
	var b bytes.Buffer
	logger := NewLogger(&b, log.InfoLevel, map[string]interface{}{"name": "john doe", "age": 18})

	type args struct {
		lvl        log.Level
		msg        string
		args       []interface{}
		lineNumber int
	}
	tests := map[string]struct {
		args args
	}{
		"debug":  {args: args{lvl: log.DebugLevel, args: []interface{}{"hello world"}, lineNumber: 58}},
		"debugf": {args: args{lvl: log.DebugLevel, msg: "Hi, %s", args: []interface{}{"John"}, lineNumber: 60}},
		"info":   {args: args{lvl: log.InfoLevel, args: []interface{}{"hello world"}, lineNumber: 64}},
		"infof":  {args: args{lvl: log.InfoLevel, msg: "Hi, %s", args: []interface{}{"John"}, lineNumber: 66}},
		"warn":   {args: args{lvl: log.WarnLevel, args: []interface{}{"hello world"}, lineNumber: 70}},
		"warnf":  {args: args{lvl: log.WarnLevel, msg: "Hi, %s", args: []interface{}{"John"}, lineNumber: 72}},
		"error":  {args: args{lvl: log.ErrorLevel, args: []interface{}{"hello world"}, lineNumber: 76}},
		"errorf": {args: args{lvl: log.ErrorLevel, msg: "Hi, %s", args: []interface{}{"John"}, lineNumber: 78}},
		"panic":  {args: args{lvl: log.PanicLevel, args: []interface{}{"hello world"}, lineNumber: 83}},
		"panicf": {args: args{lvl: log.PanicLevel, msg: "Hi, %s", args: []interface{}{"John"}, lineNumber: 87}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			defer b.Reset()

			switch tt.args.lvl {
			case log.DebugLevel:
				if tt.args.msg == "" {
					logger.Debug(tt.args.args...)
				} else {
					logger.Debugf(tt.args.msg, tt.args.args...)
				}
			case log.InfoLevel:
				if tt.args.msg == "" {
					logger.Info(tt.args.args...)
				} else {
					logger.Infof(tt.args.msg, tt.args.args...)
				}
			case log.WarnLevel:
				if tt.args.msg == "" {
					logger.Warn(tt.args.args...)
				} else {
					logger.Warnf(tt.args.msg, tt.args.args...)
				}
			case log.ErrorLevel:
				if tt.args.msg == "" {
					logger.Error(tt.args.args...)
				} else {
					logger.Errorf(tt.args.msg, tt.args.args...)
				}
			case log.PanicLevel:
				if tt.args.msg == "" {
					assert.Panics(t, func() {
						logger.Panic(tt.args.args...)
					})
				} else {
					assert.Panics(t, func() {
						logger.Panicf(tt.args.msg, tt.args.args...)
					})
				}
			}

			if tt.args.msg == "" {
				assert.Contains(t, b.String(), getMessage(tt.args.lineNumber, tt.args.lvl))
			} else {
				assert.Contains(t, b.String(), getMessagef(tt.args.lineNumber, tt.args.lvl))
			}
		})
	}
}

func getMessage(lineNumber int, lvl log.Level) string {
	return fmt.Sprintf("logger_test.go:%d: %s age=18 name=john doe hello world", lineNumber, levelMap[lvl])
}

func getMessagef(lineNumber int, lvl log.Level) string {
	return fmt.Sprintf("logger_test.go:%d: %s age=18 name=john doe Hi, John", lineNumber, levelMap[lvl])
}

var buf bytes.Buffer

func BenchmarkLogger(b *testing.B) {
	var tmpBuf bytes.Buffer
	logger := NewLogger(&tmpBuf, log.DebugLevel, map[string]interface{}{"name": "john doe", "age": 18})
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		logger.Debugf("Hello %s!", "John")
	}
	buf = tmpBuf
}
