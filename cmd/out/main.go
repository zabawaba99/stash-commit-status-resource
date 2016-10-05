package main

import (
	"encoding/json"
	"os"

	"github.com/zabawaba99/stash-commit-status-resource/resource"
)

func main() {
	var req resource.Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		resource.Error("Could not unmarshal request %s", err)
	}

	validStates := []string{"FAILED", "INPROGRESS", "SUCCESSFUL"}
	valid := false
	for _, v := range validStates {
		if v == req.Params.State {
			valid = true
		}
	}
	if !valid {
		resource.Error("Invalid state. State must be one the following %v", validStates)
	}

	if req.Params.Repository == "" {
		resource.Error("You need to specify the repository location")
	}

	if err := resource.Put(req); err != nil {
		resource.Error("Could not update resource %s", err)
	}
}
