package api_test

import (
	"context"
	"github.com/marcelblijleven/version-meister/api"
	"github.com/marcelblijleven/version-meister/jira"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHTTPClient(handler http.Handler) (*http.Client, func()) {
	server := httptest.NewServer(handler)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}

	return client, server.Close
}

const (
	searchResponse = `{
		"issues": [
			{
				"expand": "renderedFields,names,schema,operations,editmeta,changelog,versionedRepresentations",
				"fields": {
					"fixVersions": [
						{
							"archived": false,
							"id": "1",
							"name": "Test-Version",
							"releaseDate": "2020-01-02",
							"released": true,
							"self": "https://fake.url/rest/api/2/version/1"
						}
					],
					"project": {
						"id": "1",
						"key": "AB",
						"name": "Test Project",
						"projectTypeKey": "software",
						"self": "https://fake.url/rest/api/2/project/12718"
					}
				},
				"id": "1337",
				"key": "AB-123",
				"self": "https://fake.url/rest/agile/1.0/issue/1337"
			}			
		]	
	}`
)

const (
	createVersionResponse = `
	{
		"self": "https://fake.url/rest/api/2/version/1",
		"id": "1",
		"name": "Test-version",
		"archived": false,
		"released": false,
		"releaseDate": "2019-07-07",
		"overdue": false,
		"userReleaseDate": "06/Jul/19",
		"projectId": 1337
	}`
)

const (
	addCommentResponse = `
	{
		"self": "https://fake.url/rest/api/2/issue/154987/comment/1",
		"id": "1",
		"author": {
			"self": "https://fake.url/rest/api/2/user?username=username",
			"name": "Test user",
			"key": "username",
			"emailAddress": "me@mail.com",
			"displayName": "username",
			"active": true,
			"timeZone": "Europe/Amsterdam"
		},
		"body": "A fine test message",
		"updateAuthor": {
			"self": "https://fake.url/rest/api/2/user?username=username",
			"name": "Test user",
			"key": "username",
			"emailAddress": "me@mail.com",
			"displayName": "username",
			"active": true,
			"timeZone": "Europe/Amsterdam"
		},
		"created": "2020-01-08T11:15:25.865+0100",
		"updated": "2020-01-08T11:15:25.865+0100"
	}
	`
)

func TestClientSearch(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		writer.Write([]byte(searchResponse))
	})
	httpClient, closeServer := testHTTPClient(handler)
	defer closeServer()

	client, _ := api.NewClient("http://fake.com", "username", "password")
	client.SetHTTPClient(httpClient)

	issues, err := client.Search("fixVersion IS EMPTY")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(issues))
}

func TestCreateVersion(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		writer.WriteHeader(http.StatusCreated) // Set the status code to 201 - Created
		writer.Write([]byte(createVersionResponse))
	})

	httpClient, closeServer := testHTTPClient(handler)
	defer closeServer()

	client, _ := api.NewClient("http://fake.com", "username", "password")
	client.SetHTTPClient(httpClient)

	version := jira.Version{
		Name:        "Test-version",
		Released:    false,
		ReleaseDate: "2019-07-06",
		ProjectID:   1337,
	}
	err := client.CreateVersion(version)

	assert.Nil(t, err)
}

func TestCreateVersionExistingVersionDoesNotReturnError(t *testing.T) {
	errorResponse := `
	{
		"errorMessages": [],
		"errors": {
			"name": "A version with this name already exists in this project."
		}
	}`

	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		writer.WriteHeader(http.StatusBadRequest) // Set the status code to 400 - Bad Request
		writer.Write([]byte(errorResponse))
	})

	httpClient, closeServer := testHTTPClient(handler)
	defer closeServer()

	client, _ := api.NewClient("http://fake.com", "username", "password")
	client.SetHTTPClient(httpClient)

	version := jira.Version{
		Name:        "Test-version",
		Released:    false,
		ReleaseDate: "2019-07-06",
		ProjectID:   1337,
	}
	err := client.CreateVersion(version)

	assert.Nil(t, err)
}

func TestCreateVersionReturnsError(t *testing.T) {
	errorResponse := `
	{
		"errorMessages": [],
		"errors": {
			"name": "Ship is going down!"
		}
	}`

	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		writer.WriteHeader(http.StatusBadRequest) // Set the status code to 400 - Bad Request
		writer.Write([]byte(errorResponse))
	})

	httpClient, closeServer := testHTTPClient(handler)
	defer closeServer()

	client, _ := api.NewClient("http://fake.com", "username", "password")
	client.SetHTTPClient(httpClient)

	version := jira.Version{
		Name:        "Test-version",
		Released:    false,
		ReleaseDate: "2019-07-06",
		ProjectID:   1337,
	}
	err := client.CreateVersion(version)

	assert.NotNil(t, err)
	assert.Equal(t, "Ship is going down!", err.Error())
}

func TestAddVersionToIssue(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		writer.WriteHeader(http.StatusNoContent) // Set the status code to 204 - No content
	})

	httpClient, closeServer := testHTTPClient(handler)
	defer closeServer()

	client, _ := api.NewClient("http://fake.com", "username", "password")
	client.SetHTTPClient(httpClient)

	version := jira.Version{
		Name:        "Test-version",
		Released:    false,
		ReleaseDate: "2019-07-06",
		ProjectID:   1337,
	}
	issue := jira.Issue{
		ID:      "1",
		Key:     "AB-124",
		Self:    "https://fake.url/rest/api/latest/issue/1",
		Summary: "A fine test issue",
	}

	err := client.AddVersionToIssue(issue, version)

	assert.Nil(t, err)
}

func TestAddCommentToIssue(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		writer.WriteHeader(http.StatusCreated) // Set the status code to 201 - Created
		writer.Write([]byte(addCommentResponse))
	})

	httpClient, closeServer := testHTTPClient(handler)
	defer closeServer()

	client, _ := api.NewClient("http://fake.com", "username", "password")
	client.SetHTTPClient(httpClient)

	issue := jira.Issue{
		ID:      "1",
		Key:     "AB-124",
		Self:    "https://fake.url/rest/api/latest/issue/1",
		Summary: "A fine test issue",
	}
	comment := jira.Comment{
		Body: "A fine test message",
	}

	err := client.AddCommentToIssue(issue, comment)

	assert.Nil(t, err)
}

func TestAddCommentToIssueMissingBody(t *testing.T) {
	errorResponse := `
	{
		"errorMessages": [],
		"errors": {
			"name": "Comment body can not be empty!"
		}
	}`

	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		username, password, ok := req.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "username", username)
		assert.Equal(t, "password", password)
		writer.WriteHeader(http.StatusBadRequest) // Set the status code to 400 - Bad Request
		writer.Write([]byte(errorResponse))
	})

	httpClient, closeServer := testHTTPClient(handler)
	defer closeServer()

	client, _ := api.NewClient("http://fake.com", "username", "password")
	client.SetHTTPClient(httpClient)

	issue := jira.Issue{
		ID:      "1",
		Key:     "AB-124",
		Self:    "https://fake.url/rest/api/latest/issue/1",
		Summary: "A fine test issue",
	}

	comment := jira.Comment{} // Not setting body

	err := client.AddCommentToIssue(issue, comment)

	assert.NotNil(t, err)
	assert.Equal(t, "Comment body can not be empty!", err.Error())
}
