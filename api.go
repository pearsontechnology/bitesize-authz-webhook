package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Handlers handles HTTP responses
func Handlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("POST")
	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var resp *AuthorizationResponse

	w.Header().Set("Content-type", "application/json")
	rbody, err := NewAuthorizationRequest(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user := NewAuthzUser(rbody)

	if !user.IsAllowed() {
		// w.WriteHeader(http.StatusForbidden)
		b, _ := json.Marshal(user.Request())
		log.Printf("User %s forbidden, request: %s", user.Username(), string(b))
		resp = NewAuthorizationResponse(false, "Not allowed")
	} else {
		resp = NewAuthorizationResponse(true)
	}
	json.NewEncoder(w).Encode(resp)
}
