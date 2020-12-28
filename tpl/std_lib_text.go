package tpl

import (
	"bytes"
	"text/template"
)

type GoBasicText struct{}

func (*GoBasicText) Compile(templateBytes []byte, data map[string]interface{}) (string, error) {

	tmpl, err := template.New("test").Parse(string(templateBytes))
	if err != nil {
		return "", err
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, data); err != nil {
		return "", err
	}

	return output.String(), nil
}
