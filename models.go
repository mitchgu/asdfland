package main

import (
// "time"
)

type Slug struct {
	Name string
}

type Dest struct {
	Url        string
	Title      string
	CreationIP string
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
