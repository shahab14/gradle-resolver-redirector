package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	req, err := http.NewRequest(method, "https://dl.google.com/dl/android/maven2"+r.URL.Path, nil)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return -1
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", r.Header["User-Agent"][0])
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return -1
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return http.StatusNotFound
	} else {
		reqByte, _ := io.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Write(reqByte)
		return resp.StatusCode
	}
}

func SendRequestToRepo1Maven(w http.ResponseWriter, r *http.Request, method string) int {
	req, err := http.NewRequest(method, "https://repo1.maven.org/maven2"+r.URL.Path, nil)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return -1
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", r.Header["User-Agent"][0])
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return -1
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return http.StatusNotFound
	} else {
		reqByte, _ := io.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Write(reqByte)
		return resp.StatusCode
	}
}

func SendRequestToJcenterMaven(w http.ResponseWriter, r *http.Request, method string) int {
	req, err := http.NewRequest(method, "https://jcenter.bintray.com"+r.URL.Path, nil)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return -1
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", r.Header["User-Agent"][0])
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		_, _ = fmt.Fprint(w, err.Error())
		return -1
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return http.StatusNotFound
	} else {
		reqByte, _ := io.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Write(reqByte)
		return resp.StatusCode
	}
}
