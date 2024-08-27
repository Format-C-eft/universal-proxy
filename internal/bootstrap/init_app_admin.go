package bootstrap

import (
	"sync"

	"github.com/Format-C-eft/universal-proxy/internal/app/admin/handler"
	"github.com/Format-C-eft/universal-proxy/internal/proxy"
)

var appAdminHandlersStoreImplOnce sync.Once
var appAdminHandlersStoreImpl *handler.StoreImpl

func newAppAdminHandlersStoreImpl(
	proxyStore proxy.Store,
) *handler.StoreImpl {
	appAdminHandlersStoreImplOnce.Do(func() {
		appAdminHandlersStoreImpl = handler.New(
			proxyStore,
		)
	})

	return appAdminHandlersStoreImpl
}
