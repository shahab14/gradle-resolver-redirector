package handlers

import (
	"io"
	"log"
	"net/http"
	"time"
)

const (
	// MaxResponseSize limits response body size to prevent memory exhaustion (50MB)
	MaxResponseSize = 50 * 1024 * 1024
	// HTTPClientTimeout sets timeout for outgoing HTTP requests
	HTTPClientTimeout = 30 * time.Second
)

func RedirectorGETHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	statusCode := SendRequestToGoogleMaven(w, r, "GET")
	log.Println("Google Status Code: ", statusCode)
	if statusCode != http.StatusNotFound {
		return
	}
	statusCode = SendRequestToJcenterMaven(w, r, "GET")
	log.Println("JCenter Status Code: ", statusCode)
	if statusCode != http.StatusNotFound {
		return
	}
	statusCode = SendRequestToRepo1Maven(w, r, "GET")
	log.Println("Repo1 Status Code: ", statusCode)
	if statusCode != http.StatusNotFound {
		return
	}
}

func RedirectorHEADHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	log.Println(r.Method)

	statusCode := SendRequestToGoogleMaven(w, r, "HEAD")
	log.Println("Google Status Code: ", statusCode)
	if statusCode != http.StatusNotFound {
		return
	}
	statusCode = SendRequestToJcenterMaven(w, r, "HEAD")
	log.Println("JCenter Status Code: ", statusCode)
	if statusCode != http.StatusNotFound {
		return
	}
	statusCode = SendRequestToRepo1Maven(w, r, "HEAD")
	log.Println("Repo1 Status Code: ", statusCode)
	if statusCode != http.StatusNotFound {
		return
	}
}

func SendRequestToGoogleMaven(w http.ResponseWriter, r *http.Request, method string) int {
	return sendRequestToMaven(w, r, method, "https://dl.google.com/dl/android/maven2")
}

func SendRequestToRepo1Maven(w http.ResponseWriter, r *http.Request, method string) int {
	return sendRequestToMaven(w, r, method, "https://repo1.maven.org/maven2")
}

func SendRequestToJcenterMaven(w http.ResponseWriter, r *http.Request, method string) int {
	return sendRequestToMaven(w, r, method, "https://jcenter.bintray.com")
}

// sendRequestToMaven is a unified function to send requests to Maven repositories with security measures
func sendRequestToMaven(w http.ResponseWriter, r *http.Request, method string, baseURL string) int {
	req, err := http.NewRequest(method, baseURL+r.URL.Path, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return -1
	}

	// Safely get User-Agent header
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Gradle-Resolver-Redirector/1.0"
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "text/plain")

	// Create HTTP client with timeout and connection limits
	client := &http.Client{
		Timeout: HTTPClientTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to %s: %v", baseURL, err)
		// Don't expose internal errors to client
		return -1
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return http.StatusNotFound
	}

	// Limit response body size to prevent memory exhaustion
	limitedReader := io.LimitReader(resp.Body, MaxResponseSize)
	reqByte, err := io.ReadAll(limitedReader)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return -1
	}

	// Check if response was truncated
	if len(reqByte) >= MaxResponseSize {
		log.Printf("Response truncated for path: %s", r.URL.Path)
		http.Error(w, "Response too large", http.StatusInternalServerError)
		return -1
	}

	// Copy response headers (filter sensitive headers)
	for key, values := range resp.Header {
		// Skip headers that shouldn't be forwarded
		lowerKey := http.CanonicalHeaderKey(key)
		if lowerKey == "Connection" || lowerKey == "Keep-Alive" || lowerKey == "Transfer-Encoding" {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := w.Write(reqByte); err != nil {
		log.Printf("Error writing response: %v", err)
		return -1
	}

	return resp.StatusCode
}
