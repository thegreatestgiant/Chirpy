package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8000"
	const fileServerPath = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(fileServerPath)))
	corsMux := middlewareCors(mux)

	server := http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	fmt.Printf("Listening to FileServer in %s on port %v\n", fileServerPath, port)
	log.Fatal(server.ListenAndServe())
}
