package app

import (
	"github.com/ElianaTroper/cn/server/config"
	"github.com/ElianaTroper/cn/server/tree"

	"github.com/ElianaTroper/cn/server/app/benchmark"
	"github.com/ElianaTroper/cn/server/app/post"
)

func newConf(enabled bool) config.AppConf {
	return config.AppConf{Enable: enabled}
}

type App interface {
	Name() string
	Running() bool
	Stop() error
	Start(map[string]interface{}) error
	ErrChan() chan error
}

var running []App

// FUTURE: We should dynamically load apps, use RPC
// When that is done, add the ability to restart an app

func combineChan[t any](a, b chan t) chan t {
	// XXX: This assumes that the combined chans are still closed elsewhere
	combined := make(chan t)
	go func() {
		defer close(combined)
		for {
			select {
			case v, ok := <-a:
				if !ok {
					return
				}
				combined <- v
			case v, ok := <-b:
				if !ok {
					return
				}
				combined <- v
			}
		}
	}()
	return combined
}

func StartAll(apps map[string]config.AppConf) (chan error, err) {
	compiledApps := map[string]App{"Benchmark": benchmark.New(), "Post": post.New()}
	for key, val := range apps {
		if val.Enable {
			a, compiled := compiledApps[key]
			if !compiled {
				StopAll()
				return nil, fmt.Errorf("app not compiled: %v", key)
			}
			err := a.Start(val.AppSettings)
			if err != nil {
				a.Stop()
				StopAll()
				return nil, err
			}
			running = append(running, a)
		}
	}
	if running.Len() == 0 {
		return nil, nil
	} else if running.Len() == 1 {
		return running[0].ErrChan(), nil
	}
	combined := combineChan(running[0].ErrChan(), running[1].ErrChan())
	for _, a := range running[2:] {
		combined = combineChan(combined, a.ErrChan())
	}
	return combined, nil
}
func StopAll() error {
	// XXX: This only returns the first error that occurs (multiple may occur)
	var err error
	for _, a := range running {
		if err == nil {
			err = a.Stop()
			close(a.ErrChan())
		} else {
			_ = a.Stop()
			close(a.ErrChan())
		}
	}
	return err
}

// FUTURE: Add a deployment of the web app
func Deploy(appName, conf config.CnConf) error {
	file, err := os.ReadFile(conf.App.JsAppPath + '/' + appName)
	if err != nil {
		return err
	}
	// TODO: Add the updated root directly to the web app code
	return tree.Add("Root", "App/"+appName, file, conf)
}

func toggleEnable(conf *config.CnConf, appName string, enabled bool) {
	val, ok := *conf.App[appName]
	if !ok {
		*conf.App[appName] = newConf(enabled)
	} else {
		val.Enabled = enabled
	}
}

func Enable(appName, conf config.CnConf) error {
	toggleEnable(conf, appName, true)
	return daemon.Restart()
}

func Disable(appName, conf config.CnConf) error {
	toggleEnable(conf, appName, true)
	return daemon.Restart()
}
