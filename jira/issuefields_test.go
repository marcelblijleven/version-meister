package jira_test

import (
	"github.com/marcelblijleven/version-meister/jira"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIssueFieldsToJSONConversion(t *testing.T) {
	issueFields := jira.IssueFields{}
	version, _ := jira.NewVersion("Test version", false, "2019-07-06", 1337)
	project := &jira.Project{
		ID:          "1111",
		Self:        "https://url.to.self/1111",
		Key:         "AB",
		Description: "Test project",
	}

	issueFields.FixVersions = []jira.Version{*version}
	issueFields.Project = *project

	expected := "{\"project\":{\"id\":\"1111\",\"self\":\"https://url.to.self/1111\"," +
		"\"key\":\"AB\",\"description\":\"Test project\"},\"fixVersions\":[{\"name\":\"Test version\"," +
		"\"released\":false,\"releaseDate\":\"2019-07-06\",\"projectId\":1337}]}"
	jsonBytes, err := json.Marshal(issueFields)
	result := string(jsonBytes)

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
