package logger

import (
	"fmt"
	"io"
	"log"
	"strings"
)

var _ = Logger(&StdLogger{})

func NewStdLogger(writer io.Writer) *StdLogger {
	return &StdLogger{
		logger: log.New(writer, "", log.LstdFlags),
		writer: writer,
	}
}

type StdLogger struct {
	logger *log.Logger
	writer io.Writer
}

func (l *StdLogger) Info() Event {
	return &StdEvent{
		level:  "info",
		logger: log.New(l.writer, "INFO :: ", log.LstdFlags),
	}
}

func (l *StdLogger) Error() Event {
	return &StdEvent{
		level:  "error",
		logger: log.New(l.writer, "ERROR :: ", log.LstdFlags),
	}
}

func (l *StdLogger) Fatal() Event {
	return &StdEvent{
		level:  "fatal",
		logger: log.New(l.writer, "FATAL :: ", log.LstdFlags),
	}
}

func (l *StdLogger) Panic() Event {
	return &StdEvent{
		level:  "panic",
		logger: log.New(l.writer, "PANIC :: ", log.LstdFlags),
	}
}

var _ = Event(&StdEvent{})

type StdEvent struct {
	level    string
	logger   *log.Logger
	elements []string
}

func (e *StdEvent) Err(err error) Event {
	e.elements = append(e.elements, fmt.Sprintf("error: %s", err))
	return e
}

func (e *StdEvent) Msg(message string) {
	output := strings.Join(append(e.elements, fmt.Sprintf("message: %s", message)), " = ")

	switch e.level {
	case "fatal":
		e.logger.Fatal(output)
	case "panic":
		e.logger.Panic(output)
	default:
		e.logger.Println(output)
	}
}

func (e *StdEvent) Msgf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	e.Msg(message)
}
