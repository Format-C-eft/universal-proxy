package workers

import (
	"context"
)

type Runnable interface {
	Run(ctx context.Context)
	Stop(ctx context.Context)
}
