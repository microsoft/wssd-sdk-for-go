// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package marshal

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func ToJSON(data interface{}) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func FromJSON(jsonString string, object interface{}) error {
	return json.Unmarshal([]byte(jsonString), object)
}

func ToYAML(data interface{}) (string, error) {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
func FromYAMLBytes(yamlData []byte, object interface{}) error {
	return yaml.Unmarshal(yamlData, object)
}

func FromYAMLString(yamlString string, object interface{}) error {
	return FromYAMLBytes([]byte(yamlString), object)
}

func FromYAMLFile(path string, object interface{}) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return FromYAMLBytes(contents, object)
}
