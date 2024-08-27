package handler

import (
	"github.com/Format-C-eft/universal-proxy/internal/proxy"
)

type StoreImpl struct {
	proxyStore proxy.Store
}

func New(
	proxyStore proxy.Store,
) *StoreImpl {
	return &StoreImpl{
		proxyStore: proxyStore,
	}
}
