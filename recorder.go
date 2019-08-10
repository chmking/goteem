package horde

type Recorder struct{}

func (r *Recorder) Success(method, path string, latency int64) {
}

func (r *Recorder) Error(method, path string, latency int64, err error) {
}
