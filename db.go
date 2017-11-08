package main

import (
	"time"
	// "log"
	"fmt"

	"github.com/go-redis/redis"
)

type DB interface {
	DBType() string
	ReserveSlug(fp string, val string) bool
	SlugReserved(fp string, val string) bool
	DestCreate(dest *Dest) (string, bool)
	SlugCreate(slug, dest string, expire int, fp string) bool
	SlugFollow(slug string) (string, error)
	UserCreate(username, digest string) error
	UserGetDigest(username string) (string, error)
	SessionCreate(username, token string) string
	SessionLookup(token string) (bool, string, bool)
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
func userKey(username string) string  { return "user:" + username }
func sessionKey(token string) string  { return "session:" + token }

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
	rdb.Client.Set(rKey, fp, 15*time.Minute)

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
	// add to userdests list
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
	// add to destslugs list
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

func (rdb *RedisDB) UserCreate(username, digest string) error {
	uKey := userKey(username)
	userExists := rdb.Client.Exists(uKey).Val() > 0
	if userExists {
		return fmt.Errorf("User already exists")
	}
	rdb.Client.HSet(uKey, "digest", digest)
	return nil
}

func (rdb *RedisDB) UserGetDigest(username string) (string, error) {
	uKey := userKey(username)
	return rdb.Client.HGet(uKey, "digest").Result()
}

func (rdb *RedisDB) SessionCreate(username, token string) string {
	sKey := sessionKey(token)
	rdb.Client.Set(sKey, username, 30*24*time.Hour)
	return token
}

func (rdb *RedisDB) SessionLookup(token string) (bool, string, bool) {
	// Returns if the session token is valid, what the username is, and if the corresponding user exists
	username, err := rdb.Client.Get(sessionKey(token)).Result()
	if err != nil {
		return false, "", false
	} else {
		return true, username, rdb.Client.Exists(userKey(username)).Val() > 0
	}
}
