package daemon

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/ElianaTroper/cn/server/app"
	"github.com/ElianaTroper/cn/server/config"
	"github.com/ElianaTroper/cn/server/core"
	"github.com/ElianaTroper/cn/server/core/err"
	"github.com/ElianaTroper/pidfile"

	"github.com/fsnotify/fsnotify"
)

var (
	StopIfAppCrash = true // FUTURE: set via config
	ErrNotRunning  = fmt.Errorf("process is not running")
)

func run(ctx context.Context, conf config.CnConf, out io.Writer) error {

	c, err := core.New(conf)
	if err != nil {
		return err
	}
	coreErrChan, err := c.Start()
	if err != nil {
		return err
	}
	defer c.Stop()
	appErrs, err := app.StartAll(conf.App)
	if err != nil {
		return err
	}
	defer app.StopAll()

	for {
		select {
		case <-ctx.Done():
			// Cleanup
			err := app.StopAll()
			if err != nil {
				return err
			}
			return c.Stop()
		case err := <-coreErrChan:
			return err
		case err := <-appErrs:
			if StopIfAppCrash {
				return err
			} // FUTURE: Handle if this is disabled
		}
	}
}

func launch(conf config.CnConf) error {

	err := pidfile.Write(conf.Pid)
	if err != nil {
		return err
	}

	// Cancellable context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Intercepting signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	defer func() {
		signal.Stop(signalChan)
		close(signalChan)
		cancel()
	}()
	go func() {
		select {
		case s := <-signalChan:
			if s != syscall.SIGUSR1 {
				fmt.Fprintf(os.Stderr, "SIGINT/SIGTERM, exiting")
			} else {
				fmt.Fprintf(os.Stdout, "recieved stop signal, exiting")
			}
			cancel()
		case <-ctx.Done():
		}
	}()

	err = run(ctx, conf, os.Stdout)
	if err != nil {
		pidfile.Remove(conf.Pid)
		return err
	}
	return pidfile.Remove(conf.Pid)
}

func Start(conf config.CnConf) error {
	return launch(conf)
}

func Stop(conf config.CnConf) error {
	pid, running, err := pidfile.PidIsRunning(conf.Pid)
	if err != nil {
		return err
	} else if !running {
		return ErrNotRunning
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	watcher.Add(conf.Pid)
	err = syscall.Kill(pid, syscall.SIGUSR1)
	if err != nil {
		return err
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if event.Has(fsnotify.Remove) {
				return nil
			}
			if !ok {
				return nil
			}
		case err, _ := <-watcher.Errors:
			if err != nil {
				return err
			}
			return nil
		} // FUTURE: Add a timeout error
	}
}

func Restart(conf config.CnConf) error {
	err := Stop(conf)
	if err != nil && err != ErrNotRunning {
		return err
	}
	return Start(conf)
}
