package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"views-servive/services"
)

func ViewsHandler(viewsService services.ViewsServiceImp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var linkForBoost string
		err := json.NewDecoder(r.Body).Decode(&linkForBoost)
		if err != nil {
			log.Println("unmarshalling error: ", err.Error())
			w.Write([]byte("[]"))
			return
		}

		res := viewsService.Boost(linkForBoost)
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println("marshaling error: ", err)
			w.Write([]byte("[]"))
			return
		}
	}
}
