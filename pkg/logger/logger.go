package logger

import "log"

const (
	_ = iota
	FATAL
	PANIC
	PRINT
)

var IsDebugMode bool = false

func DebugLog(mode int, format string, v ...interface{}) {
	if !IsDebugMode {
		return
	}
	switch mode {
	case FATAL:
		log.Fatalf(format, v...)
	case PANIC:
		log.Panicf(format, v...)
	case PRINT:
		log.Printf(format, v...)
	default:
		log.Panic("undefined log mode")
	}
}
