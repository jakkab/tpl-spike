package assets

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Source struct {
	TemplateType string `json:"templateType"`
	TemplateURL  string `json:"templateURL"`
	JSONDataURL  string `json:"jsonDataURL"`
}

func (s *Source) DownloadTemplate() ([]byte, error) {

	templateResponse, err := http.Get(s.TemplateURL)
	if err != nil {
		return nil, err
	}
	defer templateResponse.Body.Close()

	if templateResponse.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unable to fetch %s status: %d: %s", s.TemplateURL, templateResponse.StatusCode, templateResponse.Status))
	}

	return ioutil.ReadAll(templateResponse.Body)
}

func (s *Source) DownloadDataSource() (map[string]interface{}, error) {
	dataResponse, err := http.Get(s.JSONDataURL)
	if err != nil {
		return nil, err
	}
	defer dataResponse.Body.Close()

	if dataResponse.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unable to fetch %s status: %d: %s", s.JSONDataURL, dataResponse.StatusCode, dataResponse.Status))
	}

	return parseJson(dataResponse.Body)
}

func parseJson(r io.Reader) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
