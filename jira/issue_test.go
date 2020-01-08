package jira_test

import (
	"github.com/marcelblijleven/version-meister/jira"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIssueToJSONConversion(t *testing.T) {
	issue := &jira.Issue{
		ID:   "1337",
		Self: "https://url.to.self/1337",
		Key:  "AB-1337",
	}
	expected := "{\"id\":\"1337\",\"self\":\"https://url.to.self/1337\",\"key\":\"AB-1337\"}"
	jsonBytes, err := json.Marshal(issue)
	result := string(jsonBytes)

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
