package main

import (
	"github.com/go-redis/redis"
)

var db = redis.NewClient(&redis.Options{
	Addr: 		c.GetString("redis_addr"),
	Password:	c.GetString("redis_pass"),
	DB:			c.GetInt("redis_db"),
})