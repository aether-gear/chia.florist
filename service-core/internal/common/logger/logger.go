package logger

import "context"

type Logger interface {
	Info(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Debug(ctx context.Context, msg string, fields ...Field)

	With(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value any
}
