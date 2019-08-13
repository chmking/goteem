package horde

import "context"

type contextKey int

var recorderKey contextKey

func RecorderFrom(ctx context.Context) *Recorder {
	if v, ok := ctx.Value(recorderKey).(*Recorder); ok {
		return v
	}

	return nil
}
