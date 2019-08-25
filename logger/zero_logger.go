package logger

import (
	"io"

	"github.com/rs/zerolog"
)

func NewZeroConsoleWriter(writer io.Writer) io.Writer {
	return zerolog.ConsoleWriter{Out: writer}
}

var _ = Logger(&ZeroLogger{})

func NewZeroLogger(writer io.Writer) *ZeroLogger {
	return &ZeroLogger{
		logger: zerolog.New(writer).With().Timestamp().Logger(),
	}
}

type ZeroLogger struct {
	logger zerolog.Logger
}

func (l *ZeroLogger) Info() Event {
	return &ZeroEvent{
		event: l.logger.Info(),
	}
}

func (l *ZeroLogger) Error() Event {
	return &ZeroEvent{
		event: l.logger.Error(),
	}
}

func (l *ZeroLogger) Fatal() Event {
	return &ZeroEvent{
		event: l.logger.Fatal(),
	}
}

func (l *ZeroLogger) Panic() Event {
	return &ZeroEvent{
		event: l.logger.Panic(),
	}
}

var _ = Event(&ZeroEvent{})

type ZeroEvent struct {
	event *zerolog.Event
}

func (e *ZeroEvent) Err(err error) Event {
	return &ZeroEvent{
		event: e.event.Err(err),
	}
}

func (e *ZeroEvent) Msg(message string) {
	e.event.Msg(message)
}

func (e *ZeroEvent) Msgf(format string, v ...interface{}) {
	e.event.Msgf(format, v...)
}
