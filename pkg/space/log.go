package space

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(cl *ConfigLoader) *Logger {
	lvl := slog.LevelVar{}
	lvl.Set(slog.LevelInfo)
	if cl.injectConf.Log.Level != "" {
		_ = lvl.UnmarshalText([]byte(cl.injectConf.Log.Level))
	}
	return &Logger{
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     lvl.Level(),
			AddSource: true,
		})),
	}
}
