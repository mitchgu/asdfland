package main

import (
	// "log"
	// "fmt"
	"encoding/json"

	"net/http"
)

type SlugReserveReq struct {
	Type       string
	Length     int
	Dictionary string
	CustomSlug string
}

func (a *App) SlugReserveHandler(w http.ResponseWriter, r *http.Request) {
	var srr SlugReserveReq
	err := json.NewDecoder(r.Body).Decode(&srr)
	if err != nil {
		respondBadRequest(w, "malformed JSON in request")
		return
	}
	fingerprint := GetRequestFingerprint(r)
	var slugGenerator func() string
	var attempts int
	switch srr.Type {
	case "random":
		if srr.Length < 6 {
			respondBadRequest(w, "random slug length must be >=6")
			return
		}
		slugGenerator = func() string { return GetRandString(srr.Length) }
		attempts = 3
	case "readable":
		if srr.Length < 1 || srr.Length > 6 {
			respondBadRequest(w, "readable slug length must be 1 to 6 words")
			return
		}
		slugGenerator = func() string { return GetReadableString(srr.Length) }
		attempts = 5
	case "custom":
		if len(srr.CustomSlug) < 1 {
			respondBadRequest(w, "provided custom slug is empty")
			return
		}
		slugGenerator = func() string { return srr.CustomSlug }
		attempts = 1
	default:
		attempts = 0
	}
	for i := 0; i < attempts; i++ {
		slug := slugGenerator()
		success := a.DB.ReserveSlug(fingerprint, slug)
		if success {
			respondWithJSON(w, http.StatusOK, map[string]string{"success": "true", "slug": slug})
			return
		}
	}
	respondWithJSON(w, http.StatusBadRequest, map[string]string{"success": "false"})
}

func (a *App) KeyHandler(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// dest, err := db.Get(vars["key"]).Result()
	// if err != nil {
	//        http.Error(w, "Page not found", 404)
	// } else {
	// 	http.Redirect(w, r, dest, 302)
	// }
}

func (a *App) DestCreateHandler(w http.ResponseWriter, r *http.Request) {
	// var d Dest
	// err := json.NewDecoder(r.Body).Decode(&d)
	// if err != nil {
	//     http.Error(w, "Error decoding JSON", 400)
	//     return
	// }
	// d.CreationIP = r.RemoteAddr
	// success := DestCreate(d)
	// if !success {
	// 	http.Error(w, "Error creating destination", 500)
	// 	return
	// }
	// w.Write([]byte("OK"))
}
