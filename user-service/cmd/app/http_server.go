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

	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/config"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/endpoint"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/repository"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/router"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/app/service"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/db"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/lang"
	"github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/logger"
	natstransport "github.com/ijalalfrz/event-driven-nats/user-service/internal/pkg/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
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
		startHTTPServer(ctx, cfg, cancel)
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
func startHTTPServer(ctx context.Context, cfg config.Config,
	cancel context.CancelFunc) {
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		if err := startMainApp(ctx, cfg); err != nil {
			slog.Error("failed to start main app", slog.String("error", err.Error()))
			cancel()
		}
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

func startMainApp(ctx context.Context, cfg config.Config) error {
	lang.SetSupportedLanguages(cfg.Locales.SupportedLanguages)
	lang.SetBasePath(cfg.Locales.BasePath)

	natsConn, err := nats.Connect(cfg.Nats.URL)
	if err != nil {
		slog.Error("failed to connect to NATS", slog.String("error", err.Error()))
		return err
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		slog.Error("failed to create JetStream context", slog.String("error", err.Error()))
		return err
	}

	endpts := makeEndpoints(cfg, js)

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

	natsConn.Close()
	// shutdown ctx
	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown HTTP server", slog.String("error", err.Error()))
	}

	slog.Info("HTTP server gracefully stopped")

	return nil
}

func makeEndpoints(cfg config.Config, js jetstream.JetStream) endpoint.Endpoint {
	dbConn := db.InitDB(cfg)

	// init all repo
	userRepository := repository.NewUserRepository(dbConn)

	// nats publisher
	publisher := natstransport.NewPublisher(js, natstransport.JSONEncoder)

	return endpoint.Endpoint{
		User: makeUserEndpoints(userRepository, publisher, cfg),
	}
}

func makeUserEndpoints(userRepository *repository.UserRepository,
	publisher *natstransport.Publisher, cfg config.Config) endpoint.User {
	userSvc := service.NewUserService(userRepository, publisher)

	return endpoint.NewUserEndpoint(userSvc)
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
