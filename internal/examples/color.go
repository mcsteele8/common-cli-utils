package examples

import (
	"fmt"

	"github.com/mcsteele8/common-cli-utils/color"
)

func exampleColor() {
	fmt.Println(color.Red.Paint("This text is Red."))
	fmt.Println(color.YellowBold.Paint("This text is Bold Yellow."))
	fmt.Printf("%s %s & %s\n", color.Red.Paint("red"), color.White.Paint("white"), color.Blue.Paint("blue"))
}
