package workers

import (
	"context"
)

type RunnableImpl struct {
	workers []Runnable
}

func New(workers ...Runnable) *RunnableImpl {
	return &RunnableImpl{
		workers: workers,
	}
}

func (r *RunnableImpl) Run(ctx context.Context) {
	for _, worker := range r.workers {
		worker.Run(ctx)
	}
}

func (r *RunnableImpl) Stop(ctx context.Context) {
	for _, worker := range r.workers {
		worker.Stop(ctx)
	}
}
