package main

import (
	"log"
	"time"
    "math/rand"

	"net/http"
)

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
	r := GetRouter()

    srv := &http.Server{
        Handler:      r,
        Addr:         c.GetString("server_addr"),
        // Good practice: enforce timeouts for servers you create!
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    log.Println("Starting server")
    log.Fatal(srv.ListenAndServe())
}