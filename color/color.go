package color

import (
	"fmt"
	"runtime"
)

const (
	Reset color = iota
	Red
	Green
	Yellow
	Blue
	Purple
	Cyan
	Gray
	White
	Black
	RedBold
	GreenBold
	YellowBold
	BlueBold
	PurpleBold
	CyanBold
	GrayBold
	WhiteBold
	BlackBold
	GrayDim
)

type color int

func (c color) Paint(text string) string {
	return fmt.Sprintf("%s%s%s", c.toString(), text, Reset.toString())
}

func (c color) toString() string {
	if runtime.GOOS == "windows" {
		return ""
	}
	switch c {
	case Red:
		return "\033[31m"
	case Green:
		return "\033[32m"
	case Yellow:
		return "\033[33m"
	case Blue:
		return "\033[34m"
	case Purple:
		return "\033[35m"
	case Cyan:
		return "\033[36m"
	case Gray:
		return "\033[37m"
	case White:
		return "\033[97m"
	case BlackBold:
		return "\033[1;30m"
	case RedBold:
		return "\033[1;31m"
	case GreenBold:
		return "\033[1;32m"
	case YellowBold:
		return "\033[1;33m"
	case BlueBold:
		return "\033[1;34m"
	case PurpleBold:
		return "\033[1;35m"
	case CyanBold:
		return "\033[1;36m"
	case GrayBold:
		return "\033[1;37m"
	case WhiteBold:
		return "\033[1;97m"
	case GrayDim:
		return "\033[2m"
	case Black:
		return "\033[30m"
	case Reset:
		return "\033[0m"
	default:
		return "\033[0m"
	}
}
