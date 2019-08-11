package session

type Callback func()

type Session struct {
}

func (s *Session) Scale(count int32, rate float64, wait int64, cb Callback) {
	cb()
}

func (s *Session) Stop(cb Callback) {
	cb()
}
