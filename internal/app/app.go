package app

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"gin-casbin/internal/app/config"
	"gin-casbin/pkg/logger"

	_ "gin-casbin/internal/app/swagger"
)

type options struct {
	ConfigFile string
	ModelFile  string
}

// Option
type Option func(*options)

// SetModelFile set casbin model
func SetModelFile(s string) Option {
	return func(o *options) {
		o.ModelFile = s
	}
}

// Init
func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	config.MustLoad(o.ConfigFile)
	if v := o.ModelFile; v != "" {
		config.C.Casbin.Model = v
	}
	return func() {
	}, nil
}

// Run
func Run(ctx context.Context, opts ...Option) error {
	var state int32 = 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.Printf(ctx, "Signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			atomic.CompareAndSwapInt32(&state, 1, 0)
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.Printf(ctx, "Quit")
	time.Sleep(time.Second)
	os.Exit(int(atomic.LoadInt32(&state)))
	return nil
}
