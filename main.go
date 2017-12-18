package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var githubID = os.Getenv("GITHUB_CLIENT_ID")
var githubSecret = os.Getenv("GITHUB_CLIENT_SECRET")

func main() {
	http.HandleFunc("/authenticate", func(res http.ResponseWriter, req *http.Request) {
		log.Printf("%+v", req)
		res.Write([]byte("authenticate"))
	})

	http.HandleFunc("/proxy/", func(res http.ResponseWriter, req *http.Request) {
		upstream := &url.URL{}
		upstream, err := url.Parse("https://api.github.com")
		if err != nil {
			log.Println(err.Error())
			return
		}

		rp := httputil.ReverseProxy{
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

		rp.ServeHTTP(res, req)
	})

	listenAddr := ":8080"
	if os.Getenv("LISTEN") != "" {
		listenAddr = ":" + os.Getenv("LISTEN")
	}

	log.Printf("Listening on %s", listenAddr)
	go log.Fatal(http.ListenAndServe(listenAddr, nil))
}
