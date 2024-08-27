package middleware

import (
	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/universal-proxy/internal/proxy"
	"github.com/Format-C-eft/utils/cache"
	"github.com/Format-C-eft/utils/cache/local_lru"
	"github.com/pkg/errors"
)

type StoreImpl struct {
	sessionCache cache.Store
	proxyStore   proxy.Store

	cfg config.AppRestMiddleware
}

func New(
	cfg config.AppRestMiddleware,
	proxyStore proxy.Store,
) (*StoreImpl, error) {
	sessionCache, err := local_lru.NewClient(cfg.SessionCache.Name, cfg.SessionCache.Size, cfg.SessionCache.TTL)
	if err != nil {
		return nil, errors.Wrap(err, "local_lru.NewClient")
	}

	return &StoreImpl{
		sessionCache: sessionCache,
		proxyStore:   proxyStore,

		cfg: cfg,
	}, nil
}
