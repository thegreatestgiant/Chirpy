package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/thegreatestgiant/go-server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	const port = "8000"
	const fileServerPath = "."
	const dbPath = "database.json"

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
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

	api.Post("/chirps", apiCfg.handlerCreateChirp)
	api.Post("/users", apiCfg.handlerCreateUser)
	api.Post("/login", apiCfg.handlerLogin)
	api.Post("/refresh", apiCfg.handlerRefresh)
	api.Post("/revoke", apiCfg.handlerRevoke)

	api.Put("/users", apiCfg.handlerPutUser)

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
