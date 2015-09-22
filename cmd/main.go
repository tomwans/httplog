package main

import (
	"fmt"
	log "github.com/tomwans/httplog"
	"net/http"
	"time"
)

func main() {
	l := log.New("test: ")
	s := &http.Server{Handler: l, Addr: ":21200"}
	// http2.ConfigureServer(s, nil)
	go s.ListenAndServe()

	l.Println("initial")

	i := 0
	for range time.Tick(2 * time.Second) {
		l.Println(fmt.Sprintf("test%d", i))
		i++
	}
}
