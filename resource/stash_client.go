package resource

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type StashClient struct {
	host     string
	username string
	password string
}

func NewStashClient(host, username, password string) *StashClient {
	host = strings.TrimSuffix(host, "/")
	return &StashClient{
		host:     host + "/rest",
		username: username,
		password: password,
	}
}

type Status struct {
	State       string `json:"state"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	DateAdded   int64  `json:"dateAdded,omitempty"`
}

func (c *StashClient) SetBuildStatus(commit string, status Status) error {
	body, err := json.Marshal(status)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/build-status/1.0/commits/%s", c.host, commit)
	Log("Making request... %s\n", path)
	req, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		response, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(response))
	}
	return nil
}
