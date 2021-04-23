package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var tmpl *template.Template
var tmplErr bool

func init() {
	tmpl = template.Must(template.ParseFiles("templates/layout.html"))
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("run time panic: %v\n", err)
			tmplErr = true
		}
	}()
}

func WriteResponse(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
	if len(r.Header["Accept"]) > 0 {
		if strings.Contains(r.Header["Accept"][0], "html") && !tmplErr {
			w.Header().Set("Content-Type", "text/html")

			err := tmpl.Execute(w, data)
			if err != nil {
				w.Write([]byte(`{"result":"Error"}`))
			}
			return
		} else if strings.Contains(r.Header["Accept"][0], "yaml") {
			w.Header().Set("Content-Type", "application/yaml")
			b, err := yaml.Marshal(&data)
			if err != nil {
				w.Write([]byte(`{"result":"Error"}`))
				return
			}

			w.Write(b)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(&data, "", "    ")
	if err != nil {
		w.Write([]byte(`{"result":"Error"}`))
		return
	}
	w.Write(b)
	return

}
