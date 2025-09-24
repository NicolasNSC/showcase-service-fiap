package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func SetupRoutes(router *chi.Mux, saleHandler *SaleHandler) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/listings", saleHandler.CreateListing)
	router.Get("/sales/available", saleHandler.ListAvailable)
	router.Get("/sales/sold", saleHandler.ListSold)
}
