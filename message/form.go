package message

import (
	"encoding/json"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
)

type FormItemOption struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type FormItem struct {
	Name     string            `json:"name"`
	Desc     string            `json:"desc"`
	Form     string            `json:"form"`
	Options  []*FormItemOption `json:"options"`
	Required bool              `json:"required"`
}

type MessageConfig struct {
	Items    []*FormItem `json:"items"`
	Template string      `json:"template"`
}

func journalPrompt(form []*survey.Question) (JournalEntry, error) {
	var journalEntry JournalEntry
	if err := survey.Ask(form, &journalEntry); err != nil {
		return JournalEntry{}, err
	}

	return journalEntry, nil
}

func LoadForm() (qs []*survey.Question, tmplText string, err error) {
	var msgConfig MessageConfig
	err = json.Unmarshal([]byte(defaultConfigJSON), &msgConfig)
	if err != nil {
		return nil, "", fmt.Errorf("could not unmarshal config data: %v", err)
	}

	for _, item := range msgConfig.Items {
		q := survey.Question{Name: item.Name}
		if item.Required {
			q.Validate = survey.Required
		}
		q.Prompt, q.Transform = CreatePromptAndTransform(item)
		qs = append(qs, &q)
	}

	return qs, msgConfig.Template, nil
}

func CreatePromptAndTransform(item *FormItem) (survey.Prompt, func(interface{}) interface{}) {
	switch item.Form {
	case "input":
		return &survey.Input{Message: item.Desc}, nil
	case "multiline":
		return &survey.Multiline{Message: item.Desc}, nil
	case "select":
		prompt := &survey.Select{Message: item.Desc, PageSize: 8}
		for _, option := range item.Options {
			prompt.Options = append(prompt.Options, option.Desc)
		}
		transform := func(options []*FormItemOption) func(interface{}) interface{} {
			return func(ans interface{}) (newAns interface{}) {
				if ans, ok := ans.(survey.OptionAnswer); ok {
					ans.Value = options[ans.Index].Name
					return ans
				}
				return nil
			}
		}(item.Options)
		return prompt, transform
	default:
		return nil, nil
	}
}
