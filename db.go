package main

import (
	"time"
	// "log"

	"github.com/go-redis/redis"
)

type DB interface {
	DBType() string
	ReserveSlug(fp string, val string) bool
	SlugReserved(fp string, val string) bool
	DestCreate(dest *Dest) (string, bool)
	SlugCreate(slug, dest string, expire int, fp string) bool
	SlugFollow(slug string) (string, error)
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
func destKey(destUUID string) string  { return "dest:" + destUUID }

func (rdb *RedisDB) ReserveSlug(fp string, slug string) bool {
	fpKey := fingerprintKey(fp)
	rKey := reserveKey(slug)
	sKey := slugKey(slug)

	// TODO: fix race conditions

	// check if key exists or is reserved
	slugExists := rdb.Client.Exists(sKey).Val() > 0
	if slugExists {
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

func (rdb *RedisDB) SlugReserved(fp, slug string) bool {
	fpKey := fingerprintKey(fp)
	rKey := reserveKey(slug)
	sKey := slugKey(slug)

	slugExists := rdb.Client.Exists(sKey).Val() > 0
	slugIsReserved := rdb.Client.Exists(rKey).Val() > 0
	reservedByMe, _ := rdb.Client.HGet(fpKey, "reserve").Result()
	return !slugExists && slugIsReserved && reservedByMe == rKey
}

func (rdb *RedisDB) DestCreate(dest *Dest) (string, bool) {
	destUUID := GetRandString(24)
	destKey := destKey(destUUID)
	err := rdb.Client.HMSet(destKey, *dest.ToMap()).Err()
	return destUUID, err == nil
}

func (rdb *RedisDB) SlugCreate(slug, destUUID string, expire int, fp string) bool {
	fpKey := fingerprintKey(fp)
	dKey := destKey(destUUID)
	sKey := slugKey(slug)
	rKey := reserveKey(slug)
	rdb.Client.Set(sKey, dKey, time.Duration(expire)*time.Minute).Err()
	rdb.Client.HDel(fpKey, "reserve").Err()
	rdb.Client.Del(rKey)
	return true
}

func (rdb *RedisDB) SlugFollow(slug string) (string, error) {
	sKey := slugKey(slug)
	dKey, errGet := rdb.Client.Get(sKey).Result()
	if errGet != nil {
		return "", errGet
	}
	url, errHGet := rdb.Client.HGet(dKey, "Dest").Result()
	if errHGet != nil {
		return "", errHGet
	}
	return url, errHGet
}
