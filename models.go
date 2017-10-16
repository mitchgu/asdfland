package main

import (
// "time"
)

type Dest struct {
	Dest        string
    Description string
    Password    string
    EnableAnalytics bool
}

type SlugReserveReq struct {
    Type       string
    Length     int
    Wordlist   string
    CustomSlug string
}

type SlugDestCreateReq struct {
    Slug string
    Dest string
    Expire int
    Description string
    Password string
    EnableAnalytics bool
}

func DestCreate(d Dest) bool {
	// dest_hash := map[string]interface{}{
	// 	"url": d.Url,
	// 	"title": d.Title,
	// 	"creation_ip": d.CreationIP,
	// 	"created_at": time.Now().UTC().Format(time.RFC3339),
	// }
	// dest_uuid := GetUUID(24)
	// db.HMSet("dest:" + dest_uuid, dest_hash)
	return true
}

func (dest *Dest) ToMap() *map[string]interface{} {
    destMap := map[string]interface{}{
        "Dest": dest.Dest,
        "Description": dest.Description,
        "Password": dest.Password,
        "EnableAnalytics": dest.EnableAnalytics,
    }
    return &destMap
}