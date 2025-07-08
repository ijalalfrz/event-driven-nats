package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/config"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/app/endpoint"
	httptransport "github.com/ijalalfrz/event-driven-nats/listing-view-service/internal/pkg/transport/http"
)

// MakeHTTPRouter builds the HTTP router with all the service endpoints.
func MakeHTTPRouter(
	endpts endpoint.Endpoint,
	cfg config.Config,
) *chi.Mux {
	// Initialize Router
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	router.Route("/", func(router chi.Router) {
		router.Use(
			httptransport.LoggingMiddleware(slog.Default()),
			httptransport.CORSMiddleware(cfg.HTTP.AllowedOrigin),
			httptransport.Recoverer(slog.Default()),
			render.SetContentType(render.ContentTypeJSON),
		)

		router.Route("/listings", func(router chi.Router) {
			router.Get("/", httptransport.MakeHandlerFunc(
				endpts.Listing.GetAll,
				httptransport.DecodeRequest[dto.GetAllListingsRequest],
				httptransport.ResponseWithBody,
			))
		})
	})

	return router
}
