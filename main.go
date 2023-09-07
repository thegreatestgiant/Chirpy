package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8000"

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	fmt.Printf("Listening on port %v\n", port)
	log.Fatal(server.ListenAndServe())
}
