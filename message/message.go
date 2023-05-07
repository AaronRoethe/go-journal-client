package message

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	survey "github.com/AlecAivazis/survey/v2"
)

func AssembleMessage(buf *bytes.Buffer, tmplText string, answers map[string]interface{}) error {
	tmpl, err := template.New("").Parse(tmplText)
	if err != nil {
		return err
	}

	for k, v := range answers {
		switch value := v.(type) {
		case survey.OptionAnswer:
			answers[k] = value.Value
		case string:
			answers[k] = strings.TrimSpace(value)
		}
	}

	return tmpl.Execute(buf, answers)
}

func OutputStringTemplate(answers map[string]interface{}, tmplText string) ([]byte, error) {
	var buf bytes.Buffer
	if err := AssembleMessage(&buf, tmplText, answers); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Journal() {
	form, tmplText, err := LoadForm()
	if err != nil {
		log.Fatalf("failed to fill load form: %v", err)
	}

	answers, err := FillOutForm(form)
	if err != nil {
		log.Fatalf("failed to fill out form: %v", err)
	}

	fmt.Println("Answers:", answers)

	output, err := OutputStringTemplate(answers, tmplText)
	if err != nil {
		log.Fatalf("failed to output string template: %v", err)
	}

	fmt.Println("Output:", string(output))
}
