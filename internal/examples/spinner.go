package examples

import (
	"time"

	"github.com/mcsteele8/common-cli-utils/spinner"
)

func exampleSpinner() {
	sp := spinner.New(spinner.CharSets[2], 500*time.Millisecond)
	sp.Start()
	time.Sleep(5 * time.Second)
	sp.Stop()
}
