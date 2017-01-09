package resource

import (
	"bytes"
	"crypto/tls"
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
	client   *http.Client
}

func NewStashClient(host, username, password string, skipSSLVerification bool) *StashClient {
	host = strings.TrimSuffix(host, "/")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSLVerification},
	}
	client := &http.Client{Transport: tr}

	return &StashClient{
		host:     host + "/rest",
		username: username,
		password: password,
		client:   client,
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

	resp, err := c.client.Do(req)
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
