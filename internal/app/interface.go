package app

import (
	"context"
)

type Servers interface {
	Run(ctx context.Context)
	Stop(ctx context.Context)
}
