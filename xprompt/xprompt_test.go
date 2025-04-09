package xprompt

import (
	"reflect"
	"testing"

	"github.com/AlecAivazis/survey/v2"
)

var mockAskOne func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error

func TestMultiSelect(t *testing.T) {
	tests := []struct {
		name            string
		message         string
		options         []string
		defaultSelected []string
		mockResponse    []string
		want            []string
	}{

		{
			name:            "No selection",
			message:         "Choose options:",
			options:         []string{"Option 1", "Option 2", "Option 3"},
			defaultSelected: []string{},
			mockResponse:    []string{},
			want:            []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the survey.AskOne response
			mockAskOne = func(_ survey.Prompt, response interface{}, _ ...survey.AskOpt) error {
				res, ok := response.(*[]string)
				if !ok {
					t.Fatalf("response is not of type *[]string")
				}
				*res = tt.mockResponse
				return nil
			}

			got := MultiSelect(tt.message, tt.options, tt.defaultSelected...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MultiSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}
