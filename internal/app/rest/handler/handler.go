package handler

import (
	"github.com/Format-C-eft/universal-proxy/internal/config"
)

type StoreImpl struct {
	cfg config.AppRestHandler
}

func New(
	cfg config.AppRestHandler,
) *StoreImpl {
	return &StoreImpl{
		cfg: cfg,
	}
}
