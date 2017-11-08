package main

import (
// "time"
)

type Dest struct {
	Dest            string
	Description     string
	Password        string
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

func (dest *Dest) ToMap() *map[string]interface{} {
	destMap := map[string]interface{}{
		"Dest":            dest.Dest,
		"Description":     dest.Description,
		"Password":        dest.Password,
		"EnableAnalytics": dest.EnableAnalytics,
	}
	return &destMap
}
