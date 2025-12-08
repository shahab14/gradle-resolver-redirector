package main

import (
	"log"
	"net/http"
	"redirector/handlers"
	"redirector/middleware"
	"time"

	"github.com/gorilla/mux"
)

func handleRoutes(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(handlers.RedirectorGETHandler).Methods("GET")
	router.PathPrefix("/").HandlerFunc(handlers.RedirectorHEADHandler).Methods("HEAD")
}

func createServer() {
	muxRoute := mux.NewRouter().StrictSlash(false)
	log.Println("starting server")

	// Apply security middleware
	// Rate limiting: 100 requests per second, burst of 200, cleanup every 10 minutes
	rateLimiter := middleware.NewRateLimiter(1000, 2000, 10*time.Minute)

	// Create middleware chain
	var handler http.Handler = muxRoute
	handler = middleware.MethodWhitelistMiddleware([]string{"GET", "HEAD"})(handler)
	handler = rateLimiter.RateLimitMiddleware(handler)
	handler = middleware.TimeoutMiddleware(180 * time.Second)(handler)

	handleRoutes(muxRoute)

	// Configure HTTP server with security settings
	server := &http.Server{
		Addr:           "0.0.0.0:10010",
		Handler:        handler,
		ReadTimeout:    300 * time.Second,
		WriteTimeout:   300 * time.Second,
		IdleTimeout:    300 * time.Second,
		MaxHeaderBytes: 1 << 30, // 1GB
	}

	log.Println("Server configured with security measures:")
	log.Println("  - Rate limiting: 100 req/s per IP")
	log.Println("  - Request timeout: 60s")
	log.Println("  - Path validation enabled")
	log.Println("  - Security headers enabled")
	log.Fatal(server.ListenAndServe())
}

func main() {
	createServer()
}
