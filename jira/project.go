package jira

// Project struct represents a project in JIRA
type Project struct {
	ID          string `json:"id"`
	Self        string `json:"self"`
	Key         string `json:"key"`
	Description string `json:"description"`
}
