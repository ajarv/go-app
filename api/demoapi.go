package api

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"
	"time"
)

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
