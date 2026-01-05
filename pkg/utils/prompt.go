package utils

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

/*
PromptSelect displays a selection prompt to the user and returns the selected option.

It uses the survey library to provide an interactive terminal UI.

Returns the selected choice.
*/
func PromptSelect(promptMsg string, options []string, defaultOption string) string {
	var selected string
	prompt := &survey.Select{
		Message: promptMsg,
		Options: options,
		Default: defaultOption,
	}

	err := survey.AskOne(prompt, &selected)
	if err != nil {
		fmt.Println("Cancelled:", err)
		return ""
	}
	return selected
}

/*
PromptConfirm displays a yes/no confirmation prompt to the user.

Returns true if the user selects 'yes', and false otherwise.
*/
func PromptConfirm(promptMsg string, defaultVal bool) bool {
	var confirm bool
	prompt := &survey.Confirm{
		Message: promptMsg,
		Default: defaultVal,
	}
	err := survey.AskOne(prompt, &confirm)
	if err != nil {
		return false
	}
	return confirm
}

/*
PromptInput displays one prompt and ask for the user input (e.g. public key).

Returns the content trimmed
*/
func PromptInput(promptMsg string, defaultVal string, required bool) string {
	var result string
	prompt := &survey.Input{
		Message: promptMsg,
		Default: defaultVal,
	}

	var opts []survey.AskOpt
	if required {
		// 添加必填验证器
		opts = append(opts, survey.WithValidator(survey.Required))
	}

	err := survey.AskOne(prompt, &result, opts...)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(result)
}
