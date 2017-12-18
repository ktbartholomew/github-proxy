# Github Proxy

This is a simple proxy to the GitHub API, with an OAuth endpoint. I'm sure tons of people have created things that do this, but I wanted to learn how to do it myself, especially in Go.

# HTTP API

* `/authenticate`: Handles OAuth workflow, receives a code from GitHub and uses secret client config to turn that into a usable access token
* `/proxy`: a direct proxy to https://api.github.com. The `/proxy` part of the path is removed before sending the request upstream.