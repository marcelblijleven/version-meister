package jira

import (
	"errors"
	"fmt"
	"time"
)

// Version represents a JIRA fixVersion
type Version struct {
	Name        string `json:"name"`
	Released    bool   `json:"released"`
	ReleaseDate string `json:"releaseDate"`
	ProjectID   int    `json:"projectId"`
}

// NewVersion returns a Version with the provided arguments as values
func NewVersion(name string, released bool, releaseDate string, projectID int) (*Version, error) {
	if name == "" {
		return nil, errors.New("Name cannot be empty")
	}

	layout := "2006-01-02"
	if releaseDate == "" {
		// Use current date
		releaseDate = time.Now().Format(layout)
	} else {
		// Check if date is in valid layout
		_, err := time.Parse(layout, releaseDate)

		if err != nil {
			return nil, fmt.Errorf("Received incorrect date string. Expected %v, got %v", layout, releaseDate)
		}
	}

	if projectID == 0 {
		return nil, errors.New("ProjectID cannot be 0")
	}

	return &Version{
		Name:        name,
		Released:    released,
		ReleaseDate: releaseDate,
		ProjectID:   projectID,
	}, nil
}
