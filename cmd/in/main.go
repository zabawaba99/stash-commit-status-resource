package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/zabawaba99/stash-commit-status-resource/resource"
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
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = req.Params.Repository
	sha, err := cmd.CombinedOutput()
	if err != nil {
		resource.Error("Could not fetch commit sha %s", err)
	}

	req.Version = resource.Version{Ref: string(sha)}
	result := resource.Response{
		Version: req.Version,
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
