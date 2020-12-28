package tpl

import (
	hb "github.com/aymerick/raymond"
)

type Handlebars struct{}

func (*Handlebars) Compile(templateBytes []byte, data map[string]interface{}) (string, error) {

	tmpl, err := hb.Parse(string(templateBytes))
	if err != nil {
		return "", err
	}

	return tmpl.Exec(data)
}
