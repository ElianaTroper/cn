package benchmark

import (
	"github.com/ElianaTroper/cn/server/config"

	"github.com/ElianaTroper/cn/server/core/err"
)

type benchmark struct {
	running bool
	config  config.AppConf
}

func (b *benchmark) Name() string {
	return "benchmark"
}

func (b *benchmark) Running() bool {
	return b.running
}

func (b *benchmark) Start(config.AppConf) (chan error, error) {
	if b.running {
		return nil, err.Running
	}
	// TODO: This
}

func (b *benchmark) Stop() error {
	if !b.running {
		return nil
	}
	// TODO: This
}

func New() benchmark {

}
