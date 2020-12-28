package tpl_test

import (
	"bytes"
	"encoding/json"
	"log"
)

var td  = testingData(5000)

type test struct {
	description string `json:"description"`
	price       string `json:"price"`
}

func testingData(lim int) map[string]interface{} {

	testSlice := make([]test, 1)

	for i := 0; i < lim; i++ {
		testSlice = append(testSlice, test{"test description", "test price"})
	}


	testMap :=  map[string]interface{}{
		"foo":   "bar",
		"items": testSlice,
	}

	b, err := json.Marshal(testMap)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]interface{})
	if err := json.NewDecoder(bytes.NewBuffer(b)).Decode(&m); err != nil {
		log.Fatal(err)
	}

	return m
}


