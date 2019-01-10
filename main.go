package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/gojsonq"
	yaml "gopkg.in/yaml.v2"
)

func logRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("--Request :\n%v\n------------\n", string(requestDump))
}

func getDebugData(req *http.Request) map[string]interface{} {
	var hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknown host"
	}

	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			if strings.HasPrefix(key, "secret") {
				val = b64.StdEncoding.EncodeToString([]byte(val))
				val = b64.StdEncoding.EncodeToString([]byte(val))
			}
			items[key] = val
		}
		return items
	}

	data := map[string]interface{}{
		"AppColor":   appColor,
		"Host":       hostname,
		"ApiVersion": appVersion,
		"AppName":    appName,
		"ServerTime": time.Now().String(),
		"Request":    map[string]interface{}{"Headers": req.Header},
	}
	if viewenv := req.URL.Query().Get("showenv"); viewenv != "" {
		environment := getenvironment(os.Environ(), func(item string) (key, val string) {
			splits := strings.Split(item, "=")
			key = splits[0]
			val = splits[1]
			return
		})
		data["Environment"] = environment
	}

	return data
}

var tmpl = template.Must(template.ParseFiles("templates/layout.html"))

func writeData(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	if strings.Contains(r.Header["Accept"][0], "html") {
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, data)
		if err != nil {
			w.Write([]byte(`{"result":"Error"}`))
		}
		return
	}

	if strings.Contains(r.Header["Accept"][0], "json") {
		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(&data)
		if err != nil {
			w.Write([]byte(`{"result":"Error"}`))
			return
		}
		w.Write(b)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	b, err := yaml.Marshal(&data)
	if err != nil {
		w.Write([]byte(`{"result":"Error"}`))
		return
	}
	w.Write(b)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data := getDebugData(r)
	writeData(w, r, data)
}

type formData struct {
	APIVersion string `json:"apiVersion"`
	Step       int    `json:"step"`
}

func workflowHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	var form formData

	decoder := json.NewDecoder(r.Body)
	data := getDebugData(r)

	data["Workflow"] = &form
	defer writeData(w, r, data)

	err := decoder.Decode(&form)
	if err != nil {
		data["warning"] = fmt.Sprintf("Unable to parse request %v", err)
		return
	}

	switch {
	case form.Step == 0:
		form.Step = 1
		form.APIVersion = appVersion
	case form.Step > 0 && form.APIVersion != appVersion:
		data["warning"] = fmt.Sprintf("Protocol Version Mismatch  Client - %v vs   Server - %v ", form.APIVersion, appVersion)
	case form.Step >= 3:
		form.Step = 3
		data["warning"] = fmt.Sprintf("Order already confirmed. no modifications possible")
	default:
		form.Step = form.Step + 1
		if form.Step == 3 {
			data["message"] = fmt.Sprintf("Order confirmed")
		}
	}

}

func killHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		log.Printf("Good Bye World !")
		time.Sleep(4 * time.Second)
		os.Exit(3)
	}()

	logRequest(r)
	data := getDebugData(r)
	data["message"] = "Will terminate myself on your request in a few .. good bye !!"
	writeData(w, r, data)
}

func redisHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data := getDebugData(r)
	redisdb := redis.NewClient(&redis.Options{
		Addr:        redisHost + ":" + redisPort,
		Password:    redisPassword, // no password set
		DB:          0,             // use default DB
		DialTimeout: time.Second,
	})

	count, err := redisdb.Incr("hits.go").Result()
	if err != nil {
		data["warning"] = fmt.Sprintf("Unable to connect redis server at %s:%s", redisHost, redisPort)
	} else {
		data["message"] = fmt.Sprintf("Hit Count from Redis key hits.go : %v", count)
	}
	writeData(w, r, data)
}

var redisHost = getEnv("REDIS_SERVICE_HOST", "localhost")
var redisPort = getEnv("REDIS_SERVICE_PORT", "6379")
var redisPassword = getEnv("REDIS_SERVICE_PASSWORD", "")

 

var appVersion = getEnv("APP_VERSION", "v1.0.0")
var appName = getEnv("APP_NAME", "GO_WEB")
var appColor = getEnv("APP_COLOR", "black")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseCFSettings() {
	if value, ok := os.LookupEnv("VCAP_SERVICES"); ok {
		if v := gojsonq.New().JSONString(value).Find("rediscloud.[0].credentials.hostname"); v != nil {
			redisHost = fmt.Sprintf("%v", v)
		}
		if v := gojsonq.New().JSONString(value).Find("rediscloud.[0].credentials.port"); v != nil {
			redisPort = fmt.Sprintf("%v", v)
		}
		if v := gojsonq.New().JSONString(value).Find("rediscloud.[0].credentials.password"); v != nil {
			redisPassword = fmt.Sprintf("%v", v)
		}
		fmt.Printf("Redis host:%v | port:%v ", redisHost, redisPort)
	}
}

func main() {
	var host string
	var dir string
	var port string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&host, "host", "0.0.0.0", "listen host")
	flag.StringVar(&port, "port", "8080", "listen port")
	flag.Parse()

	parseCFSettings()

	r := mux.NewRouter()
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(dir)))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/die", killHandler)
	r.HandleFunc("/redis", redisHandler)
	r.HandleFunc("/workflow", workflowHandler)

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
