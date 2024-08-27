package bootstrap

import (
	"sync"

	"github.com/Format-C-eft/universal-proxy/internal/workers"
)

var workersRunnableOnce sync.Once
var workersRunnable workers.Runnable

func newWorkersRunnable() workers.Runnable {
	workersRunnableOnce.Do(func() {
		workersRunnable = workers.New()
	})

	return workersRunnable
}
