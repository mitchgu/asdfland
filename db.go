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
	BuildDestIndex(username string) []*DestListing
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

func fingerprintKey(fp string) string     { return "fp:" + fp }
func reserveKey(slug string) string       { return "reserve:" + slug }
func slugKey(slug string) string          { return "slug:" + slug }
func destKey(destUUID string) string      { return "dest:" + destUUID }
func userKey(username string) string      { return "user:" + username }
func sessionKey(token string) string      { return "session:" + token }
func userdestsKey(username string) string { return "userdests:" + username }
func destslugsKey(dest string) string     { return "destslugs:" + dest }

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
	// add to userdests sorted set
	owner := dest.Owner
	rdb.Client.ZAdd(userdestsKey(owner), redis.Z{Score: float64(time.Now().Unix()), Member: destUUID})
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
	rdb.Client.LPush(destslugsKey(destUUID), slug)
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

type DestListing struct {
	UUID            string
	Dest            string
	Description     string
	EnableAnalytics bool
	CreatedAt       string
	Slugs           []SlugListing
}

type SlugListing struct {
	Slug    string
	Expires string
}

func (rdb *RedisDB) BuildDestIndex(username string) []*DestListing {
	var destIndex []*DestListing
	userDestsZ := rdb.Client.ZRevRangeWithScores(userdestsKey(username), 0, -1).Val()
	for _, Z := range userDestsZ {
		createdAt := time.Unix(int64(Z.Score), 0).Format(time.RFC3339)
		ud := Z.Member.(string)
		destMap := rdb.Client.HGetAll(destKey(ud)).Val()
		destSlugs := rdb.Client.LRange(destslugsKey(ud), 0, -1).Val()
		var slugIndex []SlugListing
		for _, ds := range destSlugs {
			ttl := rdb.Client.TTL(slugKey(ds)).Val()
			var expires string
			if ttl > 0 {
				expires = time.Now().Add(ttl).Format(time.RFC3339)
			} else {
				expires = ""
			}
			slugIndex = append(slugIndex, SlugListing{
				Slug:    ds,
				Expires: expires,
			})
		}
		dl := DestListing{
			UUID:            ud,
			Dest:            destMap["Dest"],
			Description:     destMap["Description"],
			EnableAnalytics: destMap["EnableAnalytics"] == "1",
			CreatedAt:       createdAt,
			Slugs:           slugIndex,
		}
		destIndex = append(destIndex, &dl)
	}
	return destIndex
}
