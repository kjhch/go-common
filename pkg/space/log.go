package space

import (
	"context"
	"log/slog"
	"os"
)

const KeyRequestID = "requestID"

type Logger struct {
	*slog.Logger
}

func NewLogger(cl *ConfigLoader) *Logger {
	lvl := slog.LevelVar{}
	lvl.Set(slog.LevelInfo)
	if cl.injectConf.Log.Level != "" {
		_ = lvl.UnmarshalText([]byte(cl.injectConf.Log.Level))
	}
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvl.Level(),
		AddSource: true,
	})
	return &Logger{slog.New(&ContextHandler{jsonHandler})}
}

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if ctx != nil {
		if requestId, ok := ctx.Value(KeyRequestID).(string); ok && requestId != "" {
			record.Add(KeyRequestID, requestId)
		}
	}
	return h.Handler.Handle(ctx, record)
}
