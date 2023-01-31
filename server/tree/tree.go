package tree

import (
	"os"

	"github.com/ElianaTroper/cn/server/config"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
)

type Request struct {
	cnPath string
	data   []byte
	dir    string
	id     string
}

func (r *Request) path() string {
	return r.dir + "/" + r.id + ".req"
}
func (r *Request) swpPath() string {
	return r.dir + "/." + r.id + ".req"
}

func (r *Request) Remove() error {
	return os.Remove(r.path())
}

func (r *Request) write() error {
	err := os.MkdirAll(r.dir, 0700)
	if err != nil {
		return err
	}
	err = os.WriteFile(r.swpPath(), nil, 0600)
	if err != nil {
		return err
	}
	err = os.WriteFile(r.path(), r.data, 0600)
	if err != nil {
		return err
	}
	return os.Remove(r.swpPath())
}

// FUTURE: Add a timeout
// FUTURE: (Maybe) Turn this watcher stuff into a library
func (r *Request) watch() (*fsnotify.Watcher, chan error, error) {
	// This function watches the
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}
	errChan := make(chan error)
	go func() {
		defer watcher.Close()
		defer close(errChan)
		for {
			select {
			case event, ok := <-watcher.Events:
				if event.Has(fsnotify.Remove) {
					return
				}
				if !ok {
					return
				}
			case err, ok := <-watcher.Errors:
				if err != nil {
					errChan <- err
					return
				}
				if !ok {
					return
				}
			}
		}
	}()
	err = watcher.Add(r.path())
	if err != nil {
		watcher.Close()
		return nil, nil, err
	}
	return watcher, errChan, nil
}

func (r *Request) writeAndWatch() error {
	// Creates a request, and waits for it to be consumed
	watcher, watchChan, err := r.watch()
	if err != nil {
		return err
	}
	err = r.write()
	if err != nil {
		watcher.Close()
		return err
	}
	err, _ = <-watchChan // FUTURE: Add a timeout
	if err != nil {
		return err
	}
	return nil
}

func newRequest(cnPath, dir string, data []byte) *Request {
	return &Request{
		cnPath: cnPath,
		data:   data,
		dir:    dir,
		id:     uuid.New().String(),
	}
}

// FUTURE: Add updating based on ipns path
func Add(app, cnPath string, file []byte, conf config.CnConf) error {
	r := newRequest(cnPath, conf.Requests+"/"+app, file)
	return r.writeAndWatch()
}
