package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Response struct {
	ReqURL   string              `json:"reqUrl"`
	Payload  map[string][]string `json:"payload"`
	Headers  map[string][]string `json:"headers"`
	ClientIP string              `json:"clientIP"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		var payload map[string][]string
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		payload = r.Form

		headers := make(map[string][]string)
		for key, values := range r.Header {
			headers[key] = values
		}

		clientIP := r.RemoteAddr

		response := Response{
			ReqURL:   r.RequestURI,
			Payload:  payload,
			Headers:  headers,
			ClientIP: clientIP,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Error encoding response", zap.Error(err))
		}
	})

	port := "8080"
	logger.Info("Starting server", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}
