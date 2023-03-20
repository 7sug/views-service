package handlers

import (
	"io"
	"log"
	"net/http"
)

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, "im alive"); err != nil {
			log.Println("smth was wrong")
		}
	}
}
