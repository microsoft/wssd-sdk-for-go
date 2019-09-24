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

package simplevault

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/keyvault"
	wssdkeyvault "github.com/microsoft/wssdagent/rpc/keyvault"
)

// Load the virtual hard disk configuration from the specified path
func LoadConfig(path string) (*keyvault.SimpleVault, error) {
	log.Infof("[LoadConfig] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	kv := &keyvault.SimpleVault{}

	err = yaml.Unmarshal(contents, kv)
	if err != nil {
		return nil, err
	}
	
	return kv, nil
}

func Print(kv *keyvault.SimpleVault) {
	str, err := yaml.Marshal(kv)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

func PrintList(kvList *[]keyvault.SimpleVault) {
	if kvList != nil {
		for _, kv := range *kvList {
			Print(&kv)
		}
	}
}

func PrintWssd(kv *wssdkeyvault.SimpleVault) {
	str, err := yaml.Marshal(kv)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

func PrintListWssd(kvList []*wssdkeyvault.SimpleVault) {
	for _, kv := range kvList {
		PrintWssd(kv)
	}
}
