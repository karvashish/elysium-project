package handlers

import (
	"fmt"
	"net/http"
)

func GetPeerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetPeer!")
}
