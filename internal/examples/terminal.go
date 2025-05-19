package examples

import "github.com/mcsteele8/common-cli-utils/terminal"

func exampleTerminal() {
	terminal.RunCommand("echo 'Hello, World!'", &terminal.RunCmdOptions{
		ShowOutput: true,
	})
}
