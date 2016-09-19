package resource

import (
	"fmt"
	"os"

	"encoding/json"
)

type Source struct {
	Host       string `json:"host"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Project    string `json:"project"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
}

func (s *Source) UnmarshalJSON(data []byte) error {
	type ss Source
	var v ss
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.Host = v.Host
	s.Username = v.Username
	s.Password = v.Password
	s.Project = v.Project
	s.Repository = v.Repository
	s.Branch = v.Branch
	if s.Branch == "" {
		s.Branch = "master"
	}
	return nil
}

type Version struct {
	Ref string `json:"ref"`
}

type Params struct {
	Name        string `json:"name"`
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
	fmt.Fprintf(os.Stdout, format, values...)
}

func Error(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, format, values...)
	os.Exit(1)
}
