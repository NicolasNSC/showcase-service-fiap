package handler

import (
	_ "github.com/NicolasNSC/showcase-service-fiap/docs"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/NicolasNSC/showcase-service-fiap/docs"
)

func SetupRoutes(router *chi.Mux, saleHandler *SaleHandler) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Post("/listings", saleHandler.CreateListing)
	router.Put("/listings/vehicle/{vehicle_id}", saleHandler.UpdateListing)
	router.Post("/webhooks/payments", saleHandler.HandlePaymentWebhook)

	router.Route("/sales/{id}", func(r chi.Router) {
		r.Post("/purchase", saleHandler.Purchase)
	})

	router.Get("/sales/available", saleHandler.ListAvailable)
	router.Get("/sales/sold", saleHandler.ListSold)
}
