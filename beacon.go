package main

import (
	"crypto/rand"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/sha3"
	"gopkg.in/gcfg.v1"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	Beacon struct {
		Port            int
		TokenSize       int
		TokenDictionary string
	}
	SSL struct {
		UseTLS     bool
		ServerCert string
		ServerKey  string
	}
	Redis struct {
		Host string
		Port int
	}
}

var (
	config = Config{}
)

func main() {
	err := gcfg.ReadFileInto(&config, "config.gcfg")
	if err != nil {
		fmt.Printf("Could not load config.gcfg, error: %s\n", err)
		return
	}

	redisaddr := fmt.Sprint(config.Redis.Host, ":", config.Redis.Port)
	bitport := fmt.Sprint(":", config.Beacon.Port)
	println("beacon running on", config.Beacon.Port)
	println("link to redis on", redisaddr)
	// initialize redis connection
	client := redis.NewClient(&redis.Options{
		Addr:     redisaddr,
		Password: "",
		DB:       0,
	})

	// all handlers. lookin funcy casue i have to pass redis handler
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/provision", func(w http.ResponseWriter, r *http.Request) {
		provision(w, r, client, "tmp")
	}).Methods("POST")
	router.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		update(w, r, client, "tmp")
	}).Methods("POST")
	router.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		remove(w, r, client, "tmp")
	}).Methods("POST")
	router.HandleFunc("/{fdata}", func(w http.ResponseWriter, r *http.Request) {
		handlerdynamic(w, r, client)
	}).Methods("GET")
	if config.SSL.UseTLS {
		log.Fatal(http.ListenAndServeTLS(bitport, config.SSL.ServerCert, config.SSL.ServerKey, router))
	} else {
		log.Fatal(http.ListenAndServe(bitport, router))
	}
}

func handlerdynamic(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	vars := mux.Vars(r)
	fdata := vars["fdata"]

	// hash the token that is passed
	clientIdHash := sha3.Sum512([]byte(fdata))

	val, err := client.Get(fmt.Sprintf("ip:%x", clientIdHash)).Result()
	if err != nil {
		log.Printf("data does not exist")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "token not found")
	} else {
		//log.Printf("data exists")
		fmt.Fprintf(w, "%s", val)
	}
}

func provision(w http.ResponseWriter, r *http.Request, client *redis.Client, state string) {
	// get file POST from index
	r.ParseForm()
	clientId := strings.TrimSpace(r.FormValue("id"))

	// sha3:512 hash the id
	clientIdHash := sha3.Sum512([]byte(clientId))

	// check if the id exists (sec:<clientId>)
	_, err := client.Get(fmt.Sprintf("sec:%x", clientIdHash)).Result()
	if err != redis.Nil {
		fmt.Println("DEBUG :: COLLISION")
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

func update(w http.ResponseWriter, r *http.Request, client *redis.Client, state string) {
	// get file POST from index
	r.ParseForm()
	clientId := strings.TrimSpace(r.FormValue("id"))
	clientSec := strings.TrimSpace(r.FormValue("sec"))
	clientAddress := strings.TrimSpace(r.FormValue("address"))

	// sha3:512 hash the id and sec
	clientIdHash := sha3.Sum512([]byte(clientId))
	clientSecHash := sha3.Sum512([]byte(clientSec))

	// check if id exists
	storedSecHash, err := client.Get(fmt.Sprintf("sec:%x", clientIdHash)).Result()
	if err != redis.Nil {
		// id exists, make sure clientSecHash is the same as the stored version
		if fmt.Sprintf("%x", clientSecHash) == storedSecHash {
			// client is authed, update clientAddress
			w.WriteHeader(http.StatusOK)
			client.Set(fmt.Sprintf("ip:%x", clientIdHash), clientAddress, 0).Err()
		} else {
			// client auth failed
			w.WriteHeader(http.StatusForbidden)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		// id does not exist
	}
}

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

func randStr(strSize int) string {
	dictionary := config.Beacon.TokenDictionary

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
