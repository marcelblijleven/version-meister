package api

import (
	"github.com/marcelblijleven/version-meister/jira"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client is the JIRA api client
type Client struct {
	baseURL    *url.URL
	username   string
	password   string
	httpClient *http.Client
}

type errorMessage struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        innerErrorMessage `json:"errors"`
}

type innerErrorMessage struct {
	Name string `json:"name"`
}

// JQLResult represents the response from the ?search requests
type jqlResult struct {
	Issues []jira.Issue `json:"issues,omitempty"`
}

// UpdateHelper allows for easy marshalling of update data
type updateHelper struct {
	FixVersion fixVersionHelper `json:"update,omitempty"`
}

// FixVersionHelper allows for easy marshalling of update data
type fixVersionHelper struct {
	SetContainers []setContainer `json:"fixVersions,omitempty"`
}

// SetContainer allows for easy marshalling of update data
type setContainer struct {
	Sets []updateSet `json:"set,omitempty"`
}

// UpdateSet allows for easy marshalling of update data
type updateSet struct {
	Name string `json:"name,omitempty"`
}

// SetHTTPClient allows for setting a custom httpClient
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

// NewClient returns a client base on the provided values
func NewClient(baseURL, username, password string) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("Base url cannot be empty")
	}
	if username == "" {
		return nil, errors.New("Username cannot be empty")
	}
	if password == "" {
		return nil, errors.New("Password cannot be empty")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := Client{
		baseURL:    parsedURL,
		username:   username,
		password:   password,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	return &client, nil
}

// Search returns a slice of JIRA issues that match the provided JQL query
func (c *Client) Search(jql string) ([]jira.Issue, error) {
	// TODO: handle maxResults, startAt
	query := url.QueryEscape(jql)
	endpoint := fmt.Sprintf("rest/api/latest/search?jql=%s", query)
	endpointURL, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	resolvedURL := c.baseURL.ResolveReference(endpointURL)
	req, err := http.NewRequest("GET", resolvedURL.String(), nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		var result jqlResult
		if err = json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		return result.Issues, nil
	}

	return nil, fmt.Errorf("Search response status is %v", resp.StatusCode)
}

// CreateVersion creates a new JIRA fixVersion based on the provided Version
func (c *Client) CreateVersion(version jira.Version) error {
	endpoint := "rest/api/latest/version"
	endpointURL, err := url.Parse(endpoint)

	if err != nil {
		return err
	}

	resolvedURL := c.baseURL.ResolveReference(endpointURL)

	buffer := new(bytes.Buffer)
	if json.NewEncoder(buffer).Encode(version); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", resolvedURL.String(), buffer)

	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode != http.StatusCreated {
			msg, err := handleErrorMessage(resp)

			if err != nil {
				return fmt.Errorf("CreateVersion response status is %v", resp.StatusCode)
			}

			if msg.Errors.Name == "A version with this name already exists in this project." {
				fmt.Println(msg.Errors.Name, "Using existing version")
				return nil
			}

			return fmt.Errorf(msg.Errors.Name)
		}

	}

	fmt.Println("Successfully created version", version.Name)
	return nil
}

// AddVersionToIssue adds the provided JIRA version to the provided JIRA issue as a fixVersion
func (c *Client) AddVersionToIssue(issue jira.Issue, version jira.Version) error {
	endpoint, err := url.Parse(fmt.Sprintf("rest/api/latest/issue/%s", issue.ID))

	if err != nil {
		return err
	}

	resolvedURL := c.baseURL.ResolveReference(endpoint)
	container := setContainer{
		Sets: []updateSet{updateSet{Name: version.Name}},
	}
	update := updateHelper{FixVersion: fixVersionHelper{
		SetContainers: []setContainer{container},
	}}

	buffer := new(bytes.Buffer)
	if json.NewEncoder(buffer).Encode(update); err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", resolvedURL.String(), buffer)

	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("AddVersion response status is %v", resp.StatusCode)
	}

	fmt.Printf("Successfully added version %v to issue %v", version.Name, issue.Key)
	return nil
}

// AddCommentToIssue adds the provided JIRA comment as a user comment on the provided JIRA issue
func (c *Client) AddCommentToIssue(issue jira.Issue, comment jira.Comment) error {
	endpoint, err := url.Parse(fmt.Sprintf("rest/api/latest/issue/%s/comment", issue.ID))

	if err != nil {
		return err
	}

	resolvedURL := c.baseURL.ResolveReference(endpoint)
	buffer := new(bytes.Buffer)
	if json.NewEncoder(buffer).Encode(comment); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", resolvedURL.String(), buffer)

	if err != nil {
		return err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		msg, err := handleErrorMessage(resp)

		if err != nil {
			return fmt.Errorf("AddComment response status is %v", resp.StatusCode)
		}

		return fmt.Errorf(msg.Errors.Name)
	}

	return nil
}

func handleErrorMessage(resp *http.Response) (errorMessage, error) {
	defer resp.Body.Close()

	var msg errorMessage
	body, err := ioutil.ReadAll(resp.Body)

	if err = json.Unmarshal(body, &msg); err != nil {
		return msg, err
	}

	return msg, nil
}
