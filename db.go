package main

import (
	"github.com/go-redis/redis"
)

type DB interface {
	DBType() string
}

type RedisDB struct {
	Client *redis.Client
}

func (rdb *RedisDB) Init(addr, pass string, dbnum int) {
	rdb.Client = redis.NewClient(&redis.Options{
		Addr: 		addr,
		Password:	pass,
		DB:			dbnum,
	})
}

func (rdb *RedisDB) DBType() string {
	return "redis"
}