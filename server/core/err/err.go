package err

import "fmt"

var (
	Running = fmt.Errorf("already running")
	Watcher = fmt.Errorf("an error occurred while watching files")
)
