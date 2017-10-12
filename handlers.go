package main

import (
	// "log"
	// "fmt"
    "strings"
    "io/ioutil"
	"encoding/json"

	"net/http"
    "github.com/gorilla/mux"
)

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
    respondBadRequest(w, "could not reserve a slug witn provided params")
}

func (a *App) SlugDestCreateHandler(w http.ResponseWriter, r *http.Request) {
    var sdcr SlugDestCreateReq
    var dest Dest
    bodyBuf, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(bodyBuf, &sdcr)
    errDest := json.Unmarshal(bodyBuf, &dest)
    if err != nil || errDest != nil {
        respondBadRequest(w, "malformed JSON in request2" + err.Error())
        return
    }
    fingerprint := GetRequestFingerprint(r)
    if (!a.DB.SlugReserved(fingerprint, sdcr.Slug)) {
        respondBadRequest(w, "slug hasn't been reserved yet" + sdcr.Slug)
        return
    }
    destUUID, success := a.DB.DestCreate(&dest)
    if (!success) {
        respondServerError(w, "problem creating destination")
        return
    }
    success = a.DB.SlugCreate(sdcr.Slug, destUUID, sdcr.Expire, fingerprint)
    if (!success) {
        respondServerError(w, "problem creating slug")
        return
    }
    respondOK(w)
}

func (a *App) KeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dest, err := a.DB.SlugFollow(vars["slug"])
	if err != nil {
	    http.Error(w, "Page not found", 404)
	} else {
        if strings.Index(dest, "http://") != 0 || strings.Index(dest, "https://") != 0 {
            dest = "http://" + dest
        }
		http.Redirect(w, r, dest, 302)
	}
}
