
package main

import (

	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"fmt"
)

//utilities

// Get env var url or if none use default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		
		return value
	}
	return fallback
}

//Getters

// Get the port to listen on
func getListenAddress() string {
	port := getEnv("PORT", "8080")
	return ":" + port
}


//Get URL to evidence 1
func getProxyUrl() string {

	evidence1 := os.Getenv("EVIDENCE_1")
	fmt.Println(evidence1)
	return evidence1
}

//Logging


// Log the env variables required for a reverse proxy
func logSetup() {

	log.Printf("Server will run on: %s\n", getListenAddress())
	log.Printf("Connecting with evidence 1: %s\n", getProxyUrl())
}

//Reverse proxy

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	proxy.ServeHTTP(res, req)
}


// Receive and handle request, sending it to URL of evidence 1
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	url := getProxyUrl()
	serveReverseProxy(url, res, req)
}

//Main

func main() {
	// Log setup values
	logSetup()

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
	if err := http.ListenAndServe(getListenAddress(), nil); err != nil {
		panic(err)
	}
}