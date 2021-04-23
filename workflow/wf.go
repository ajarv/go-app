package workflow

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ajarv/go-app/api"
)

type formData struct {
	APIVersion string `json:"apiVersion"`
	Step       int    `json:"step"`
}

func WorkflowHandler(w http.ResponseWriter, r *http.Request) {
	var form formData

	decoder := json.NewDecoder(r.Body)
	data := api.GetDebugData(r)
	defer api.WriteResponse(w, r, data)

	data["Workflow"] = &form

	err := decoder.Decode(&form)
	if err != nil {
		data["warning"] = fmt.Sprintf("Unable to parse request %v", err)
		return
	}

	switch {
	case form.Step == 0:
		form.Step = 1
		form.APIVersion = api.Info.Version
	case form.Step > 0 && form.APIVersion != api.Info.Version:
		data["warning"] = fmt.Sprintf("Protocol Version Mismatch  Client - %v vs   Server - %v ", form.APIVersion, api.Info.Version)
	case form.Step >= 3:
		form.Step = 3
		data["warning"] = "Order already confirmed. no modifications possible"
	default:
		form.Step = form.Step + 1
		if form.Step == 3 {
			data["message"] = "Order confirmed"
		}
	}

}
