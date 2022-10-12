package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const proxyURL = "https://gcsim.app"

func main() {
	remote, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL)
			r.Host = remote.Host
			p.ServeHTTP(w, r)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)
	http.HandleFunc("/", handler(proxy))
	err = http.ListenAndServe(":3030", nil)
	if err != nil {
		panic(err)
	}
}
