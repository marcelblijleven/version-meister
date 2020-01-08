package jira

import (
	"errors"
)

// Comment represents a comment on a JIRA issue
type Comment struct {
	Body string `json:"body,omitempty"`
}

// NewComment returns a Comment with the provided message as body
func NewComment(message string) (*Comment, error) {
	if message == "" {
		return nil, errors.New("Message cannot be empty")
	}

	return &Comment{Body: message}, nil
}
