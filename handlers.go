package main

import (
	// "log"
	"fmt"
	// "encoding/json"

	"net/http"
)

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