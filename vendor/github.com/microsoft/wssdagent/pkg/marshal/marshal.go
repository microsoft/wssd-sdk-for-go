// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package marshal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Duplicate(data interface{}, duplicatedData interface{}) error {
	dataBytes, err := ToJSONBytes(data)
	if err != nil {
		return err
	}
	err = FromJSONBytes(dataBytes, duplicatedData)
	if err != nil {
		return err
	}
	return nil
}
func ToString(data interface{}) string {
	return fmt.Sprintf("%+v", data)
}

func ToJSON(data interface{}) (string, error) {
	jsonBytes, err := ToJSONBytes(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
func ToJSONBytes(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// ToJSONFile writes the data to path in YAML format
func ToJSONFile(data interface{}, path string) error {
	enc, err := ToJSONBytes(data)
	if err != nil {
		return err

	}

	err = ioutil.WriteFile(path, enc, 0644)
	if err != nil {
		return err
	}
	return nil
}

func FromJSON(jsonString string, object interface{}) error {
	return json.Unmarshal([]byte(jsonString), object)
}

func FromJSONBytes(jsonBytes []byte, object interface{}) error {
	return json.Unmarshal(jsonBytes, object)
}

func FromJSONFile(path string, object interface{}) error {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return FromJSONBytes(contents, object)
}

func ToBase64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func FromBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func ToYAML(data interface{}) (string, error) {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}
func ToYAMLBytes(data interface{}) ([]byte, error) {
	return yaml.Marshal(data)
}

// ToYAMLFile writes the data to path in YAML format
func ToYAMLFile(data interface{}, path string) error {
	enc, err := ToYAMLBytes(data)
	if err != nil {
		return err

	}

	err = ioutil.WriteFile(path, enc, 0644)
	if err != nil {
		return err
	}
	return nil
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

func ToTSV(data interface{}) (string, error) {
	jsonBytes, err := ToTSVBytes(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func ToTSVBytes(data interface{}) ([]byte, error) {
	return marshalTSV(data)
}

func marshalTSV(result interface{}) ([]byte, error) {
	var str []byte
	switch v := result.(type) {
	case string:
		str = []byte(v)
	case map[string]interface{}:
		var tabString string
		for _, value := range v {
			typ, ok := value.(string)
			if ok && typ != "" {
				tabString += typ + "\t"
			}
		}
		str = []byte(tabString)
	default:
		return nil, fmt.Errorf("Unsupported Format")
	}
	return str, nil
}

func ToCSV(data interface{}) (string, error) {
	jsonBytes, err := marshalCSV(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func ToCSVBytes(data interface{}) ([]byte, error) {
	return marshalCSV(data)
}

func marshalCSV(result interface{}) ([]byte, error) {
	var str []byte
	switch v := result.(type) {
	case string:
		str = []byte(v)
	case map[string]interface{}:
		var tabString string
		for _, value := range v {
			typ, ok := value.(string)
			if ok && typ != "" {
				tabString += typ + ","
			}
		}
		str = []byte(tabString)
	default:
		return nil, fmt.Errorf("Unsupported Format")
	}
	return str, nil
}
