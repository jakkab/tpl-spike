package main

import (
	"encoding/json"
	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

var (
	tmplFile           = "./sample/template.html"
	jsonFilePath       = "./sample/sample-data.json"
	outputHtmlFilename = "output.html"
	outputPdfFileName = "/media/sf_PLJAKAB/output.pdf"
)

func main() {

	/// Part 1, main microservice

	dataBytes, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}

	dataMap, err := parseJson(dataBytes)
	if err != nil {
		log.Fatalf("unable to convert json: %s", err.Error())
	}

	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		log.Fatalf("unable to parse: %s", err.Error())
	}

	output, err := os.Create(outputHtmlFilename)
	if err != nil {
		log.Fatalf("unable to create file: %s, %s", outputHtmlFilename, err.Error())
	}

	if err = tmpl.Execute(output, dataMap); err != nil {
		log.Fatalf("unable to execute template: %s", err.Error())
	}

	/// Part 2, separate microservice to convert html file to pdf, I guess

	pdfGen, err := pdf.NewPDFGenerator()
	if err != nil {
		log.Fatalf("unable to init pdf generator: %s", err.Error())
	}

	htmlContent, err := os.Open(outputHtmlFilename)
	if err != nil {
		log.Fatal(err)
	}

	pdfGen.AddPage(pdf.NewPageReader(htmlContent))

	err = pdfGen.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfGen.WriteFile(outputPdfFileName)
	if err != nil {
		log.Fatal(err)
	}

	//if err := os.Remove(outputHtmlFilename); err != nil {
	//	fmt.Printf("unable to remove intermediate html file %s", outputHtmlFilename)
	//}
}

func parseJson(b []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
