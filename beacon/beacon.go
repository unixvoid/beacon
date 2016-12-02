package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/unixvoid/glogger"
	"golang.org/x/crypto/sha3"
	"gopkg.in/gcfg.v1"
	"gopkg.in/redis.v3"
)

type Config struct {
	Beacon struct {
		Port            int
		TokenSize       int
		TokenDictionary string
		Loglevel        string
	}
	SSL struct {
		UseTLS     bool
		ServerCert string
		ServerKey  string
	}
	Redis struct {
		Host     string
		Password string
	}
}

var (
	config = Config{}
)

func main() {
	// read in config
	readConf()

	// init logger
	initLogger()

	// initialize redis connection
	client, err := initRedisConnection()
	if err != nil {
		glogger.Debug.Println("redis conneciton cannot be made, trying again in 1 second")
		client, err = initRedisConnection()
		if err != nil {
			glogger.Error.Println("redis connection cannot be made.")
			os.Exit(1)
		}
	}
	glogger.Debug.Println("connection to redis succeeded.")
	glogger.Info.Println("link to redis on", config.Redis.Host)

	// router routes/handlers
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/provision", func(w http.ResponseWriter, r *http.Request) {
		provision(w, r, client)
	}).Methods("POST")
	router.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		update(w, r, client)
	}).Methods("POST")
	router.HandleFunc("/rotate", func(w http.ResponseWriter, r *http.Request) {
		rotate(w, r, client)
	}).Methods("POST")
	router.HandleFunc("/remove", func(w http.ResponseWriter, r *http.Request) {
		remove(w, r, client)
	}).Methods("POST")
	router.HandleFunc("/{fdata}", func(w http.ResponseWriter, r *http.Request) {
		handlerdynamic(w, r, client)
	}).Methods("GET")

	if config.SSL.UseTLS {
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			ClientSessionCache: tls.NewLRUClientSessionCache(128),
		}
		glogger.Info.Println("beacon running https on", config.Beacon.Port)
		tlsServer := &http.Server{Addr: fmt.Sprintf(":%d", config.Beacon.Port), Handler: router, TLSConfig: tlsConfig}
		log.Fatal(tlsServer.ListenAndServeTLS(config.SSL.ServerCert, config.SSL.ServerKey))
	} else {
		glogger.Info.Println("beacon running http on", config.Beacon.Port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Beacon.Port), router))
	}
}

func readConf() {
	err := gcfg.ReadFileInto(&config, "config.gcfg")
	if err != nil {
		fmt.Printf("Could not load config.gcfg, error: %s\n", err)
		return
	}
}

func initLogger() {
	// init logger
	if config.Beacon.Loglevel == "debug" {
		glogger.LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else if config.Beacon.Loglevel == "cluster" {
		glogger.LogInit(os.Stdout, os.Stdout, ioutil.Discard, os.Stderr)
	} else if config.Beacon.Loglevel == "info" {
		glogger.LogInit(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
	} else {
		glogger.LogInit(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	}
}

func initRedisConnection() (*redis.Client, error) {
	// initialize redis connection
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       0,
	})
	_, redisErr := client.Ping().Result()
	return client, redisErr
}

func handlerdynamic(w http.ResponseWriter, r *http.Request, client *redis.Client) {
	vars := mux.Vars(r)
	fdata := vars["fdata"]

	// hash the token that is passed
	clientIdHash := sha3.Sum512([]byte(fdata))

	val, err := client.Get(fmt.Sprintf("ip:%x", clientIdHash)).Result()
	if err != nil {
		glogger.Debug.Printf("data does not exist")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "token not found")
	} else {
		//log.Printf("data exists")
		fmt.Fprintf(w, "%s", val)
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
