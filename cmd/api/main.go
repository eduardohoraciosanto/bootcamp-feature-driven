package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/cache"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/config"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/cart"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/health"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
	transport "github.com/eduardohoraciosanto/bootcamp-feature-driven/transport/http"
	"github.com/go-redis/redis/v8"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	tracer.Start()
	defer tracer.Stop()

	conf := config.New()

	l := logger.NewLogger(conf.JSONLogs, "shopping cart api")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.RedisServer,
		Password: conf.RedisPassword,
	})

	cacheClient := cache.NewRedisCache(
		l.WithField("svc", "cache"),
		0,
		redisClient,
	)

	isvc := item.NewExternalService(l.WithField("svc", "external service"), &http.Client{
		Timeout: time.Second * 10,
	})

	hsvc := health.NewService(
		cacheClient,
		isvc,
		l.WithField("svc", "health service"),
	)

	csvc := cart.NewCartService(
		config.GetVersion(),
		l.WithField("svc", "cart service"),
		cacheClient,
		isvc,
	)

	httpTransportRouter := transport.NewHTTPRouter(hsvc, csvc, isvc)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", conf.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      httpTransportRouter,
	}
	l.WithField(
		"transport", "http").
		WithField(
			"port", conf.Port).
		Info(context.Background(), "Transport Start")
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.WithField(
				"transport", "http").
				WithError(err).
				Info(context.Background(), "Transport Stopped")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	l.Info(context.Background(), "Service gracefully shutted down")
	os.Exit(0)
}
