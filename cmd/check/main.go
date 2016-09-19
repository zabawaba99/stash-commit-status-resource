package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/zabawaba99/stash-commit-status-resource/resource"
	"github.com/zabawaba99/stash-commit-status-resource/stash"
)

func main() {
	var req resource.Request
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		resource.Error("Could not unmarshal request %s", err)
	}

	if err := check(req); err != nil {
		resource.Error("Could not check resource %s", err)
	}
}

func check(req resource.Request) error {
	src := req.Source
	client := stash.NewClient(src.Host, src.Username, src.Password)
	shas, err := client.CommitsSince(src.Project, src.Repository, src.Branch, req.Version.Ref)
	if err != nil {
		return err
	}

	return outputResult(shas)
}

func outputResult(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Printf("%s", b)
	return nil
}
