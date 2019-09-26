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

	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	wssdkeyvault "github.com/microsoft/wssdagent/rpc/security"
)

// Load the virtual hard disk configuration from the specified path
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

func ExportList(srtList *[]keyvault.Secret, path string) error {
	log.Infof("[ExportList] [%s]", path)
	var fileToWrite string 
	if srtList != nil {
		for _, srt := range *srtList {
			str, err := yaml.Marshal(srt)
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

func Print(srt *keyvault.Secret) {
	str, err := yaml.Marshal(srt)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

func PrintList(srtList *[]keyvault.Secret) {
	if srtList != nil {
		for _, srt := range *srtList {
			Print(&srt)
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
