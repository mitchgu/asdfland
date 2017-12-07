package main

import (
	"log"
	"math/rand"
	"time"
)

var version = "master"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	log.Printf("Starting asdfland version %s", version)

	a := App{}
	if c.GetString("db_kind") == "redis" {
		a.InitRedis(
			c.GetString("redis_addr"),
			c.GetString("redis_pass"),
			c.GetInt("redis_dbnum"))
	} else {
		log.Fatalf("Database type not supported: %s", c.GetString("db_kind"))
	}
	a.InitRouter()

	log.Printf("Serving at localhost:" + c.GetString("port"))
	a.Run(c.GetString("port"))
}
