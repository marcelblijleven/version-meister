package jira_test

import (
	"github.com/marcelblijleven/version-meister/jira"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewComment(t *testing.T) {
	testMessage := "This is a test"
	comment, err := jira.NewComment(testMessage)

	assert.Equal(t, testMessage, comment.Body)
	assert.Nil(t, err)
}

func TestNewCommentReturnsErrorWithEmptyMessage(t *testing.T) {
	testMessage := ""
	comment, err := jira.NewComment(testMessage)

	assert.Equal(t, errors.New("Message cannot be empty"), err)
	assert.Nil(t, comment)
}

func TestCommentToJSONConversion(t *testing.T) {
	expected := "{\"body\":\"This is a test\"}"

	testMessage := "This is a test"
	comment, err := jira.NewComment(testMessage)
	jsonBytes, err := json.Marshal(comment)
	result := string(jsonBytes)

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
