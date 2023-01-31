package core

import (
	"sync"
	"time"

	"github.com/ElianaTroper/cn/server/config"
	"github.com/ElianaTroper/cn/server/core/err"
	"github.com/ElianaTroper/cn/server/tree"

	"github.com/fsnotify/fsnotify"
)

// TODO: Use accounts to point to pubkeys (initially)
/* FUTURE: add accounts elsewhere
type Account struct {
	Name string // FUTURE: Add a pubkey later
}
*/

type NodePointer struct {
	Ipfs string // Points to the last seen ipfs hash of the node
	Ipns string // Points to the canonical ipns address of the node
}

type Node struct {
	Index    uint64       // Starts at 0, goes up by 1 if the node is iterated on
	Previous *NodePointer `json:",omitempty"` // Points at the previous version of the node - optional, depending on the application or if index == 0
	// Owner *Account `json:",omitempty"` // FUTURE: Move accounts to another place
} // FUTURE: Tell node to only maintain a max depth at certain points
// FUTURE: Add a time item
// FUTURE: Add a signature

// The root node of the system
type Root struct {
	Node
	App []NodePointer
	// Source map[string]NodePointer FUTURE: Add the source at the root
}

type Core struct {
	err     chan error
	running bool
	conf    config.CnConf
	done    chan bool
	root    *Root
	mu      sync.Mutex
	updates []tree.Request
}

const (
	rootNodeFile = "root.cn" // TODO: Set in the config
	runManager   = true      // FUTURE: Set in the config
)

func (c *Core) Process() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: Get all files in dir
	// TODO: For all non-swp files in dir, check if there's a swp file
	// TODO: If there isn't a swp file, and it's not in updates, process
	// TODO: If it's newer than the root, update the local root and add to ipfs
	// XXX: This works well for the root (no intermediate nodes needed)
	//		but may have issues if used directly for, e.g. post
}

// FUTURE: Add a "client" mode
func (c *Core) startManager() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	go func() {
		// This function watches for new additions to the root
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if event.Has(fsnotify.Write) {
					err := c.Process()
					if err != nil {
						c.err <- err
						return
					}
				}
				if !ok {
					return
				}
			case err, _ := <-watcher.Errors:
				if err != nil {
					c.err <- err
				}
				return
			case <-c.done:
				return
			}
		}
	}()
	err = watcher.Add(c.conf.Requests + "/Root")
	if err != nil {
		watcher.Close()
		return err
	}
	ticker := time.NewTicker(time.Duration(c.conf.Root.Tick) * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := func() error {
					c.mu.Lock()
					defer c.mu.Unlock()
					if len(c.updates) > 0 {
						inErr := c.root.deploy()
						if inErr != nil {
							return inErr
						}
						for _, u := range c.updates {
							inErr = u.Remove()
							if inErr != nil {
								return inErr
							} // XXX: We still may want to remove the files?
						}
						c.updates = nil
					}
				}()
				if err != nil {
					c.err <- err
					return
				}
			case <-c.done:
				return
			}
		}
	}()
	return nil
}

func (c *Core) Start() (chan error, error) {
	if c.running {
		return c.errChan, err.Running
	}
	c.errChan = make(chan error)
	err := c.startManager()
	if err != nil {
		c.Stop()
		return nil, err
	}
	c.running = true
	return c.errChan, nil
}

func (c *Core) Stop() error {
	c.done <- true
	close(c.done)
	close(c.err)
	c.running = false
	return nil
}

func (c *Core) loadRoot() (*Root, error) {
	rootFile, err := os.ReadFile(c.conf.Storage.Path + "/" + rootNodeFile)
	// TODO: Also walk back until Retain is reached, whether
	//		 that's over ipfs or locally
	// TODO: Also check that root is pinned properly
	if err != nil {
		return nil, err
	}
	var res Root
	err = json.Unmarshall(rootfile, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *Root) deploy() error {
	// TODO: This
}

func (c *Core) startRoot() (*Root, error) {
	// FUTURE: Allow init from ipfs if you're not the main node
	root, err := c.loadRoot()
	if err == nil {
		c.root = root
		return c.root, nil
	} else if os.IsNotExist(err) {
		root = &Root{}
		err = root.deploy()
		if err != nil {
			return nil, err
		}
		c.root = root
		return c.root, err
	}
	return nil, err
}

func (c *Core) getRoot() (*Root, error) {
	if c.root == nil {
		return c.startRoot()
	}
	return c.root, nil
}

func New(conf config.CnConf) (*Core, error) {
	c := &Core{
		conf: conf,
	}
	_, err := c.getRoot()
	if err != nil {
		return nil, err
	}
	return c, nil
}
