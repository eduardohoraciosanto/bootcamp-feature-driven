package logger

import (
	"context"
	"reflect"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Logger interface {
	Info(ctx context.Context, message string)
	Error(ctx context.Context, message string)
	Debug(ctx context.Context, message string)
	Warn(ctx context.Context, message string)
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
}

type logger struct {
	log *logrus.Entry
}

func NewLogger(jsonFormatted bool, service string) Logger {
	l := logrus.New()
	if jsonFormatted {
		l.SetFormatter(&logrus.JSONFormatter{})
	} else {
		l.SetFormatter(&logrus.TextFormatter{})
	}

	return &logger{
		log: l.WithField("log_svc", service),
	}
}

func (l *logger) Info(ctx context.Context, message string) {
	l.injectTracingInfo(ctx).Info(message)
}
func (l *logger) Error(ctx context.Context, message string) {
	l.injectTracingInfo(ctx).WithField("dd.error.msg", message).Error(message)
}
func (l *logger) Debug(ctx context.Context, message string) {
	l.injectTracingInfo(ctx).Debug(message)
}
func (l *logger) Warn(ctx context.Context, message string) {
	l.injectTracingInfo(ctx).Warn(message)
}
func (l *logger) WithField(key string, value interface{}) Logger {
	return &logger{
		log: l.log.WithField(key, value),
	}
}
func (l *logger) WithError(err error) Logger {
	return &logger{
		log: l.log.WithField("dd.error.type", reflect.TypeOf(err)).
			WithField("dd.error.stack", string(debug.Stack())).WithField("error", err),
	}
}

func (l *logger) injectTracingInfo(ctx context.Context) *logrus.Entry {
	//add our correlation id if present
	cid := ctx.Value("correlation_id")
	entry := l.log
	if cid != nil {
		entry = entry.WithField("correlation_id", ctx.Value("correlation_id"))
	}

	//add datadog information if available
	if span, ok := tracer.SpanFromContext(ctx); ok {
		entry = entry.WithField("dd.trace_id", span.Context().TraceID()).
			WithField("dd.span_id", span.Context().SpanID())
	}

	return entry
}
