package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (a *App) SlugReserveHandler(w http.ResponseWriter, r *http.Request) {
	var srr SlugReserveReq
	err := json.NewDecoder(r.Body).Decode(&srr)
	if err != nil {
		respondBadRequest(w, "Malformed JSON in request "+err.Error())
		return
	}
	fingerprint := r.Context().Value("Username").(string)
	var slugGenerator func() (string, error)
	var attempts int
	switch srr.Type {
	case "random":
		if srr.Length < 6 {
			respondBadRequest(w, "Random slug length must be >=6")
			return
		}
		slugGenerator = func() (string, error) { return GetRandString(srr.Length), nil }
		attempts = 3
	case "readable":
		if srr.Length < 1 || srr.Length > 6 {
			respondBadRequest(w, "Readable slug length must be 1 to 6 words")
			return
		}
		slugGenerator = func() (string, error) { return GetReadableString(srr.Wordlist, srr.Length) }
		attempts = 5
	case "custom":
		if len(srr.CustomSlug) < 1 {
			respondBadRequest(w, "Provided custom slug is empty")
			return
		}
		slugGenerator = func() (string, error) { return srr.CustomSlug, nil }
		attempts = 1
	default:
		attempts = 0
	}
	for i := 0; i < attempts; i++ {
		if slug, err := slugGenerator(); err != nil {
			respondBadRequest(w, err.Error())
			return
		} else {
			success := a.DB.ReserveSlug(fingerprint, slug)
			if success {
				respondWithJSON(w, http.StatusOK, map[string]string{"success": "true", "slug": slug})
				return
			}
		}
	}
	respondBadRequest(w, "Could not reserve a slug witn provided params")
}

func (a *App) SlugDestCreateHandler(w http.ResponseWriter, r *http.Request) {
	var sdcr SlugDestCreateReq
	var dest Dest
	bodyBuf, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBuf, &sdcr)
	errDest := json.Unmarshal(bodyBuf, &dest)
	if err != nil || errDest != nil {
		respondBadRequest(w, "Malformed JSON in request "+err.Error())
		return
	}
	log.Print(dest.Dest)
	destUrl, err := url.Parse(dest.Dest)
	if err != nil {
		respondBadRequest(w, "Destination URL could not be parsed")
		return
	}
	if !destUrl.IsAbs() {
		destUrl.Scheme = "http"
	}
	dest.Dest = destUrl.String()
	fingerprint := r.Context().Value("Username").(string)
	if !a.DB.SlugReserved(fingerprint, sdcr.Slug) {
		respondBadRequest(w, "Slug hasn't been reserved yet"+sdcr.Slug)
		return
	}
	dest.Username = r.Context().Value("Username").(string)
	destUUID, success := a.DB.DestCreate(&dest)
	if !success {
		respondServerError(w, "Problem creating destination")
		return
	}
	success = a.DB.SlugCreate(sdcr.Slug, destUUID, sdcr.Expire, fingerprint)
	if !success {
		respondServerError(w, "Problem creating slug")
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
		http.Redirect(w, r, dest, 302)
	}
}

func (a *App) KeyPreviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dest, err := a.DB.SlugFollow(vars["slug"])
	if err != nil {
		http.Error(w, "Page not found", 404)
	} else {
		fmt.Fprintf(w, "Destination is %s", dest)
	}
}

func (a *App) UserCreateHandler(w http.ResponseWriter, r *http.Request) {
	var ucr UserCreateReq
	if err := json.NewDecoder(r.Body).Decode(&ucr); err != nil {
		respondBadRequest(w, "Malformed JSON in request "+err.Error())
		return
	}
	if len(ucr.Password) < 6 {
		respondBadRequest(w, "Password must be at least 6 characters")
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(ucr.Password), c.GetInt("bcrypt_cost"))
	digest := string(bytes)
	err = a.DB.UserCreate(ucr.Username, digest)
	if err != nil {
		respondBadRequest(w, err.Error())
		return
	}
	token := GetRandString(32)
	a.DB.SessionCreate(ucr.Username, token)
	cookie := http.Cookie{Name: "session", Value: token, Path: "/", Expires: time.Now().Add(30 * 24 * time.Hour)}
	http.SetCookie(w, &cookie)
	respondOK(w)
}

func (a *App) UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	var ulr UserLoginReq
	if r.Context().Value("IsRegistered").(bool) {
		respondBadRequest(w, "Already logged in")
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&ulr); err != nil {
		respondBadRequest(w, "Malformed JSON in request "+err.Error())
		return
	}
	digest, err := a.DB.UserGetDigest(ulr.Username)
	if err != nil {
		respondBadRequest(w, "User not found")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(digest), []byte(ulr.Password))
	token := GetRandString(32)
	a.DB.SessionCreate(ulr.Username, token)
	cookie := http.Cookie{Name: "session", Value: token, Path: "/", Expires: time.Now().Add(30 * 24 * time.Hour)}
	http.SetCookie(w, &cookie)
	respondWithJSON(w, http.StatusOK, map[string]string{
		"success": strconv.FormatBool(err == nil),
	})
}

func (a *App) UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx.Value("IsRegistered").(bool) {
		// overwrite old session cookie with new unregistered one
		token := GetRandString(32)
		username := GetRandString(32)
		a.DB.SessionCreate(username, token)
		cookie := http.Cookie{Name: "session", Value: token, Path: "/", Expires: time.Now().Add(30 * 24 * time.Hour)}
		http.SetCookie(w, &cookie)
		ctx = context.WithValue(ctx, "IsLoggedIn", true)
		ctx = context.WithValue(ctx, "Username", username)
		ctx = context.WithValue(ctx, "IsRegistered", false)
	}
	respondOK(w)
}

func (a *App) SessionGetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	respondWithJSON(w, http.StatusOK, map[string]string{
		"Username":     ctx.Value("Username").(string),
		"IsRegistered": strconv.FormatBool(ctx.Value("IsRegistered").(bool)),
	})
}
