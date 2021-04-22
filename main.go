package main

import (
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

	// "github.com/go-redis/redis"

	"github.com/ajarv/go-app/api"
	"github.com/ajarv/go-app/tlsserver"
	"github.com/gorilla/mux"
	yaml "gopkg.in/yaml.v2"
)

func logRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("--Request :\n%v\n------------\n", string(requestDump))
}

var tmpl = template.Must(template.ParseFiles("templates/layout.html"))

func writeData(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	if len(r.Header["Accept"]) > 0 {
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
	data := api.GetDebugData(r)
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
	data := api.GetDebugData(r)

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

func apiInfoHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data := api.GetDebugData(r)
	data["info"] = appInfo
	writeData(w, r, data)
}

func apiResourceHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	vars := mux.Vars(r)
	data := api.GetDebugData(r)
	data["resource"] = vars
	writeData(w, r, data)
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		log.Printf("Good Bye World !")
		time.Sleep(4 * time.Second)
		os.Exit(3)
	}()

	logRequest(r)
	data := api.GetDebugData(r)
	data["message"] = "Will terminate myself on your request in a few .. good bye !!"
	writeData(w, r, data)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	message := "ok"
	w.Write([]byte(message))
}

var appVersion = getEnv("APP_VERSION", "v1.0.0")
var appName = getEnv("APP_NAME", "GO_WEB")
var appColor = getEnv("APP_COLOR", "black")
var appInfo = getEnv("APP_INFO", "-")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func registerRoutes(dir string) *mux.Router {
	r := mux.NewRouter()
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(dir)))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/die", killHandler)
	r.HandleFunc("/api/v1/info", apiInfoHandler)
	r.HandleFunc("/api/v2/{type}", apiResourceHandler)
	r.HandleFunc("/api/v2/{type}/{id}", apiResourceHandler)
	// r.HandleFunc("/redis", redisHandler)
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/workflow", workflowHandler)
	r.HandleFunc("/", indexHandler)
	return r
}

func main() {
	var appAdmin string
	var host string
	var dir string
	var port string
	var secure bool
	var serverkey string
	var serverCert string

	flag.StringVar(&appAdmin, "admin", "raviraaj", "Dummy var to signify the admin user")
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&host, "host", "0.0.0.0", "listen host")
	flag.StringVar(&port, "port", "8080", "listen port")
	flag.StringVar(&serverkey, "sslKey", "", "ssl key file")
	flag.StringVar(&serverCert, "sslCert", "", "ssl cert file")
	flag.BoolVar(&secure, "secure", false, "ssl listener")
	flag.Parse()
	var err error
	if secure {
		if !strings.HasSuffix(port, "443") {
			port = "8443"
		}
		if serverkey == "" {
			serverkey, err = tlsserver.GetDefaultKeyFilePath()
			if err != nil {
				log.Fatal("Unable to create default key file", err)
			}
			serverCert, err = tlsserver.GetDefaultCertFilePath()
			if err != nil {
				log.Fatal("Unable to create default cert file", err)
			}
			defer os.Remove(serverkey)
			defer os.Remove(serverCert)
		}
	}

	router := registerRoutes(dir)

	srv := &http.Server{
		Handler: router,
		Addr:    host + ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Fprintf(os.Stdout, "Admin User :%s\n", appAdmin)
	fmt.Fprintf(os.Stdout, "Server listening HTTP %s:%s\n", host, port)
	if secure {
		log.Fatal(srv.ListenAndServeTLS(serverCert, serverkey))
		return
	}
	log.Fatal(srv.ListenAndServe())
}
