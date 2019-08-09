package goteem

import "context"

type contextKey string

const recorder = contextKey("recorder")

func RecorderFrom(ctx context.Context) *Recorder {
	if v, ok := ctx.Value(recorder).(*Recorder); ok {
		return v
	}

	return nil
}
