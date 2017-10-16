package main

import (
	"bufio"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"log"
	"fmt"
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

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func initWordlists() (map[string][]string) {
	wls := make(map[string][]string)
	for name, file := range c.GetStringMapString("wordlists") {
		if wordlist, err := readLines(file); err != nil {
			log.Printf("Word list %s at %s could not be read", name, file)
		} else {
			wls[name] = wordlist
		}
	}
	return wls
}

var wordlists = initWordlists()

func GetWordlistWord(wl string) (string, error) {
	if wordlist, ok := wordlists[wl]; ok {
		return wordlist[rand.Intn(len(wordlist))], nil
	} else {
		return "", fmt.Errorf("GetWordlistWord: word list %s not found", wl)
	}
}

func GetReadableString(wl string, len int) (string, error) {
	var words []string
	for i := 0; i < len; i++ {
		if word, err := GetWordlistWord(wl); err != nil {
			return "", err
		} else {
			words = append(words, word)
		}
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
