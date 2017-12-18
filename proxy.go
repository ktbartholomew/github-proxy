package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var upstream *url.URL = &url.URL{
	Scheme: "https",
	Host:   "api.github.com",
}

var GithubProxy httputil.ReverseProxy = httputil.ReverseProxy{
	Director: func(r *http.Request) {
		r.Host = upstream.Host
		r.URL.Host = upstream.Host
		r.URL.Scheme = upstream.Scheme
		r.URL.Path = r.URL.Path[6:]

		if _, ok := r.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			r.Header.Set("User-Agent", "")
		}

		log.Printf("Proxying request: %s %s", r.Method, r.URL)
	},
}
