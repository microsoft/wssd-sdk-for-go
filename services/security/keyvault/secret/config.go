// Copyright 2019 (c) Microsoft and contributors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package secret

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"encoding/json"

	log "k8s.io/klog"

	"github.com/jmespath/go-jmespath"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	wssdkeyvault "github.com/microsoft/wssdagent/rpc/security"
)

// Load the secret configuration from the specified path
func LoadConfig(path string) (*keyvault.Secret, error) {
	log.Infof("[LoadConfig] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	srt := &keyvault.Secret{}

	err = yaml.Unmarshal(contents, srt)
	if err != nil {
		return nil, err
	}
	
	return srt, nil
}

// Load the secret configuration from the specified path
func LoadValue(path string) (*string, error) {
	log.Infof("[LoadValue] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	value := string(contents)

	return &value, nil
}

func ExportList(srtList *[]keyvault.Secret, path string, query string, outputType string) error {
	log.Infof("[ExportList] [%s]", path)
	var fileToWrite string 
	if srtList != nil {
		for _, srt := range *srtList {
			str, err := marshalOutput(&srt, query, outputType)
			if err != nil {
				fmt.Printf("%v", err)
				return err
			}
			fileToWrite += string(str)
		}
	}
	err := ioutil.WriteFile(
		path,
		[]byte(fileToWrite),
		0644)
	return err
}

func Print(srt *keyvault.Secret, query string, outputType string) {
	marshaledByte, err := marshalOutput(srt, query, outputType)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(marshaledByte))
}

func PrintList(srtList *[]keyvault.Secret, query string,  outputType string) {
	if srtList != nil {
		for _, srt := range *srtList {
			Print(&srt, query,  outputType)
		}
	}
}

func PrintWssd(srt *wssdkeyvault.Secret) {
	str, err := yaml.Marshal(srt)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

func PrintListWssd(srtList []*wssdkeyvault.Secret) {
	for _, srt := range srtList {
		PrintWssd(srt)
	}
}


func marshalOutput(srt *keyvault.Secret, query string, outputType string) ([]byte, error) {
	var queryTarget interface{}
	var result interface{}
	var err error

	jsonByte, err := json.Marshal(&srt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(jsonByte, &queryTarget)

	if query != "" {	
		result, err = jmespath.Search(query, queryTarget)
		if err != nil {
			return nil, err
		}
	} else {
		result = queryTarget
	}

	var marshaledByte []byte
	if (outputType == "json") {
		marshaledByte, err = json.Marshal(result)
	} else if (outputType == "tsv") {
		marshaledByte, err = marshalTSV(result)
	} else {
		marshaledByte, err = yaml.Marshal(result)
	}

	if err != nil {
		return nil, err
	}

	return marshaledByte, nil
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
