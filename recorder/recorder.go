package recorder

import (
	"time"

	"github.com/chmking/horde/protobuf/private"
	"github.com/chmking/horde/protobuf/public"
)

func New() *Recorder {
	return &Recorder{
		results: make(chan private.Result, 1000),
	}
}

type Recorder struct {
	results chan private.Result
}

func (r *Recorder) Success(method, path string, latency int64) {
	r.record(public.Code_CODE_SUCCESS, method, path, latency, nil)
}

func (r *Recorder) Error(method, path string, latency int64, err error) {
	r.record(public.Code_CODE_ERROR, method, path, latency, err)
}

func (r *Recorder) Panic(method, path string, latency int64, err error) {
	r.record(public.Code_CODE_PANIC, method, path, latency, err)
}

func (r *Recorder) Results() []*private.Result {
	var results []*private.Result

	for {
		select {
		case result := <-r.results:
			results = append(results, &result)
		default:
			return results
		}
	}
}

func (r *Recorder) record(code public.Code, method, path string, latency int64, err error) {
	result := private.Result{
		Second:  time.Now().Unix(),
		Code:    code,
		Method:  method,
		Path:    path,
		Latency: latency,
	}

	if err != nil {
		result.Message = err.Error()
	}

	r.results <- result
}
