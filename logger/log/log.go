package log

import (
	"os"

	"github.com/chmking/horde/logger"
)

// Logger is the global logger.
var Logger = logger.NewStdLogger(os.Stderr)

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event
func Info() logger.Event {
	return Logger.Info()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event
func Error() logger.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level.
//
// You must call Msg on the returned event in order to send the event
func Fatal() logger.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level.
//
// You must call Msg on the returned event in order to send the event
func Panic() logger.Event {
	return Logger.Panic()
}
