package main

import (
	"fmt"
	"log"
	"mailapi/api"
	"mailapi/env"
	"mailapi/utils"
	"net/http"
	"os"
	"strconv"
)

func main() {
	err := env.Load(".env")
	if err != nil {
		log.Println(err)
	}
	port, err := strconv.Atoi(os.Getenv("PORT_SERVER"))
	secret := os.Getenv("SECRET_ACCESS")
	if secret == "" {
		log.Println("Failed to get secret from .env file")
	}
	if err != nil {
		log.Println(err)
	}
	if port == 0 {
		log.Println("Failed to get port from .env file")
	}
	tmpPath := os.Getenv("STORAGE")
	if tmpPath == "" {
		log.Println("Failed to get tmp path from .env file")
	}
	mailAdmin := os.Getenv("MAIL_ADMIN")
	if mailAdmin == "" {
		log.Println("Failed to get mail admin from .env file")
	}
	homePage := os.Getenv("HOME_PAGE")
	if homePage == "" {
		log.Println("Failed to get home page from .env file")
	}

	utils.InitStorage()
	mux := http.NewServeMux()
	m := api.Midleware
	mux.HandleFunc("POST /send", m(api.SendEmailAPI))
	mux.HandleFunc("POST /send/template/{name}", m(api.SendEmailTemplate))
	mux.HandleFunc("/", m(api.NotFoundHandler))
	fmt.Printf("Server running on port %d\n", port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
