package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/zabawaba99/stash-commit-status-resource/resource"
	"github.com/zabawaba99/stash-commit-status-resource/stash"
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

	if err := out(req); err != nil {
		resource.Error("Could not update resource %s", err)
	}
}

func out(req resource.Request) error {
	src := req.Source
	client := stash.NewClient(src.Host, src.Username, src.Password)

	commit, err := ioutil.ReadFile(os.Args[1] + "/" + req.Params.Commit + "/commit")
	if err != nil {
		return err
	}
	status := stash.Status{
		State:       req.Params.State,
		Key:         os.Getenv("BUILD_JOB_NAME"),
		Name:        fmt.Sprintf("%s-%s", os.Getenv("BUILD_JOB_NAME"), os.Getenv("BUILD_ID")),
		Description: req.Params.Description,
		URL:         os.Getenv("ATC_EXTERNAL_URL"),
	}

	if client.SetBuildStatus(string(commit), status); err != nil {
		return err
	}

	version := resource.Version{Ref: string(commit)}
	result := resource.Response{
		Version: version,
		Metadata: resource.Metadata{
			{Name: "commit", Value: version.Ref},
			{Name: "date_added", Value: strconv.FormatInt(status.DateAdded, 10)},
			{Name: "description", Value: status.Description},
			{Name: "key", Value: status.Key},
			{Name: "name", Value: status.Name},
			{Name: "state", Value: status.State},
			{Name: "url", Value: status.URL},
		},
	}

	return outputResult(result)
}

func outputResult(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Printf("%s", b)
	return nil
}
