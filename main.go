// Reqlog is a simple http proxy that prints the contents
// of requests and responses.
package main

import (
	"net/url"
	"os"
	"net/http"
	"net/http/httputil"
	"log"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 3 {
		log.Fatal("Usage: reqlog addr url")
	}
	u, err := url.Parse(os.Args[2])
	if err != nil {
		panic(err)
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	rp.Transport = &Transport{http.DefaultTransport}
	http.Handle("/", rp)
	panic(http.ListenAndServeTLS(os.Args[1], "localhost.crt.pem", "localhost.key.pem", nil))
}

type Transport struct {
	next http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	logReq(req)
	resp, err := t.next.RoundTrip(req)
	logResp(resp)
	return resp, err
}

func logReq(r *http.Request) {
	if dump, err := httputil.DumpRequestOut(r, true); err != nil {
		log.Println("reqlog: error dumping request:", err)
	} else {
		log.Println(string(dump))
	}
}

func logResp(r *http.Response) {
	if dump, err := httputil.DumpResponse(r, true); err != nil {
		log.Println("reqlog: error dumping response:", err)
	} else {
		log.Println(string(dump))
	}
}
