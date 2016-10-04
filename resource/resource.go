package resource

import (
	"fmt"
	"os"
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

func Error(format string, values ...interface{}) {
	fmt.Fprintf(os.Stderr, format, values...)
	os.Exit(1)
}
