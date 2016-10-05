package resource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type Source struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Version struct {
	Ref string `json:"ref"`
}

type Params struct {
	Repository  string `json:"repository"`
	Commit      string `json:"commit"`
	State       string `json:"state"`
	Description string `json:"description"`
}

type Metadata []MetadataItem

type MetadataItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Response struct {
	Version  Version  `json:"version"`
	Metadata Metadata `json:"metadata"`
}

type Request struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
	Params  Params  `json:"params"`
}

func Log(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, format, values...)
}

func Error(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, format, values...)
	os.Exit(1)
}

func Output(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s", b)
	return nil
}

func Put(req Request) error {
	src := req.Source
	client := NewStashClient(src.Host, src.Username, src.Password)

	cmd := exec.Command("git", "rev-parse", "--short=40", "HEAD")
	cmd.Dir = fmt.Sprintf("%s/%s", os.Args[1], req.Params.Repository)
	commit, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	commit = bytes.TrimSuffix(commit, []byte("\n"))
	Log("Setting build status for %s\n", commit)

	status := Status{
		State:       req.Params.State,
		Key:         os.Getenv("BUILD_JOB_NAME"),
		Name:        fmt.Sprintf("%s-%s", os.Getenv("BUILD_JOB_NAME"), os.Getenv("BUILD_ID")),
		Description: req.Params.Description,
		URL:         os.Getenv("ATC_EXTERNAL_URL"),
	}

	Log("Build status %#v\n", status)
	if err := client.SetBuildStatus(string(commit), status); err != nil {
		return err
	}
	Log("Status set successfully\n")

	version := Version{Ref: string(commit)}
	result := Response{
		Version: version,
		Metadata: Metadata{
			{Name: "commit", Value: version.Ref},
			{Name: "date_added", Value: strconv.FormatInt(status.DateAdded, 10)},
			{Name: "description", Value: status.Description},
			{Name: "key", Value: status.Key},
			{Name: "name", Value: status.Name},
			{Name: "state", Value: status.State},
			{Name: "url", Value: status.URL},
		},
	}

	return Output(result)
}
