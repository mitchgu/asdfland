package main

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func GetRandString(l int) string {
	token := make([]byte, int((l*3+3)/4))
	crand.Read(token)
	str := base64.RawURLEncoding.EncodeToString(token)
	return str[:l]
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
