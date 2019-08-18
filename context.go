package horde

import (
	"context"

	"github.com/chmking/horde/recorder"
)

type contextKey int

var recorderKey contextKey

func RecorderFrom(ctx context.Context) *recorder.Recorder {
	if v, ok := ctx.Value(recorderKey).(*recorder.Recorder); ok {
		return v
	}

	return nil
}
