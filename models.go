package main

import (
// "time"
)

type Dest struct {
	Dest            string
	Description     string
	Password        string
	Owner           string
	EnableAnalytics bool
}

type SlugReserveReq struct {
	Type       string
	Length     int
	Wordlist   string
	CustomSlug string
}

type SlugDestCreateReq struct {
	Slug            string
	Dest            string
	Expire          int
	Description     string
	Password        string
	EnableAnalytics bool
}

type UserCreateReq struct {
	Username string
	Password string
}

type UserLoginReq struct {
	Username string
	Password string
}

func (dest *Dest) ToMap() *map[string]interface{} {
	destMap := map[string]interface{}{
		"Dest":            dest.Dest,
		"Description":     dest.Description,
		"Password":        dest.Password,
		"Owner":           dest.Owner,
		"EnableAnalytics": dest.EnableAnalytics,
	}
	return &destMap
}
