package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("handlers.BaseHandler -> called")
	fmt.Fprintf(w, "Hello, World!")
}
