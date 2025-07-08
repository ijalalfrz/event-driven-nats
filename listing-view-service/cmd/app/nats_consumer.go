package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	gokitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/config"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/endpoint"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/repository"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/service"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/pkg/db"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/pkg/logger"
	natstransport "github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/pkg/transport/nats"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/spf13/cobra"
)

var (
	userCreatedSubject    = "user.created"
	listingCreatedSubject = "listing.created"
	streamName            = "listing_view_event"
)

var natsConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consume events from NATS",
	Run: func(_ *cobra.Command, _ []string) {
		slog.Debug("command line flags", slog.String("config_path", cfgFilePath))
		cfg := config.MustInitConfig(cfgFilePath)

		logger.InitStructuredLogger(cfg.LogLevel)

		runNATS(cfg)
	},
}

func runNATS(cfg config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nc, err := nats.Connect(cfg.NATS.URL,
		nats.MaxReconnects(cfg.NATS.MaxReconnects),
		nats.ReconnectWait(cfg.NATS.ReconnectWait))
	if err != nil {
		slog.Error("failed to connect to NATS", "error", err)
		return
	}

	js, err := jetstream.New(nc)
	if err != nil {
		slog.Error("failed to create NATS JetStream", "error", err)
		return
	}

	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name: streamName,
		Subjects: []string{
			userCreatedSubject,
			listingCreatedSubject,
		},
	})
	if err != nil {
		slog.Error("failed to create stream", "error", err)
		return
	}

	middlewares := []gokitendpoint.Middleware{
		natstransport.AutoAckMiddleware(),
	}
	endpoints := makeNatsEndpoints(cfg)

	userCreatedConsumer, err := natstransport.NewSubscriber(
		ctx,
		stream,
		userCreatedSubject,
		endpoints.User.OnCreated,
		natstransport.NewDecoder[dto.UserCreated](),
		middlewares,
	)
	if err != nil {
		slog.Error("failed to create user created consumer", "error", err)
		return
	}

	listingCreatedConsumer, err := natstransport.NewSubscriber(
		ctx,
		stream,
		listingCreatedSubject,
		endpoints.Listing.OnCreated,
		natstransport.NewDecoder[dto.ListingCreated](),
		middlewares,
	)
	if err != nil {
		slog.Error("failed to create listing created consumer", "error", err)
		return
	}

	userCreatedConsumer.Start(ctx)
	listingCreatedConsumer.Start(ctx)

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigChannel:
		cancel()
		slog.InfoContext(ctx, "received OS signal. Exiting...", slog.String("signal", sig.String()))
	case <-ctx.Done():
		slog.Info("context done. Exiting...", "error", ctx.Err())
	}

	userCreatedConsumer.Stop()
	listingCreatedConsumer.Stop()
	nc.Close()

	slog.Info("nats consumer stopped")
}

func makeNatsEndpoints(cfg config.Config) endpoint.Endpoint {
	dbConn := db.InitDB(cfg)

	// init all repo
	userRepo := repository.NewUserRepository(dbConn)
	listingRepo := repository.NewListingRepository(dbConn)

	return endpoint.Endpoint{
		User:    makeUserEndpoint(userRepo),
		Listing: makeListingEndpoint(listingRepo, userRepo),
	}
}

func makeUserEndpoint(userRepo *repository.UserRepository) endpoint.User {
	userSvc := service.NewUserService(userRepo)

	return endpoint.NewUserEndpoint(userSvc)
}

func makeListingEndpoint(listingRepo *repository.ListingRepository, userRepo *repository.UserRepository) endpoint.Listing {
	listingSvc := service.NewListingService(listingRepo, userRepo)
	listingViewSvc := service.NewListingViewService(listingRepo)

	return endpoint.NewListingEndpoint(listingViewSvc, listingSvc)
}
