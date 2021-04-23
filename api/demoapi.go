package api

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type AppInfo struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Info    string `json:"info"`
}

var Info AppInfo

func init() {
	Info.Version = getEnv("APP_VERSION", "v1.0.0")
	Info.Name = getEnv("APP_NAME", "GO_WEB")
	Info.Color = getEnv("APP_COLOR", "black")
	Info.Info = getEnv("APP_INFO", "-")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetDebugData(req *http.Request) map[string]interface{} {
	var hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknown host"
	}

	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			if strings.HasPrefix(key, "secret") {
				val = base64.StdEncoding.EncodeToString([]byte(val))
				val = base64.StdEncoding.EncodeToString([]byte(val))
			}
			items[key] = val
		}
		return items
	}

	data := map[string]interface{}{
		"Host":       hostname,
		"ServerTime": time.Now().Format("2006-01-02 15:04:05"),
		"Request":    map[string]interface{}{"Headers": req.Header},
	}
	data["info"] = Info.Info

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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := GetDebugData(r)
	WriteResponse(w, r, data)
}

func ApiResourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data := GetDebugData(r)
	defer WriteResponse(w, r, data)
	data["resource"] = vars
}
