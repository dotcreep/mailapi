package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/dotcreep/mailapi/docs"
	"github.com/dotcreep/mailapi/internal/service/api"
	"github.com/dotcreep/mailapi/internal/service/api/send_template"
	"github.com/dotcreep/mailapi/internal/utils"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// @title						Postmail API
// @version					1.0
// @description				Documentation for Postmail
// @license.name				MIT
// @license.url				https://opensource.org/licenses/MIT
// @BasePath					/
// @SecurityDefinitions.apikey	X-Auth-Key
// @name						X-Auth-Key
// @in							header
// @description				Input your token authorized
func main() {
	Json := utils.Json{}
	cfg, err := utils.OpenYAML()
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(cfg.Config.Port)
	if err != nil {
		panic(err)
	}

	utils.InitStorage()
	r := chi.NewRouter()
	m := api.Midleware
	cors := cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Auth-Key"},
		AllowCredentials: true,
	})
	r.Use(cors)
	r.Route("/send", func(r chi.Router) {
		r.Use(m)
		r.Post("/", api.SendEmailAPI)
		r.Post("/template/welcome", send_template.Welcome)
		r.Post("/template/registered", send_template.Registered)
		r.Post("/template/credentials", send_template.Credentials)
		r.Post("/template/rejected", send_template.Rejected)
		r.Post("/template/approved", send_template.Approved)
		r.Post("/template/invoice-paid", send_template.InvoicePaid)
		r.Post("/template/invoice", send_template.Invoice)
	})
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		Json.NewResponse(false, w, nil, "404 not found", http.StatusNotFound, nil)
	})
	fmt.Printf("Server running on port %d\n", port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
