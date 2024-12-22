package handlers

import (
	"elysium-backend/config"
	"fmt"
	"log"
	"net/http"
)

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("handlers.BaseHandler -> called")
	}
	fmt.Fprintf(w, "Hello, World!")
}
