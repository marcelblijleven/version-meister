package jira_test

import (
	"github.com/marcelblijleven/version-meister/jira"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewVersion(t *testing.T) {
	version, err := jira.NewVersion("Test version", false, "2019-07-06", 1337)
	expectedVersion := &jira.Version{
		Name:        "Test version",
		Released:    false,
		ReleaseDate: "2019-07-06",
		ProjectID:   1337,
	}

	assert.EqualValues(t, expectedVersion, version)
	assert.Nil(t, err)
}

func TestNewVersionEmptyNameReturnsError(t *testing.T) {
	version, err := jira.NewVersion("", false, "2019-07-06", 1337)
	assert.Equal(t, errors.New("Name cannot be empty"), err)
	assert.Nil(t, version)
}

func TestNewVersionInvalidDateReturnsError(t *testing.T) {
	version, err := jira.NewVersion("Test version", false, "06-07-2019", 1337)
	assert.Equal(t, fmt.Errorf("Received incorrect date string. Expected %v, got %v", "2006-01-02", "06-07-2019"), err)
	assert.Nil(t, version)
}

func TestNewVersionInvalidProjectIDReturnsError(t *testing.T) {
	version, err := jira.NewVersion("Test version", false, "2019-07-06", 0)
	assert.Equal(t, errors.New("ProjectID cannot be 0"), err)
	assert.Nil(t, version)
}

func TestVersionToJSONConversion(t *testing.T) {
	expected := "{\"name\":\"Test version\",\"released\":false,\"releaseDate\":\"2019-07-06\",\"projectId\":1337}"
	version, err := jira.NewVersion("Test version", false, "2019-07-06", 1337)
	jsonBytes, err := json.Marshal(version)
	result := string(jsonBytes)

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
