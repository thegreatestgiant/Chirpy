package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thegreatestgiant/go-server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const port = "8000"
	const fileServerPath = "."
	const dbPath = "database.json"

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
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
	api.Get("/chirps", apiCfg.handlerGetChirps)
	api.Get("/reset", apiCfg.handlerResetDB)
	api.Get("/chirps/{chirpID}", apiCfg.handlerGetChirpsByID)
	api.Get("/users", apiCfg.handlerGetusers)
	api.Get("/users/{UserID}", apiCfg.handlerGetUserByID)

	api.Post("/chirps", apiCfg.handlerCreateChirp)
	api.Post("/users", apiCfg.handlerCreateUser)

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
