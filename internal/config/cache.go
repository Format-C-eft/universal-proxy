package config

import (
	"time"
)

type LocalCache struct {
	Name string
	Size int
	TTL  time.Duration
}
