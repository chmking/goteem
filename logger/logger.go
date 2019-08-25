package logger

type Event interface {
	Err(error) Event
	Msg(string)
	Msgf(string, ...interface{})
}

type Logger interface {
	Info() Event
	Error() Event
	Fatal() Event
	Panic() Event
}
