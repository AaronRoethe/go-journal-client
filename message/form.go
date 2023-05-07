package message

import (
	"encoding/json"
	"log"
	"os"

	survey "github.com/AlecAivazis/survey/v2"
)

type formItemOption struct {
	Name string
	Desc string
}

type formItem struct {
	Name     string
	Desc     string
	Form     string
	Options  []*formItemOption
	Required bool
}

type messageConfig struct {
	Items    []*formItem
	Template string
}

func FillOutForm(form []*survey.Question) (map[string]interface{}, error) {
	answers := map[string]interface{}{}
	if err := survey.Ask(form, &answers); err != nil {
		return nil, err
	}

	return answers, nil
}

func LoadForm() (qs []*survey.Question, tmplText string, err error) {
	configFile, err := os.Open("templates/default.json")
	if err != nil {
		log.Printf("could not open config file: %v", err)
		return nil, "", err
	}
	defer configFile.Close()

	var msgConfig messageConfig
	err = json.NewDecoder(configFile).Decode(&msgConfig)
	if err != nil {
		log.Printf("could not decode config file: %v", err)
		return nil, "", err
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

func CreatePromptAndTransform(item *formItem) (survey.Prompt, func(interface{}) interface{}) {
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
		transform := func(options []*formItemOption) func(interface{}) interface{} {
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
