package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/config"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/dto"
	"github.com/ijalalfrz/event-driven-nats/gateway-service/internal/app/endpoint"
	httptransport "github.com/ijalalfrz/event-driven-nats/gateway-service/internal/pkg/transport/http"
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

		router.Route("/public", func(router chi.Router) {
			router.Route("/listings", func(router chi.Router) {
				router.Post("/", httptransport.MakeHandlerFunc(
					endpts.PublicListing.Create,
					httptransport.DecodeRequest[dto.CreateListingRequest],
					httptransport.ResponseWithBody,
				))

				router.Get("/", httptransport.MakeHandlerFunc(
					endpts.PublicListing.GetAll,
					httptransport.DecodeRequest[dto.GetAllListingsRequest],
					httptransport.ResponseWithBody,
				))
			})

			router.Route("/users", func(router chi.Router) {
				router.Post("/", httptransport.MakeHandlerFunc(
					endpts.PublicUser.Create,
					httptransport.DecodeRequest[dto.CreateUserRequest],
					httptransport.ResponseWithBody,
				))
			})
		})
	})

	return router
}
