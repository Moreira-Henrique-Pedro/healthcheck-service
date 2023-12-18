package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "OK"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	host := flag.String("host", "localhost", "Server Host")
	port := flag.Int("port", 3000, "Server Port")
	ssl := flag.Bool("ssl", false, "Use SSL (HTTPS)")

	flag.Parse()

	address := fmt.Sprintf("%s:%d", *host, *port)

	http.Handle("/", loggingMiddleware(http.HandlerFunc(healthCheckHandler)))

	if *ssl {
		fmt.Printf("Server running at https://%s\n", address)
		err := http.ListenAndServeTLS(address, "cert.pem", "key.pem", nil)
		if err != nil {
			fmt.Println("Server error:", err)
		}
	} else {
		fmt.Printf("Server running at http://%s\n", address)
		err := http.ListenAndServe(address, nil)
		if err != nil {
			fmt.Println("Server error:", err)
		}
	}
}
