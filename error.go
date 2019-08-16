package horde

const (
	ErrStatusUnexpected = Error("unexpected status")
	ErrStatusUnknown    = Error("invalid status switch from unknown")
	ErrStatusIdle       = Error("invalid status switch from idle")
	ErrStatusScaling    = Error("invalid status switch from scaling")
	ErrStatusRunning    = Error("invalid status switch from running")
	ErrStatusStopping   = Error("invalid status switch from stopping")
	ErrStatusQuitting   = Error("invalid status switch from quitting")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
