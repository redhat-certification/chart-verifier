package report

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBuilder(t *testing.T) {

	var commands []string
	commands = append(commands, AllCommandsName)
	commands = append(commands, DigestsCommandName)
	commands = append(commands, ResultsCommandName)
	commands = append(commands, MetadataCommandName)
	commands = append(commands, AnnotationsCommandName)

	allCommands := ReportCommandRegistry().AllCommands()

	assert.Equal(t, len(commands), len(allCommands), "Number of commands expected/found differs")

	for _, command := range commands {
		found := false
		for commandname, _ := range allCommands {
			if strings.Compare(commandname, command) == 0 {
				found = true
				break
			}
		}
		assert.True(t, found, "Command not found: %s", command)
	}

}
