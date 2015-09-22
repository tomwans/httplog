package main

import (
	"fmt"
	log "github.com/tomwans/httplog"
	"net/http"
	"time"
)

func main() {
	l := log.New()

	go http.ListenAndServe(":21200", l)

	l.Println("initial")

	i := 0
	for range time.Tick(10 * time.Second) {
		l.Println(fmt.Sprintf("test%d", i))
		i++
	}
}
