package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"views-servive/services"
)

func TestParseHandler(parseService services.ParseServiceImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proxies := parseService.Parse()
		err := json.NewEncoder(w).Encode(proxies)
		if err != nil {
			log.Println("marshaling error: ", err)
			w.Write([]byte("[]"))
			return
		}
	}
}
