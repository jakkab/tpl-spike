package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
)

type Source struct {
	TemplateURL string `json:"templateURL"`
	JSONDataURL string `json:"jsonDataURL"`
}

func (s *Source) Compile(w io.Writer) error {

	dataResponse, err := http.Get(s.JSONDataURL)
	if err != nil {
		return err
	}
	defer dataResponse.Body.Close()

	if dataResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Unable to fetch %s status: %d: %s", s.JSONDataURL, dataResponse.StatusCode, dataResponse.Status))
	}

	dataMap, err := parseJson(dataResponse.Body)
	if err != nil {
		return err
	}

	templateResponse, err := http.Get(s.TemplateURL)
	if err != nil {
		return err
	}
	defer templateResponse.Body.Close()

	if templateResponse.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Unable to fetch %s status: %d: %s", s.TemplateURL, templateResponse.StatusCode, templateResponse.Status))
	}

	htmlTemplateBytes, err := ioutil.ReadAll(templateResponse.Body)
	if err != nil {
		return err
	}

	tmpl, err := template.New("test").Parse(string(htmlTemplateBytes))
	if err != nil {
		return err
	}

	return tmpl.Execute(w, dataMap)
}

func parseJson(r io.Reader) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
