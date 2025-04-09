package examples

import (
	"fmt"

	"github.com/mcsteele8/common-cli-utils/xprompt"
)

func exampleConformationPrompt() {

	if xprompt.ConformationPrompt("Conformation prompt") {
		fmt.Println("you selected 'yes'")
	} else {
		fmt.Println("you selected 'no'")
	}

}
