package cli_test

import (
	"github.com/marcelblijleven/version-meister/cli"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestParseCreateCommand(t *testing.T) {
	args := []string{"-name", "Test-Version", "-date", "2019-07-06", "-project", "1337"}
	name, projectID, component, date, dryRun := cli.ParseCreateCommand(args)
	assert.Equal(t, "Test-Version", name)
	assert.Equal(t, 1337, projectID)
	assert.Equal(t, "", component)
	assert.Equal(t, "2019-07-06", date)
	assert.False(t, dryRun)
}

func TestParseCreateCommandExitsWithMissingName(t *testing.T) {
	args := []string{"-name", "-date", "2020-01-07", "-project", "1337"}

	if os.Getenv("DETACHED_PARSE_CREATE_COMMAND") == "1" {
		// In subprocess
		cli.ParseCreateCommand(args)
		return
	}

	// Create a command to run as subprocess
	cmd := exec.Command(os.Args[0], "-test.run=TestParseCreateCommandExitsWithMissingName")
	cmd.Env = append(os.Environ(), "DETACHED_PARSE_CREATE_COMMAND=1")
	err := cmd.Run()
	// Cast err as ExitError
	e, ok := err.(*exec.ExitError)

	assert.True(t, ok && !e.Success())
}
