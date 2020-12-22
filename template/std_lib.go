package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const outputFilenameFmt = "compiled-%s.html"

type GoBasic struct {}

func (*GoBasic) Compile(templateURL, jsonDataURL string) (string, error) {

	dataResponse, err := http.Get(jsonDataURL)
	if err != nil {
		return "", err
	}
	defer dataResponse.Body.Close()

	if dataResponse.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Unable to fetch %s status: %d: %s", jsonDataURL, dataResponse.StatusCode, dataResponse.Status))
	}

	dataMap, err := parseJson(dataResponse.Body)
	if err != nil {
		return "", err
	}

	templateResponse, err := http.Get(templateURL)
	if err != nil {
		return "", err
	}
	defer templateResponse.Body.Close()

	if templateResponse.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Unable to fetch %s status: %d: %s", templateURL, templateResponse.StatusCode, templateResponse.Status))
	}

	htmlTemplateBytes, err := ioutil.ReadAll(templateResponse.Body)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("test").Parse(string(htmlTemplateBytes))
	if err != nil {
		return "", err
	}

	outputFilename := fmt.Sprintf(outputFilenameFmt, time.Now().Format(time.RFC3339Nano))
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return "", err
	}

	return outputFilename, tmpl.Execute(outputFile, dataMap)
}

func parseJson(r io.Reader) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
