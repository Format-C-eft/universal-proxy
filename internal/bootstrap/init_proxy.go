package bootstrap

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/universal-proxy/internal/proxy"
)

var proxyStoreImplOnce = sync.Once{}
var proxyStoreImpl *proxy.StoreImpl

func newProxyStoreImpl(ctx context.Context, cfg config.ProxyConfig) (*proxy.StoreImpl, error) {
	var err error
	proxyStoreImplOnce.Do(func() {
		proxyStoreImpl, err = proxy.New(ctx, cfg)
		if err != nil {
			err = errors.Wrap(err, "proxy.New")
		}
	})

	return proxyStoreImpl, err
}
