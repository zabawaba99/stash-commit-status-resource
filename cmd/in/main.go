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

	if req.Version.Ref == "" {
		req.Version.Ref = "1"
	}

	resp := resource.Response{
		Version: req.Version,
	}
	if err := resource.Output(resp); err != nil {
		resource.Error("failed to output body %#v with error %s", req, err)
	}
}
