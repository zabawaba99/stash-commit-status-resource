package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func main() {
	var out struct {
		Version interface{} `json:"version"`
	}

	inBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(inBytes, &out); err != nil {
		panic(err)
	}
	if out.Version == nil {
		os.Stderr.WriteString("missing version")
		os.Exit(1)
	}
	outBytes, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	os.Stdout.WriteString(string(outBytes))
}
