# version-meister
Create versions and release issues in JIRA

## Example uses

```go
package main

import (
    "fmt"
	"github.com/marcelblijleven/version-meister/api"
	"github.com/marcelblijleven/version-meister/cli"
    "github.com/marcelblijleven/version-meister/jira"
    "os"
)

// Example JQL queries
const (
    jqlTemplate = "project = %v AND status = \"%v\" and fixVersion Is EMPTY"
)

func main() {
    // Get credentials from env variables
    baseURL := os.Getenv("JIRA_URL")
    username := os.Getenv("JIRA_USERNAME")
    password := os.Getenv("JIRA_PASSWORD")

    // Create a new api client
    client, err := api.NewClient(baseUrl, username, password)

    if err != nil {
        panic(err)
    }

    // Create a version
    version := jira.Version{
        Name:        name,
        Released:    false,
        ReleaseDate: "2020-01-16",
        ProjectID:   1337,
    }

    // Find issues
    jql := fmt.Sprintf(jqlTemplate, 1337, "Ready for Release")
    issues, err := client.Search(jql)

    if err != nil {
        panic(err)
    }

    // Assign version to issues
    for _, issue := range issues {
        client.AddVersionToIssue(issue, version)
    }
}
```