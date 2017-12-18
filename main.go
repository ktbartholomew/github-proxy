package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	config()
	http.HandleFunc("/oauth/authorize", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Location", GetAuthorizeURL(req.URL.Query().Get("state")))
		res.WriteHeader(302)
	})

	http.HandleFunc("/oauth/code", func(res http.ResponseWriter, req *http.Request) {
		code := req.URL.Query().Get("code")
		state := req.URL.Query().Get("state")

		token, err := GetAccessToken(code, state)
		if err != nil {
			log.Println(err.Error())
			res.WriteHeader(500)
			res.Write([]byte(err.Error()))
			return
		}

		res.WriteHeader(200)
		res.Write(BuildRedirectPage(token))
	})

	http.HandleFunc("/proxy/", func(res http.ResponseWriter, req *http.Request) {
		GithubProxy.ServeHTTP(res, req)
	})

	listenAddr := ":8080"
	if os.Getenv("LISTEN") != "" {
		listenAddr = ":" + os.Getenv("LISTEN")
	}

	log.Printf("Listening on %s", listenAddr)
	go log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func config() {
	if os.Getenv("GITHUB_CLIENT_ID") == "" {
		log.Println("GITHUB_CLIENT_ID not set")
		os.Exit(1)
	}

	if os.Getenv("GITHUB_CLIENT_SECRET") == "" {
		log.Println("GITHUB_CLIENT_SECRET not set")
		os.Exit(1)
	}
}
