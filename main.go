package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	// "github.com/go-redis/redis"

	"github.com/ajarv/go-app/api"
	"github.com/ajarv/go-app/tlsserver"
	"github.com/ajarv/go-app/workflow"
	"github.com/gorilla/mux"
)

func killHandler(w http.ResponseWriter, r *http.Request) {
	go func() {
		log.Printf("Good Bye World !")
		time.Sleep(4 * time.Second)
		os.Exit(3)
	}()

	data := api.GetDebugData(r)
	defer api.WriteResponse(w, r, data)
	data["message"] = "Will terminate myself on your request in a few .. good bye !!"
}

func healthz(w http.ResponseWriter, r *http.Request) {
	message := "ok"
	w.Write([]byte(message))
}

func registerRoutes(dir string) *mux.Router {
	r := mux.NewRouter()
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(dir)))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/die", killHandler)
	r.HandleFunc("/api/v1/info", api.IndexHandler)
	r.HandleFunc("/api/v2/{type}", api.ApiResourceHandler)
	r.HandleFunc("/api/v2/{type}/{id}", api.ApiResourceHandler)
	// r.HandleFunc("/redis", redisHandler)
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/workflow", workflow.WorkflowHandler)
	r.HandleFunc("/", api.IndexHandler)
	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		requestDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(requestDump))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
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
	router.Use(loggingMiddleware)

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
