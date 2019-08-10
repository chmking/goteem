package horde

const ErrStatusStopping = Error("agent is in the process of stopping")
const ErrStatusQuitting = Error("agent is in the process of quitting")

type Error string

func (e Error) Error() string {
	return string(e)
}
