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

	if err := in(req); err != nil {
		resource.Error("Could not get resource %s", err)
	}
}

type metadataItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func in(req resource.Request) error {
	src := req.Source
	client := stash.NewClient(src.Host, src.Username, src.Password)
	status, err := client.BuildStatus(req.Version.Ref)
	if err != nil {
		return err
	}
	result := resource.Response{
		Version: req.Version,
		Metadata: resource.Metadata{
			{Name: "commit", Value: req.Version.Ref},
			{Name: "date_added", Value: strconv.FormatInt(status.DateAdded, 10)},
			{Name: "description", Value: status.Description},
			{Name: "key", Value: status.Key},
			{Name: "name", Value: status.Name},
			{Name: "state", Value: status.State},
			{Name: "url", Value: status.URL},
		},
	}

	if err := ioutil.WriteFile(os.Args[1]+"/commit", []byte(req.Version.Ref), 0777); err != nil {
		return err
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
