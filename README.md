# Github Proxy

This is a simple proxy to the GitHub API, with an OAuth endpoint. I'm sure tons of people have created things that do this, but I wanted to learn how to do it myself, especially in Go.

# HTTP API

* `/oauth/authorize`: Starts OAuth workflow, redirecting to the GitHub OAuth authorization page
* `/oauth/code`: Consumes an OAuth code from GitHub, then retrieves an access token and renders a page that attempts to send the token back to the opening window via a `postMessage` call.
* `/proxy`: a direct proxy to https://api.github.com. The `/proxy` part of the path is removed before sending the request upstream.

## Environment variables

* `GITHUB_CLIENT_ID`: The client ID for the GitHub app
* `GITHUB_CLIENT_SECRET`: The client secret for the GitHub app
