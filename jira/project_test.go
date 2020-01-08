package jira_test

import (
	"github.com/marcelblijleven/version-meister/jira"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProjectToJSONConversion(t *testing.T) {
	project := &jira.Project{
		ID:          "1111",
		Self:        "https://url.to.self/1111",
		Key:         "AB",
		Description: "Test project",
	}
	expected := "{\"id\":\"1111\",\"self\":\"https://url.to.self/1111\",\"key\":\"AB\",\"description\":\"Test project\"}"
	jsonBytes, err := json.Marshal(project)
	result := string(jsonBytes)

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
