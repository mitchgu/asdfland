package main

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

func GetRequestFingerprint(r *http.Request) string {
	return strings.Split(r.RemoteAddr, ":")[0]
}

func GetRandString(l int) string {
	token := make([]byte, int((l*3+3)/4))
	crand.Read(token)
	str := base64.RawURLEncoding.EncodeToString(token)
	return str[:l]
}

func initWordlists() map[string][]string {
	fname_regex := regexp.MustCompile(`^wordlists\/(.*)\.txt$`)
	wls := make(map[string][]string)
	for _, wl_name := range AssetNames() {
		fname_match := fname_regex.FindStringSubmatch(wl_name)
		if fname_match == nil {
			continue
		}
		wl_bytes := MustAsset(wl_name)
		scanner := bufio.NewScanner(bytes.NewReader(wl_bytes))
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		wls[fname_match[1]] = lines
		log.Printf("Wordlist loaded: %s", wl_name)
	}
	return wls
}

var wordlists = initWordlists()

func GetReadableString(wl_name string, count int) (string, error) {
	var words []string
	wl, exists := wordlists[wl_name]
	if !exists {
		return "", fmt.Errorf("GetReadableString: word list %s not found", wl)
	}
	wl_size := len(wl)
	for i := 0; i < count; i++ {
		words = append(words, wl[rand.Intn(wl_size)])
	}
	return strings.Join(words, c.GetString("word_delimiter")), nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondBadRequest(w http.ResponseWriter, reason string) {
	respondWithJSON(w, http.StatusBadRequest, map[string]string{
		"success": "false",
		"msg":     reason,
	})
}

func respondServerError(w http.ResponseWriter, reason string) {
	respondWithJSON(w, http.StatusInternalServerError, map[string]string{
		"success": "false",
		"msg":     reason,
	})
}

func respondOK(w http.ResponseWriter) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"success": "true",
	})
}
