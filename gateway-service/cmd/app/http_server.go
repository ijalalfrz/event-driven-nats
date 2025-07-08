package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/config"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/endpoint"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/router"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/service"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/lang"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var timeout = 30 * time.Second

var httpServerCmd = &cobra.Command{
	Use:   "http",
	Short: "Serve incoming requests from REST HTTP/JSON API",
	Run: func(_ *cobra.Command, _ []string) {
		slog.Debug("command line flags", slog.String("config_path", cfgFilePath))
		cfg := config.MustInitConfig(cfgFilePath)

		logger.InitStructuredLogger(cfg.LogLevel)

		runHTTPServer(cfg)
	},
}

func runHTTPServer(cfg config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var waitGroup sync.WaitGroup

	slog.InfoContext(ctx, "starting...", slog.String("log_level", string(cfg.LogLevel)))

	// Starts the server in a go routine
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		startHTTPServer(ctx, cfg)
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigChannel:
		cancel()
		slog.InfoContext(ctx, "received OS signal. Exiting...", slog.String("signal", sig.String()))
	case <-ctx.Done():
		slog.ErrorContext(ctx, "failed to start HTTP server")
	}

	waitGroup.Wait()
	slog.Info("All servers stopped")
}

// startHTTPServer loads config and starts HTTP server.
func startHTTPServer(ctx context.Context, cfg config.Config) {
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		startMainApp(ctx, cfg)
	}()

	if cfg.HTTP.PprofEnabled {
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()
			startPprof(ctx, cfg)
		}()
	}

	waitGroup.Wait()
}

func startMainApp(ctx context.Context, cfg config.Config) {
	lang.SetSupportedLanguages(cfg.Locales.SupportedLanguages)
	lang.SetBasePath(cfg.Locales.BasePath)

	endpts := makeEndpoints(cfg)

	router := router.MakeHTTPRouter(
		endpts,
		cfg,
	)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		WriteTimeout: cfg.HTTP.Timeout,
		ReadTimeout:  cfg.HTTP.Timeout,
	}

	slog.Info("running HTTP server...", slog.Int("port", cfg.HTTP.Port))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	// shutdown ctx
	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown HTTP server", slog.String("error", err.Error()))
	}

	slog.Info("HTTP server gracefully stopped")
}

func makeEndpoints(cfg config.Config) endpoint.Endpoint {
	// init all service clients
	userServiceClient := service.NewUserServiceClient(cfg.UserService.URL,
		service.WithMaxRetries(cfg.UserService.MaxRetry),
		service.WithTimeout(cfg.UserService.Timeout),
	)
	listingViewServiceClient := service.NewListingViewServiceClient(cfg.ListingViewService.URL,
		service.WithMaxRetries(cfg.ListingViewService.MaxRetry),
		service.WithTimeout(cfg.ListingViewService.Timeout),
	)
	listingServiceClient := service.NewListingServiceClient(cfg.ListingService.URL,
		service.WithMaxRetries(cfg.ListingService.MaxRetry),
		service.WithTimeout(cfg.ListingService.Timeout),
	)

	return endpoint.Endpoint{
		PublicUser: makePublicUserEndpoints(userServiceClient),
		PublicListing: makePublicListingEndpoints(listingViewServiceClient,
			listingServiceClient, userServiceClient),
	}
}

func makePublicUserEndpoints(
	userServiceClient *service.UserServiceClient,
) endpoint.PublicUser {
	userSvc := service.NewPublicUserService(userServiceClient)

	return endpoint.NewPublicUserEndpoint(userSvc)
}

func makePublicListingEndpoints(
	listingViewServiceClient *service.ListingViewServiceClient,
	listingServiceClient *service.ListingServiceClient,
	userServiceClient *service.UserServiceClient,
) endpoint.PublicListing {
	listingSvc := service.NewPublicListingService(listingViewServiceClient,
		listingServiceClient, userServiceClient)

	return endpoint.NewPublicListingEndpoint(listingSvc)
}

func startPprof(ctx context.Context, cfg config.Config) {
	// manually register pprof handlers with custom path.
	http.HandleFunc("/internal/pprof/", pprof.Index)
	http.HandleFunc("/internal/pprof/cmdline", pprof.Cmdline)
	http.HandleFunc("/internal/pprof/profile", pprof.Profile)
	http.HandleFunc("/internal/pprof/symbol", pprof.Symbol)
	http.HandleFunc("/internal/pprof/trace", pprof.Trace)
	slog.Info("running pprof...", slog.Int("port", cfg.HTTP.PprofPort))

	server := &http.Server{
		Addr:              fmt.Sprintf("localhost:%d", cfg.HTTP.PprofPort),
		ReadHeaderTimeout: cfg.HTTP.Timeout,
	}

	// limit access by binding port to localhost only.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("pprof server error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown pprof server", slog.String("error", err.Error()))
	}

	slog.Info("pprof server stopped")
}
