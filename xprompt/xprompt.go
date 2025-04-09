package xprompt

import (
	"errors"
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
	"github.com/mcsteele8/common-cli-utils/color"
)

type PromptOptions struct {
	MaskInput    bool
	DefaultValue string
}

type DropdownPromptOptions struct {
	DefaultValue string
}

type noBellStdoutStruct struct{}

func (n *noBellStdoutStruct) Write(p []byte) (int, error) {
	if len(p) == 1 && p[0] == '\a' {
		return 0, nil
	}
	return readline.Stdout.Write(p)
}

func (n *noBellStdoutStruct) Close() error {
	return readline.Stdout.Close()
}

var noBellStdout = &noBellStdoutStruct{}

func ConformationPrompt(message string) bool {
	ask := promptui.Prompt{
		Label: color.Cyan.Paint(message + " (y/n)"),
		Validate: func(input string) error {
			if input != "y" && input != "Y" && input != "n" && input != "N" {
				return errors.New("needs \"y\" or \"n\" response")
			}
			return nil
		},
	}

	result, err := ask.Run()
	if err != nil {
		log.Fatalf("conformation prompt \"%s\" FAILED with | error: %s", message, err.Error())
		return false
	}

	if result == "y" || result == "Y" {
		return true
	}

	return false
}

func DropdownPrompt(message string, options []string, dropdownOptionsRaw ...DropdownPromptOptions) string {
	var dropdownOptions DropdownPromptOptions
	if len(dropdownOptionsRaw) > 0 {
		dropdownOptions = dropdownOptionsRaw[0]
	} else {
		dropdownOptions = DropdownPromptOptions{}
	}

	var defaultLabel = ""
	if dropdownOptions.DefaultValue != "" {
		for index, key := range options {
			if key == dropdownOptions.DefaultValue {
				defaultLabel = color.Red.Paint("(DEFAULT)") + " " + dropdownOptions.DefaultValue
				options = append(options[:index], options[index+1:]...)
				options = append([]string{defaultLabel}, options...)
			}
		}
	}

	ask := promptui.Select{
		Label: color.Cyan.Paint(message),
		Items: options,
		Size:  len(options),
	}

	_, result, err := ask.Run()
	if err != nil {
		log.Fatalf("dropdown prompt: \"%s\" FAILED with | error: %s", message, err.Error())
		return ""
	}

	if result == defaultLabel {
		return dropdownOptions.DefaultValue
	}

	return result
}

func ValidatePrompt(message, defaultValue string, validate func(input string) error) string {

	ask := promptui.Prompt{
		Label:    color.Cyan.Paint(message),
		Validate: validate,
		Default:  defaultValue,
		Pointer:  promptui.PipeCursor,
		Stdout:   noBellStdout,
	}

	result, err := ask.Run()
	if err != nil {
		log.Fatalf("validate prompt \"%s\" FAILED with | error: %s", message, err.Error())
		return ""
	}

	return result
}

func mapOptionsOntoPrompt(optionsList *[]PromptOptions, prompt *promptui.Prompt) {
	if optionsList != nil && len(*optionsList) > 0 && prompt != nil {
		options := (*optionsList)[0]
		prompt.Default = options.DefaultValue
		if options.MaskInput {
			prompt.Mask = '*'
			prompt.HideEntered = true
		}
	}

	if prompt != nil {
		prompt.Stdout = noBellStdout
	}
}

func Prompt(message string, options ...PromptOptions) string {
	ask := promptui.Prompt{
		Label: color.Cyan.Paint(message),
		Validate: func(input string) error {
			if input == "" {
				return errors.New("needs response")
			}
			return nil
		},
		Pointer: promptui.PipeCursor,
	}

	mapOptionsOntoPrompt(&options, &ask)

	result, err := ask.Run()
	if err != nil {
		log.Fatalf("text prompt \"%s\" FAILED with | error: %s", message, err.Error())
		return ""
	}

	return result
}

func MultiSelect(message string, options []string, defaultSelected ...string) []string {

	results := []string{}
	ask := &survey.MultiSelect{
		Message: color.Cyan.Paint(message),
		Options: options,
		Default: defaultSelected,
	}

	survey.AskOne(ask, &results)

	return results
}
