package bootstrap

import (
	"context"

	"github.com/Format-C-eft/universal-proxy/internal/app"
	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/universal-proxy/internal/workers"
	"github.com/pkg/errors"
)

func InitializeServers(ctx context.Context) (app.Servers, error) {
	proxyStore, err := newProxyStoreImpl(ctx, config.GetProxyConfig())
	if err != nil {
		return nil, errors.Wrap(err, "newProxyStoreImpl")
	}

	appRestMiddlewareStore, err := newAppRestMiddlewareStoreImpl(config.GetAppRestMiddleware(), proxyStore)
	if err != nil {
		return nil, errors.Wrap(err, "newAppRestMiddlewareStoreImpl")
	}

	return app.New(
		config.GetAppAdmin(),
		newAppAdminHandlersStoreImpl(proxyStore),
		config.GetAppRest(),
		appRestMiddlewareStore,
		newAppRestHandlerStoreImpl(config.GetAppRestHandler()),
	), nil
}

func InitializeWorker(_ context.Context) (workers.Runnable, error) {
	return newWorkersRunnable(), nil
}
