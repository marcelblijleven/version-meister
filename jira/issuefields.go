package jira

// IssueFields represent the fields property on the JIRA issue
type IssueFields struct {
	Project     Project   `json:"project"`
	FixVersions []Version `json:"fixVersions"`
}
