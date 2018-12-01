package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"html/template"
)

const appVersion = "3.0.0"

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func logRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("--Request :%v\n", string(requestDump))
}

func getDebugResponseString(req *http.Request) string {
	hdrs, _ := json.Marshal(req.Header)
	var hostname, _ = os.Hostname()
	const stemplate = `---
App Version : %s
Server Host : %s
<< Request 
Headers:
%s
---

`
	return fmt.Sprintf(stemplate, appVersion, hostname, string(hdrs))
}

func handler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Write([]byte(getDebugResponseString(r)))
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	w.Write([]byte(getDebugResponseString(r)))
	w.Write([]byte("\n\nWill terminate myself on your request in a few .. good bye !!"))

	go func() {
		time.Sleep(4 * time.Second)
		os.Exit(3)
	}()
}

type RedisPageData struct {
	PageTitle string
	Count int64
	Error error
	DebugString string
}

var tmpl = template.Must(template.ParseFiles("templates/layout.html"))

func redishandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	// w.Write([]byte(getDebugResponseString(r)))

	redisdb := redis.NewClient(&redis.Options{
		Addr:        redisHost + ":" + redisPort,
		Password:    "", // no password set
		DB:          0,  // use default DB
		DialTimeout: time.Second * 5,
	})
	
	result, err := redisdb.Incr("hits.go").Result()
	data := RedisPageData {
		PageTitle : "Redis Page",
		Count: result,
		Error: err,
		DebugString:  getDebugResponseString(r),
	}
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)

	// if err != nil {
	// 	fmt.Fprintf(w, "Redis Error: %v\n", err)
	// } else {
	// 	fmt.Fprintf(w, "(redis) Hit Count : %v\n", result)
	// }
}

var redisHost = getEnv("REDIS_SERVICE_HOST", "localhost")
var redisPort = getEnv("REDIS_SERVICE_PORT", "6379")

func main() {
	var host string
	var dir string
	var port string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&host, "host", "0.0.0.0", "listen host")
	flag.StringVar(&port, "port", "8080", "listen port")

	flag.Parse()

	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler( http.FileServer(http.Dir(dir)))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/", handler)
	r.HandleFunc("/redis", redishandler)
	r.HandleFunc("/kill", killHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    host + ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Fprintf(os.Stdout, "Server listening %s:%s\n", host, port)

	log.Fatal(srv.ListenAndServe())
}
