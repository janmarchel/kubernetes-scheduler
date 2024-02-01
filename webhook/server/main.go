package server

import (
	"encoding/json"
	"net/http"
)

func HandleMutate(w http.ResponseWriter, r *http.Request) {
	// Przykładowa logika mutacji
	// Tutaj można modyfikować obiekty podów zgodnie z logiką szeregowania
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "success"}
	json.NewEncoder(w).Encode(response)
}
