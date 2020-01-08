package cli

import (
	"flag"
	"os"
)

// ParseCreateCommand uses Args to determine which flags were called
func ParseCreateCommand(args []string) (string, int, string, string, bool) {
	command := flag.NewFlagSet("create", flag.ExitOnError)
	releaseName := command.String("name", "", "Name of the version")
	projectID := command.Int("project", 0, "ID for the JIRA project")
	component := command.String("component", "", "Optional JIRA Component to include in the JQL query")
	date := command.String("date", "", "Optional date string to include as release date. Use format 2006-01-02")
	dryRun := command.Bool("dryRun", false, "Use dry run to preview which issues would be affected")

	command.Parse(args)

	if *releaseName == "" || *projectID == 0 {
		command.PrintDefaults()
		os.Exit(1)
	}

	return *releaseName, *projectID, *component, *date, *dryRun
}
