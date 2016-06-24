package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/sha3"
	"gopkg.in/gcfg.v1"
	"gopkg.in/redis.v3"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
//================================================================
// general strategy:
// we take in a file, the filename is a hashed random string.
// the file is stored with its filename as the hased string.
// the random string (token) is returned back to the user.
//
// now when the user wants to retrive the file, he puts in the
// token (random string from earlier). his request is hashed and
// the stored hash is returned. ez
//================================================================
*/

type Config struct {
	Beacon struct {
		Port            int
		TokenSize       int
		LinkTokenSize   int
		TokenDictionary string
		TTL             time.Duration
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
	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		upload(w, r, client, "tmp")
	}).Methods("POST")
	router.HandleFunc("/{fdata}", func(w http.ResponseWriter, r *http.Request) {
		handlerdynamic(w, r, client)
	}).Methods("GET")
	log.Fatal(http.ListenAndServe(bitport, router))
}

func handlerdynamic(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	vars := mux.Vars(r)
	fdata := vars["fdata"]

	// hash the token that is passed
	hash := sha3.Sum512([]byte(fdata))
	hashstr := fmt.Sprintf("%x", hash)

	val, err := client.Get(hashstr).Result()
	if err != nil {
		log.Printf("data does not exist")
		fmt.Fprintf(w, "token not found")
	} else {
		//log.Printf("data exists")
		ip := strings.Split(r.RemoteAddr, ":")[0]
		log.Printf("Responsing to %s :: from: %s", fdata, ip)

		decodeVal, _ := base64.StdEncoding.DecodeString(val)

		file, _ := os.Create("tmpfile")
		io.WriteString(file, string(decodeVal))
		file.Close()

		http.ServeFile(w, r, "tmpfile")
		os.Remove("tmpfile")
	}
}

func upload(w http.ResponseWriter, r *http.Request, client *redis.Client, state string) {
	// get file POST from index
	//fmt.Println("method:", r.Method)
	r.ParseForm()
	address := strings.TrimSpace(r.FormValue("address"))

	// generate token and hash to save
	token := tokenGen(config.Beacon.TokenSize, client)
	w.Header().Set("token", token)
	fmt.Fprintf(w, "%s", token)

	// done with client, rest is server side
	hash := sha3.Sum512([]byte(token))
	hashstr := fmt.Sprintf("%x", hash)
	fmt.Println("uploading:", token)

	fileBase64Str := base64.StdEncoding.EncodeToString([]byte(address))

	//println("uploading ", "file")
	client.Set(hashstr, fileBase64Str, 0).Err()
	if strings.EqualFold(state, "tmp") {
		client.Expire(hashstr, (config.Beacon.TTL * time.Hour)).Err()
		//fmt.Println("expire link generated")
	}
}

func tokenGen(strSize int, client *redis.Client) string {
	// generate new token
	token := randStr(strSize)
	// hash token
	hash := sha3.Sum512([]byte(token))
	hashstr := fmt.Sprintf("%x", hash)

	_, err := client.Get(hashstr).Result()

	for err != redis.Nil {
		fmt.Println("DEBUG :: COLLISION")
		token = randStr(strSize)
		hash := sha3.Sum512([]byte(token))
		hashstr := fmt.Sprintf("%x", hash)

		_, err = client.Get(hashstr).Result()
		// do not ddos box if db is full
		time.Sleep(time.Second * 1)
	}
	return token
}

func randStr(strSize int) string {
	//dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	dictionary := config.Beacon.TokenDictionary

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
