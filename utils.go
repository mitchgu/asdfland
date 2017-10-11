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

var googleWords, _ = readLines("wordlists/google.txt")

func GetDictionaryWord() string {
	return googleWords[rand.Intn(len(googleWords))]
}

func GetReadableString(len int) string {
	slug := ""
	for i := 0; i < len; i++ {
		slug += GetDictionaryWord()
	}
	return slug
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
