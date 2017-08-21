package main

import (
	"math"
	"math/rand"
	"encoding/base64"
)


func GetUUID(l int) string {
    token := make([]byte, int(math.Ceil(float64(l)*0.75)))
    rand.Read(token)
    uuid := base64.RawURLEncoding.EncodeToString(token)
    return uuid
}