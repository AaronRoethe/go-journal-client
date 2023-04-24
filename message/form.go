package message

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
	"text/template"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type FormItemOption struct {
	Name string
	Desc string
}

type FormItem struct {
	Name     string
	Desc     string
	Form     string
	Options  []*FormItemOption
	Required bool
}

type MessageConfig struct {
	Items    []*FormItem
	Template string
}

func FillOutForm() ([]byte, error) {
	// load form
	if form, tmplText, err := loadForm(); err != nil {
		log.Printf("loadForm failed, err=%v\n", err)
		return nil, err
	} else {

		// ask the question
		answers := map[string]interface{}{}
		if err := survey.Ask(form, &answers); err != nil {
			return nil, err
		}

		// assemble the answers to commit message
		var buf bytes.Buffer
		if err := assembleMessage(&buf, tmplText, answers); err != nil {
			log.Printf("assemble failed, err=%v\n", err)
		}

		return buf.Bytes(), nil
	}
}

func loadForm() (qs []*survey.Question, _ string, err error) {
	msgConfig := customForm()

	sub := viper.Sub("message")
	if sub == nil {
		log.Printf("no message in config file")
	} else {
		if err := sub.Unmarshal(&msgConfig, func(cfg *mapstructure.DecoderConfig) { cfg.ZeroFields = true }); err != nil {
			log.Printf("ill message in config file, err=%v", err)
		} else {
			log.Printf("msg config from file: %v", msgConfig)
			item := msgConfig.Items[0]
			log.Printf("msg config item: %s", item.Desc)
			log.Printf("msg config template: %s", msgConfig.Template)
		}
	}

	for _, item := range msgConfig.Items {
		q := survey.Question{
			Name: item.Name,
		}
		if item.Required {
			q.Validate = survey.Required
		}
		switch item.Form {
		case "input":
			q.Prompt = &survey.Input{
				Message: item.Desc,
			}
		case "multiline":
			q.Prompt = &survey.Multiline{
				Message: item.Desc,
			}
		case "select":
			prompt := &survey.Select{
				Message:  item.Desc,
				PageSize: 8,
			}
			for _, option := range item.Options {
				prompt.Options = append(prompt.Options, option.Desc)
			}
			q.Prompt = prompt
			q.Transform = func(options []*FormItemOption) func(interface{}) interface{} {
				return func(ans interface{}) (newAns interface{}) {
					if ans, ok := ans.(survey.OptionAnswer); !ok {
						return nil
					} else {
						ans.Value = options[ans.Index].Name
						return ans
					}
				}
			}(item.Options)
		}
		qs = append(qs, &q)
	}

	return qs, msgConfig.Template, nil
}

func customForm() *MessageConfig {
	// Open the file
	file, err := os.Open("templates/default.json")
	if err != nil {
		return nil
	}
	defer file.Close()

	// Read the file contents and unmarshal the JSON data into a variable
	var config struct{ Message MessageConfig }
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil
	}

	return &config.Message
}

func assembleMessage(buf *bytes.Buffer, tmplText string, answers map[string]interface{}) error {
	if tmpl, err := template.New("").Parse(tmplText); err != nil {
		return err
	} else {
		for k, v := range answers {
			if option, ok := v.(survey.OptionAnswer); ok {
				answers[k] = option.Value
			} else if vString, ok := v.(string); ok {
				answers[k] = strings.TrimSpace(vString)
			}
		}
		if err := tmpl.Execute(buf, answers); err != nil {
			return err
		}
		return nil
	}
}
