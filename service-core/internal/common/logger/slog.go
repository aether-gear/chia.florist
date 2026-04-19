package logger

import (
	"context"
	"log/slog"
	"os"
)

type slogLogger struct {
	log *slog.Logger
}

func New() Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &slogLogger{
		log: slog.New(handler),
	}
}

func (l *slogLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log.InfoContext(ctx, msg, toAttrs(fields)...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, fields ...Field) {
	l.log.ErrorContext(ctx, msg, toAttrs(fields)...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log.WarnContext(ctx, msg, toAttrs(fields)...)
}

func (l *slogLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log.DebugContext(ctx, msg, toAttrs(fields)...)
}

func (l *slogLogger) With(fields ...Field) Logger {
	return &slogLogger{
		log: l.log.With(toAttrs(fields)...),
	}
}

func toAttrs(fields []Field) []any {
	attrs := make([]any, 0, len(fields))
	for _, f := range fields {
		attrs = append(attrs, f.Key, f.Value)
	}
	return attrs
}
