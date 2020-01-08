package jira

// Issue represent a JIRA issue
type Issue struct {
	ID      string       `json:"id,omitempty"`
	Self    string       `json:"self,omitempty"`
	Key     string       `json:"key,omitempty"`
	Summary string       `json:"summary,omitempty"`
	Fields  *IssueFields `json:"fields,omitempty"`
}
