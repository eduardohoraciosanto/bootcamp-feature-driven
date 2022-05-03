package logger

import (
	"context"
	"runtime/debug"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/config"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Logger interface {
	Info(ctx context.Context, message string)
	Error(ctx context.Context, message string)
	Debug(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Sync() error
}

type logger struct {
	internal       *zap.SugaredLogger
	tracingEnabled bool
}

// NewLogger initializes de logger and allows the usage of the Logger Interface
// user MUST defer the call to Sync() immediately afterwards.

func NewLogger(service string, tracingEnabled bool) Logger {
	l, err := zap.NewProduction()
	if err != nil {
		panic("unable to initialize logger: " + err.Error())
	}

	return &logger{
		internal: l.WithOptions(zap.AddCallerSkip(1)).Sugar().
			With("version", config.GetVersion()).
			With("service", service),
		tracingEnabled: tracingEnabled,
	}
}

// Sync MUST be defered to flush any buffered logs prior to shutting down the application.
func (l *logger) Sync() error {
	return l.internal.Sync()
}

// Info allows for a message with info lever to be logged
func (l *logger) Info(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Info(message)
		return
	}
	l.internal.Info(message)
}

// Error allows for a message with error lever to be logged
func (l *logger) Error(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Error(message)
		return
	}
	l.internal.Error(message)
}

// Debug allows for a message with debug lever to be logged
func (l *logger) Debug(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Debug(message)
		return
	}
	l.internal.Debug(message)
}

// Warn allows for a message with warn lever to be logged
func (l *logger) Warn(ctx context.Context, message string) {
	if l.tracingEnabled {
		l.injectTracing(ctx).Warn(message)
		return
	}
	l.internal.Warn(message)

}

// WithField allows for the inclusion of a key-value into the log
func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{
		internal: l.internal.With(key, value),
	}
}

// WithError allows for the inclusion of an error into the log. It also populates DDog stack field
func (l *logger) WithError(err error) Logger {
	newLogger := l.internal.With(
		"error", err,
	)

	if l.tracingEnabled {
		newLogger = newLogger.With("dd.error.stack", string(debug.Stack()))
	}

	return &logger{
		internal: newLogger,
	}
}

// injectTracing enters correlation ID, if any, and DDog tracing information into the log.
func (l *logger) injectTracing(ctx context.Context) *zap.SugaredLogger {
	//add our correlation id if present
	cid := ctx.Value("correlation_id")
	entry := l.internal
	if cid != nil {
		entry = entry.With("correlation_id", ctx.Value("correlation_id"))
	}

	//add datadog information if available
	if span, ok := tracer.SpanFromContext(ctx); ok {
		entry = entry.With("dd.trace_id", span.Context().TraceID()).
			With("dd.span_id", span.Context().SpanID())
	}

	return entry
}
