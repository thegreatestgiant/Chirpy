package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const port = "8000"
	const fileServerPath = "."

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()
	api := chi.NewRouter()
	admin := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(fileServerPath))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	admin.Get("/metrics", apiCfg.getMetrics)

	api.Get("/healthz", handlerReadiness)
	api.Get("/reset", apiCfg.resetMetrics)

	api.Post("/validate_chirp", validateLength)

	r.Mount("/api", api)
	r.Mount("/admin", admin)

	corsMux := middlewareCors(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Listening to FileServer in %s on port %v\n", fileServerPath, port)
	log.Fatal(server.ListenAndServe())
}
