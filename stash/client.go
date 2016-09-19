package stash

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type Client struct {
	host     string
	username string
	password string
}

func NewClient(host, username, password string) *Client {
	host = strings.TrimSuffix(host, "/")
	return &Client{
		host:     host + "/rest",
		username: username,
		password: password,
	}
}

type Ref struct {
	ID string `json:"ref"`
}

func (c *Client) CommitsSince(project, repo, branch, lastCommit string) ([]Ref, error) {
	query := url.Values{}
	query.Add("limit", "999")
	query.Add("until", branch)
	if lastCommit != "" {
		query.Add("since", lastCommit)
	}
	path := fmt.Sprintf("%s/api/1.0/projects/%s/repos/%s/commits?%s", c.host, project, repo, query.Encode())
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(response))
	}

	var response struct {
		Values commits `json:"values"`
	}
	if json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	sort.Sort(response.Values)
	refs := make([]Ref, len(response.Values))
	for i, v := range response.Values {
		refs[i] = Ref{ID: v.ID}
	}

	return refs, nil
}

type commit struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"authorTimestamp"`
}

type commits []commit

func (d commits) Len() int {
	return len(d)
}

func (d commits) Less(i, j int) bool {
	return d[i].Timestamp < d[j].Timestamp
}

func (d commits) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type Status struct {
	State       string `json:"state"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	DateAdded   int64  `json:"dateAdded,omitempty"`
}

func (c *Client) BuildStatus(commit string) (*Status, error) {
	path := fmt.Sprintf("%s/build-status/1.0/commits/%s", c.host, commit)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(response))
	}
	var response struct {
		Statuses []Status `json:"values"`
	}
	if json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response.Statuses) == 0 {
		return new(Status), nil
	}

	return &response.Statuses[0], nil
}

func (c *Client) SetBuildStatus(commit string, status Status) error {
	body, err := json.Marshal(status)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/build-status/1.0/commits/%s", c.host, commit)
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
