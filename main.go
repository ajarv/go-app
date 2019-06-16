package main

import (
	"crypto/tls"
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/yalp/jsonpath"
	yaml "gopkg.in/yaml.v2"
)

var webserverCertificate = `-----BEGIN CERTIFICATE-----
MIIDHjCCAoegAwIBAgIJAJH05SI8aNd0MA0GCSqGSIb3DQEBCwUAMIGnMQswCQYD
VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJU2FuIERpZWdv
MSAwHgYDVQQKDBdCbHVlIFNreSB3aGl0ZSBTYW5kIExMQzELMAkGA1UECwwCSVQx
HjAcBgNVBAMMFWJsdWUtc2t5LndoaXRlc2FuZC5pbzEgMB4GCSqGSIb3DQEJARYR
YWphcnZAZXhhbXBsZS5jb20wHhcNMTkwMjI4MDEzMTQxWhcNMjQwMjI3MDEzMTQx
WjCBpzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcM
CVNhbiBEaWVnbzEgMB4GA1UECgwXQmx1ZSBTa3kgd2hpdGUgU2FuZCBMTEMxCzAJ
BgNVBAsMAklUMR4wHAYDVQQDDBVibHVlLXNreS53aGl0ZXNhbmQuaW8xIDAeBgkq
hkiG9w0BCQEWEWFqYXJ2QGV4YW1wbGUuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQDwELpzGB7XElUH1cxJ8J7gOfRXoVDFu6hQ1f/YuHv1zAQk/zzCHoxy
VxU/ogtTEt7jqGYvHrIh4tjUwb2BNScAmWsPWkWlVBmG+pq1YenNlqOJ52tIZcZ5
YAu4V9N0Nz2nC8+ij7FrcuxbMEAFOVCf2JTy3iYatMm7b0ViPKmUvQIDAQABo1Aw
TjAdBgNVHQ4EFgQUg+vv2w007ilS0M3gDgyQ+0XvJCMwHwYDVR0jBBgwFoAUg+vv
2w007ilS0M3gDgyQ+0XvJCMwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOB
gQDgcGlcqIvRTQRyqgCgrurEG5KUy3tBlbqX6bnPsfMxqYcoWzIzLlrMWe8Hzlta
xJETxc+0v7wdj011z5anfxXonbTS5AF0NWFbtkC378tJ4Nm0N8FLO/+KV9iVtGKB
Ud84CYDn7LLKC5NnOm4MDmUCAdMKYqYNRhfmPSxxvLbStw==
-----END CERTIFICATE-----`

var webserverPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDwELpzGB7XElUH1cxJ8J7gOfRXoVDFu6hQ1f/YuHv1zAQk/zzC
HoxyVxU/ogtTEt7jqGYvHrIh4tjUwb2BNScAmWsPWkWlVBmG+pq1YenNlqOJ52tI
ZcZ5YAu4V9N0Nz2nC8+ij7FrcuxbMEAFOVCf2JTy3iYatMm7b0ViPKmUvQIDAQAB
AoGBANn4RmJUR0Q+R+haTifgi1DKLjoWpVE0ByqGc8viDeNqf2TcPt1+gUUcHpXt
Wtzt6GTKtSUZeOHdp8TduGQFz8cwVk04No6ygYVbV30VA9OYuYba0XxCyR0cKYf8
i2mV3pAfpobndec70wXTw6tn3GEIO/RISFY1cFWcu+scZyh9AkEA+Lmx/ooWD020
gDPOcXuzlQCIWoP/U+d40Cgi795YdpdDk2FDarkM5WGZnww4sym80HzAUn56kaS1
p4F60nBxcwJBAPcWMSbD8wFjlfOeyiWfVB0Ak1YREVI+IcpvMCQ7LD3TL914BxIW
vnDLO+0DOPZEYOxOkxOMeg9RY/nPW2MJlQ8CQQC+5uAL6vZlhpGcuKaiCXzbR05g
kuFdB9N9iODP1It3ckAmlUeGWUPhptie72Vxdf560tVWO8ddk9rtFv8rF6yrAkEA
oDmx0dOLR0FOwdYce90f7FatNEiJFO3Zd642Z6g/fi/ugA0PeLlq8TW5PG60h227
9EDXuvuDQ1+iFyJRvp0+HQJAHeireIT7yU+KY3GCdn2Vq4EjrxTu1MKn6loInkwo
JpF8fprPVvHTLj3c0+srGwlKgDzsOrNIh86D+nmA5ES/qw==
-----END RSA PRIVATE KEY-----`

type embeddedServer struct {

	/*
		custom struct that embeds golang's standard http.Server type
		another way of looking at this is that embeddedServer "inherits" from http.Server,
		though this is not strictly accurate. Have a look at note below for additional information
	*/

	http.Server
	webserverCertificate string
	webserverKey         string
}

func (srv *embeddedServer) ListenAndServeTLS(addr string) error {

	/*
		This is where we "hide" or "override" the default "ListenAndServeTLS" method so we modify it to accept
		hardcoded certificates and keys rather than the default filenames

		The default implementation of ListenAndServeTLS was obtained from:
		https://github.com/zenazn/goji/blob/master/graceful/server.go#L33

		and tls.X509KeyPair (http://golang.org/pkg/crypto/tls/#X509KeyPair) is used,
		rather than the default tls.LoadX509KeyPair
	*/

	config := &tls.Config{
		MinVersion: tls.VersionTLS10,
	}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair([]byte(srv.webserverCertificate), []byte(srv.webserverKey))
	if err != nil {
		return err
	}

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(conn, config)
	return srv.Serve(tlsListener)
}

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
		"ServerTime": time.Now().Format("2006-01-02 15:04:05"),
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

func apiInfoHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data := getDebugData(r)
	data["info"] = appInfo

	writeData(w, r, data)
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
func healthz(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	message := "ok"
	w.Write([]byte(message))
}

var redisHost = getEnv("REDIS_SERVICE_HOST", "localhost")
var redisPort = getEnv("REDIS_SERVICE_PORT", "6379")
var redisPassword = getEnv("REDIS_SERVICE_PASSWORD", "")

var appVersion = getEnv("APP_VERSION", "v1.0.0")
var appName = getEnv("APP_NAME", "GO_WEB")
var appColor = getEnv("APP_COLOR", "black")
var appInfo = getEnv("APP_INFO", "{}")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseCFSettings() {
	if value, ok := os.LookupEnv("VCAP_SERVICES"); ok {
		raw := []byte(value)
		hostnameFilter, _ := jsonpath.Prepare("$..credentials.hostname")
		portFilter, _ := jsonpath.Prepare("$..credentials.port")
		passwordFilter, _ := jsonpath.Prepare("$..credentials.password")

		var data interface{}
		if err := json.Unmarshal(raw, &data); err != nil {
			return
		}
		if v, err := hostnameFilter(data); err == nil {
			if a := v.([]interface{}); len(a) > 0 {
				redisHost = fmt.Sprintf("%v", a[0])
			}
		}
		if v, err := portFilter(data); err == nil {
			if a := v.([]interface{}); len(a) > 0 {
				redisPort = fmt.Sprintf("%v", a[0])
			}
		}
		if v, err := passwordFilter(data); err == nil {
			if a := v.([]interface{}); len(a) > 0 {
				redisPassword = fmt.Sprintf("%v", a[0])
			}
		}

		// if v := gojsonq.New().JSONString(value).Find("rediscloud.[0].credentials.hostname"); v != nil {
		// 	redisHost = fmt.Sprintf("%v", v)
		// }
		// if v := gojsonq.New().JSONString(value).Find("rediscloud.[0].credentials.port"); v != nil {
		// 	redisPort = fmt.Sprintf("%v", v)
		// }
		// if v := gojsonq.New().JSONString(value).Find("rediscloud.[0].credentials.password"); v != nil {
		// 	redisPassword = fmt.Sprintf("%v", v)
		// }
		fmt.Printf("Redis host:%v | port:%v ", redisHost, redisPort)
	}
}

func main() {
	var host string
	var dir string
	var port string
	var secure bool
	var serverkey string
	var serverCert string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&host, "host", "0.0.0.0", "listen host")
	flag.StringVar(&port, "port", "8080", "listen port")
	flag.StringVar(&serverkey, "key", "", "ssl key")
	flag.StringVar(&serverCert, "cert", "", "ssl cert")
	flag.BoolVar(&secure, "secure", false, "ssl listener")
	flag.Parse()

	parseCFSettings()

	r := mux.NewRouter()
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(dir)))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/die", killHandler)
	r.HandleFunc("/api/v1/info", apiInfoHandler)
	r.HandleFunc("/redis", redisHandler)
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/workflow", workflowHandler)
	r.HandleFunc("/", indexHandler)

	if secure {
		if !strings.HasSuffix(port, "443") {
			port = "8443"
		}
		if serverkey == "" {
			embeddedTLSserver := &embeddedServer{
				webserverCertificate: webserverCertificate,
				webserverKey:         webserverPrivateKey,
			}
			embeddedTLSserver.Server.Handler = r
			fmt.Fprintf(os.Stdout, "Server listening HTTPS %s:%s with embedded Cert\n", host, port)
			log.Fatal(embeddedTLSserver.ListenAndServeTLS(host + ":" + port))
			return
		}

		srv := &http.Server{
			Handler: r,
			Addr:    host + ":" + port,
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		fmt.Fprintf(os.Stdout, "Server listening HTTPS %s:%s with cert from %s\n", host, port, serverCert)
		log.Fatal(srv.ListenAndServeTLS(serverCert, serverkey))
		return
	}

	srv := &http.Server{
		Handler: r,
		Addr:    host + ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Fprintf(os.Stdout, "Server listening HTTP %s:%s\n", host, port)
	log.Fatal(srv.ListenAndServe())
}
