package main

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/sha3"
	"gopkg.in/redis.v3"
)

func remove(w http.ResponseWriter, r *http.Request, client *redis.Client, state string) {
	// get file POST from index
	r.ParseForm()
	clientId := strings.TrimSpace(r.FormValue("id"))
	clientSec := strings.TrimSpace(r.FormValue("sec"))

	// sha3:512 hash the id and sec
	clientIdHash := sha3.Sum512([]byte(clientId))
	clientSecHash := sha3.Sum512([]byte(clientSec))

	// check if id exists
	storedSecHash, err := client.Get(fmt.Sprintf("sec:%x", clientIdHash)).Result()
	if err != redis.Nil {
		// id exists, make sure clientSecHash is the same as the stored version
		if fmt.Sprintf("%x", clientSecHash) == storedSecHash {
			// client is authed
			w.WriteHeader(http.StatusOK)
			client.Del(fmt.Sprintf("ip:%x", clientIdHash))
			client.Del(fmt.Sprintf("sec:%x", clientIdHash))
		} else {
			// client auth failed
			w.WriteHeader(http.StatusForbidden)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		// id does not exist
	}
}
