package color

import (
	"runtime"
	"testing"
)

func Test_color_toString(t *testing.T) {
	tests := []struct {
		name string
		c    color
		want string
	}{
		{"Red Color", Red, "\033[31m"},
		{"Green Color", Green, "\033[32m"},
		{"Yellow Color", Yellow, "\033[33m"},
		{"Blue Color", Blue, "\033[34m"},
		{"Purple Color", Purple, "\033[35m"},
		{"Cyan Color", Cyan, "\033[36m"},
		{"Gray Color", Gray, "\033[37m"},
		{"White Color", White, "\033[97m"},
		{"Black Bold", BlackBold, "\033[1;30m"},
		{"Red Bold", RedBold, "\033[1;31m"},
		{"Green Bold", GreenBold, "\033[1;32m"},
		{"Yellow Bold", YellowBold, "\033[1;33m"},
		{"Blue Bold", BlueBold, "\033[1;34m"},
		{"Purple Bold", PurpleBold, "\033[1;35m"},
		{"Cyan Bold", CyanBold, "\033[1;36m"},
		{"Gray Bold", GrayBold, "\033[1;37m"},
		{"White Bold", WhiteBold, "\033[1;97m"},
		{"Gray Dim", GrayDim, "\033[2m"},
		{"Black", Black, "\033[30m"},
		{"Reset", Reset, "\033[0m"},
		{"Invalid Color", color(999), "\033[0m"}, // Test for invalid input
	}

	// If on Windows, expected output should always be an empty string
	if runtime.GOOS == "windows" {
		for i := range tests {
			tests[i].want = ""
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.toString(); got != tt.want {
				t.Errorf("color.toString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_color_Paint(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		c    color
		args args
		want string
	}{
		{"Paint with Red", Red, args{"Hello"}, "\033[31mHello\033[0m"},
		{"Paint with Green", Green, args{"World"}, "\033[32mWorld\033[0m"},
		{"Paint with Blue", Blue, args{"Test"}, "\033[34mTest\033[0m"},
		{"Paint with Purple", Purple, args{"Example"}, "\033[35mExample\033[0m"},
		{"Paint with Reset", Reset, args{"Plain"}, "\033[0mPlain\033[0m"},
		{"Paint with Invalid Color", color(999), args{"Invalid"}, "\033[0mInvalid\033[0m"},
	}

	// If on Windows, the output should not include escape codes
	if runtime.GOOS == "windows" {
		for i := range tests {
			tests[i].want = tests[i].args.text
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Paint(tt.args.text); got != tt.want {
				t.Errorf("color.Paint() = %v, want %v", got, tt.want)
			}
		})
	}
}
