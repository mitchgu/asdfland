package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	a := App{}
	if c.GetString("db_kind") == "redis" {
		a.InitRedis(
			c.GetString("redis_addr"),
			c.GetString("redis_pass"),
			c.GetInt("redis_dbnum"),
			c.GetString("frontend_dir"))
	} else {
		log.Fatalf("Database type not supported: %s", c.GetString("db_kind"))
	}
	a.Run(c.GetString("server_addr"))
}
