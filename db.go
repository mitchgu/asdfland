package main

import (
	"time"
	// "log"

	"github.com/go-redis/redis"
)

type DB interface {
	DBType() string
	ReserveSlug(fp string, val string) bool
}

type RedisDB struct {
	Client *redis.Client
}

func (rdb *RedisDB) Init(addr, pass string, dbnum int) {
	rdb.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       dbnum,
	})
}

func (rdb *RedisDB) DBType() string {
	return "redis"
}

func fingerprintKey(fp string) string { return "fp:" + fp }
func reserveKey(slug string) string   { return "reserve:" + slug }
func slugKey(slug string) string      { return "slug:" + slug }

func (rdb *RedisDB) ReserveSlug(fp string, slug string) bool {
	fpKey := fingerprintKey(fp)
	rKey := reserveKey(slug)
	slugKey := slugKey(slug)

	// TODO: fix race conditions

	// check if key exists or is reserved
	isSlug := rdb.Client.Exists(slugKey).Val() > 0
	if isSlug {
		return false
	}
	isReserved := rdb.Client.Exists(rKey).Val() > 0

	// Get currently reserved slug
	reservedByMe, _ := rdb.Client.HGet(fpKey, "reserve").Result()
	if isReserved && reservedByMe != rKey {
		return false // It's already reserved but by someone else
	}
	// Clear previously reserved slug
	if reservedByMe != "" {
		rdb.Client.Del(reservedByMe) // Unreserve the previously reserved slug
	}

	// Set the new one as reserved
	rdb.Client.HSet(fpKey, "reserve", rKey)
	rdb.Client.Set(rKey, fp, 5*time.Minute)

	return true
}
