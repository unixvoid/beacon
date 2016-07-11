package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v3"
)

func provision(w http.ResponseWriter, r *http.Request, client *redis.Client, state string) {
	// get file POST from index
	r.ParseForm()
	clientId := strings.TrimSpace(r.FormValue("id"))

	// check if clientId is set
	if len(clientId) == 0 {
		glogger.Debug.Println("id not set, exiting")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// sha3:512 hash the id
	clientIdHash := sha3.Sum512([]byte(clientId))

	// check if the id exists (sec:<clientId>)
	_, err := client.Get(fmt.Sprintf("sec:%x", clientIdHash)).Result()
	if err != redis.Nil {
		glogger.Debug.Println("COLLISION")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// generate token, store hashed token in db
		token := randStr(config.Beacon.TokenSize)
		tokenHash := sha3.Sum512([]byte(token))

		// return token to client
		w.Header().Set("token", token)
		fmt.Fprintf(w, "%s", token)

		// done with client, rest is server side
		// sec:<hashed clientId> : hashed password
		client.Set(fmt.Sprintf("sec:%x", clientIdHash), fmt.Sprintf("%x", tokenHash), 0).Err()

		// if temp objects, set an expire link on them
		if strings.EqualFold(state, "tmp") {
		}
	}
}
