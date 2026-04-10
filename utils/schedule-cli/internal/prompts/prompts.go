package prompts

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/amterp/jsoncolor"
	"github.com/fatih/color"
)

// Style definitions
var (
	questionStyle = color.New(color.FgCyan, color.Bold)
	errorStyle    = color.New(color.FgRed, color.Bold)
	successStyle  = color.New(color.FgGreen, color.Bold)
)

// AskUsername запрашивает имя пользователя
func AskUsername(defaultVal string) string {
	var username string
	prompt := &survey.Input{
		Message: questionStyle.Sprint("Username:"),
		Default: defaultVal,
		Help:    "Enter your username (3-20 chars, letters, digits, underscore)",
	}
	survey.AskOne(prompt, &username, survey.WithValidator(survey.Required))
	return username
}

// AskPassword запрашивает пароль со звездочками
func AskPassword(confirm bool) string {
	var password string
	prompt := &survey.Password{
		Message: questionStyle.Sprint("Password:"),
		Help:    "Password will be hidden",
	}
	survey.AskOne(prompt, &password, survey.WithValidator(survey.Required))

	if confirm {
		var confirmPwd string
		confirmPrompt := &survey.Password{
			Message: questionStyle.Sprint("Confirm password:"),
		}
		survey.AskOne(confirmPrompt, &confirmPwd, survey.WithValidator(survey.Required))
		if password != confirmPwd {
			errorStyle.Println("Passwords do not match")
			return AskPassword(true) // рекурсивный запрос
		}
	}
	return password
}

// AskNewPassword запрашивает старый и новый пароль для смены
func AskNewPassword() (old, new string) {
	oldPrompt := &survey.Password{
		Message: questionStyle.Sprint("Old password:"),
	}
	survey.AskOne(oldPrompt, &old, survey.WithValidator(survey.Required))

	newPrompt := &survey.Password{
		Message: questionStyle.Sprint("New password:"),
	}
	survey.AskOne(newPrompt, &new, survey.WithValidator(survey.Required))

	var confirm string
	confirmPrompt := &survey.Password{
		Message: questionStyle.Sprint("Confirm new password:"),
	}
	survey.AskOne(confirmPrompt, &confirm, survey.WithValidator(survey.Required))

	if new != confirm {
		errorStyle.Println("New passwords do not match")
		return AskNewPassword()
	}
	return old, new
}

// AskString запрашивает произвольную строку с валидацией
func AskString(message string, required bool, validator func(val any) error) string {
	var value string
	prompt := &survey.Input{
		Message: questionStyle.Sprint(message),
	}
	opts := []survey.AskOpt{}
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}
	if validator != nil {
		opts = append(opts, survey.WithValidator(validator))
	}
	survey.AskOne(prompt, &value, opts...)
	return value
}

// AskInt запрашивает целое число
func AskInt(message string, min, max int) int {
	var value int
	prompt := &survey.Input{
		Message: questionStyle.Sprint(message),
	}
	validators := []survey.Validator{
		func(val any) error {
			if str, ok := val.(string); ok {
				// проверка на число
				_, err := fmt.Sscanf(str, "%d", &value)
				if err != nil {
					return errors.New("must be a number")
				}
				if value < min || (max != 0 && value > max) {
					return fmt.Errorf("value must be between %d and %d", min, max)
				}
			}
			return nil
		},
		survey.Required,
	}
	survey.AskOne(prompt, &value, survey.WithValidator(survey.ComposeValidators(validators...)))
	return value
}

// AskSelect запрашивает выбор из списка
func AskSelect(message string, options []string) string {
	var result string
	prompt := &survey.Select{
		Message: questionStyle.Sprint(message),
		Options: options,
	}
	survey.AskOne(prompt, &result)
	return result
}

// AskMultiSelect запрашивает множественный выбор
func AskMultiSelect(message string, options []string) []string {
	var result []string
	prompt := &survey.MultiSelect{
		Message: questionStyle.Sprint(message),
		Options: options,
	}
	survey.AskOne(prompt, &result)
	return result
}

// AskConfirm запрашивает подтверждение (Y/n)
func AskConfirm(message string, def bool) bool {
	var result bool
	prompt := &survey.Confirm{
		Message: questionStyle.Sprint(message),
		Default: def,
	}
	survey.AskOne(prompt, &result)
	return result
}

// ShowSuccess выводит зеленое сообщение об успехе
func ShowSuccess(message string, data any) {
	if data == nil {
		successStyle.Println("✓ " + message)
	} else {
		msg, err := jsoncolor.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Failed to marshal data: ", err)
			fmt.Println("Data: ", data)
		}
		successStyle.Println("✓ " + message + "\n" + string(msg))
	}
}

// ShowError выводит красное сообщение об ошибке
func ShowError(message string) {
	errorStyle.Println("✗ " + message)
}
