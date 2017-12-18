package main

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var GithubID = os.Getenv("GITHUB_CLIENT_ID")
var GithubSecret = os.Getenv("GITHUB_CLIENT_SECRET")
var OpenerOrigin = "*"

func GetAuthorizeURL(state string) string {
	u := &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/login/oauth/authorize",
	}

	q := u.Query()
	q.Set("client_id", GithubID)
	q.Set("state", state)
	u.RawQuery = q.Encode()

	return u.String()
}

func GetAccessToken(code, state string) (string, error) {
	client := http.Client{
		Timeout:   time.Second * 20,
		Transport: http.DefaultTransport,
	}
	tr := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "/login/oauth/access_token",
		},
	}

	q := tr.URL.Query()
	q.Set("client_id", GithubID)
	q.Set("client_secret", GithubSecret)
	q.Set("code", code)
	q.Set("state", state)

	tr.URL.RawQuery = q.Encode()

	log.Printf("getting access token: %s", tr.URL)
	res, err := client.Do(tr)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.New("error with access token request: " + res.Status)
	}

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	response, _ := url.ParseQuery(string(rawbody))
	if response.Get("error") != "" {
		return "", errors.New(response.Get("error_description"))
	}

	return string(response.Get("access_token")), nil
}

func BuildRedirectPage(token string) []byte {
	t, err := template.New("access-token.html").Parse(getTemplateFile())
	if err != nil {
		return []byte("There was a problem parsing the file needed to show this page.")
	}

	o := bytes.Buffer{}
	err = t.Execute(&o, templateData{
		AccessToken: token,
		Origin:      OpenerOrigin,
	})
	if err != nil {
		return []byte("There was a problem building the file needed to show this page.")
	}

	a, _ := ioutil.ReadAll(&o)
	return a
}

func getTemplateFile() string {
	wd, _ := os.Getwd()
	f := filepath.Join(wd, "access-token.html")

	contents, err := ioutil.ReadFile(f)
	if err != nil {
		return "There was a problem reading the file needed to show this page."
	}

	return string(contents)
}

type templateData struct {
	AccessToken string
	Origin      string
}
