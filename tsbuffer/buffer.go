package tsbuffer

import (
	"sync"
	"time"

	"github.com/chmking/horde/protobuf/private"
)

type Millisecond func() int64

var unixMillisecond = func() int64 {
	return time.Now().UnixNano() / 1e6
}

func New(window time.Duration) *Buffer {
	return newInjected(window, unixMillisecond)
}

func newInjected(window time.Duration, millisecond Millisecond) *Buffer {
	windowMilliseconds := window.Nanoseconds() / 1e6
	return &Buffer{
		millisecond: millisecond,
		window:      windowMilliseconds,
		pointer:     millisecond() - windowMilliseconds,
		results:     make(map[int64][]*private.Result, 0),
	}
}

type Buffer struct {
	millisecond Millisecond
	window      int64
	pointer     int64
	results     map[int64][]*private.Result
	mtx         sync.Mutex
}

func (b *Buffer) Add(result *private.Result) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	ts := result.Millisecond
	if result == nil || ts < b.pointer {
		return
	}

	b.results[ts] = append(b.results[ts], result)
}

func (b *Buffer) Collect() map[int64][]*private.Result {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	var results map[int64][]*private.Result

	dest := b.millisecond() - b.window
	diff := dest - b.pointer

	for i := int64(0); i < diff; i++ {
		current := b.pointer + i
		values, ok := b.results[current]
		if !ok {
			continue
		}

		if results == nil {
			results = make(map[int64][]*private.Result, 0)
		}

		results[current] = values
		delete(b.results, current)
	}

	b.pointer = dest

	return results
}
