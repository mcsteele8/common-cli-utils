package examples

import "github.com/mcsteele8/common-cli-utils/terminal"

func exampleTerminal() {
	terminal.RunCommand("ls -l", &terminal.RunCmdOptions{
		ShowOutput: true,
	})
}
