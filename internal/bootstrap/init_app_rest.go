package bootstrap

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/Format-C-eft/universal-proxy/internal/app/rest/handler"
	"github.com/Format-C-eft/universal-proxy/internal/app/rest/middleware"
	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/universal-proxy/internal/proxy"
)

var appRestMiddlewareStoreImplOnce sync.Once
var appRestMiddlewareStoreImpl *middleware.StoreImpl

func newAppRestMiddlewareStoreImpl(
	cfg config.AppRestMiddleware,
	proxyStore proxy.Store,
) (*middleware.StoreImpl, error) {
	var err error
	appRestMiddlewareStoreImplOnce.Do(func() {
		appRestMiddlewareStoreImpl, err = middleware.New(
			cfg,
			proxyStore,
		)
		if err != nil {
			err = errors.Wrap(err, "middleware.New")
		}
	})

	return appRestMiddlewareStoreImpl, err
}

var appRestHandlerStoreImplOnce sync.Once
var appRestHandlerStoreImpl *handler.StoreImpl

func newAppRestHandlerStoreImpl(
	cfg config.AppRestHandler,
) *handler.StoreImpl {
	appRestHandlerStoreImplOnce.Do(func() {
		appRestHandlerStoreImpl = handler.New(
			cfg,
		)
	})

	return appRestHandlerStoreImpl
}
