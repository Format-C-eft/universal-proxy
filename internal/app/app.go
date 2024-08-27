package app

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	appAdmin "github.com/Format-C-eft/universal-proxy/internal/app/admin"
	appAdminHandlers "github.com/Format-C-eft/universal-proxy/internal/app/admin/handler"
	appRest "github.com/Format-C-eft/universal-proxy/internal/app/rest"
	appRestHandlers "github.com/Format-C-eft/universal-proxy/internal/app/rest/handler"
	appRestMiddleware "github.com/Format-C-eft/universal-proxy/internal/app/rest/middleware"
	appStatus "github.com/Format-C-eft/universal-proxy/internal/app/status"
	"github.com/Format-C-eft/universal-proxy/internal/config"
	"github.com/Format-C-eft/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type servers struct {
	admin  *http.Server
	rest   *http.Server
	status *http.Server

	isReady *atomic.Value
}

func New(
	appAdminConfig config.AppAdmin,
	appAdminHandlersStore appAdminHandlers.Store,
	appRestConfig config.AppRest,
	appRestMiddlewareStore appRestMiddleware.Store,
	appRestHandlerStore appRestHandlers.Store,
) Servers {
	mode := gin.ReleaseMode
	if config.LocalRunFlag {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)

	isReady := &atomic.Value{}
	isReady.Store(false)

	return &servers{
		admin: &http.Server{
			Addr:              fmt.Sprintf(":%d", config.GetPort("admin")),
			Handler:           appAdmin.New(appAdminConfig, appAdminHandlersStore),
			ReadHeaderTimeout: time.Second * 2,
		},
		rest: &http.Server{
			Addr:              fmt.Sprintf(":%d", config.GetPort("rest")),
			Handler:           appRest.New(appRestMiddlewareStore, appRestHandlerStore, appRestConfig),
			ReadHeaderTimeout: time.Second * 2,
		},
		status: &http.Server{
			Addr:              fmt.Sprintf(":%d", config.GetPort("status")),
			Handler:           appStatus.New(isReady),
			ReadHeaderTimeout: time.Second * 2,
		},
		isReady: isReady,
	}
}

func (s *servers) Run(ctx context.Context) {
	go func() {
		logger.InfoF(ctx, "Admin Server is listening on port: %s", s.admin.Addr)
		if err := s.admin.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.FatalKV(ctx, "Failed running admin server", "err", err)
		}
	}()

	go func() {
		logger.InfoF(ctx, "Rest Server is listening on port: %s", s.rest.Addr)
		if err := s.rest.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.FatalKV(ctx, "Failed running rest server", "err", err)
		}
	}()

	go func() {
		logger.InfoF(ctx, "Status Server is listening on port: %s", s.status.Addr)
		if err := s.status.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.FatalKV(ctx, "Failed running status server", "err", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		s.isReady.Store(true)
		logger.Info(ctx, "The service is ready to accept requests")
	}()
}

func (s *servers) Stop(ctx context.Context) {
	s.isReady.Store(false)

	if err := s.admin.Shutdown(ctx); err != nil {
		logger.FatalKV(ctx, "admin.Shutdown", "err", err)
	} else {
		logger.Info(ctx, "Admin server shut down correctly")
	}

	if err := s.rest.Shutdown(ctx); err != nil {
		logger.FatalKV(ctx, "rest.Shutdown", "err", err)
	} else {
		logger.Info(ctx, "Rest server shut down correctly")
	}

	if err := s.status.Shutdown(ctx); err != nil {
		logger.FatalKV(ctx, "status.Shutdown", "err", err)
	} else {
		logger.Info(ctx, "Status server shut down correctly")
	}
}
