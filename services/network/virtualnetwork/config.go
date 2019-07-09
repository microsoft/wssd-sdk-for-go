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

package virtualnetwork

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/network"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
)

// Load the virtual machine configuration from the specified path
func LoadConfig(path string) (*network.VirtualNetwork, error) {
	log.Infof("[LoadConfig] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vnet := &network.VirtualNetwork{}

	err = yaml.Unmarshal(contents, vnet)
	if err != nil {
		return nil, err
	}

	return vnet, nil
}

// Print
func Print(vnet *network.VirtualNetwork) {
	str, err := yaml.Marshal(vnet)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

// PrintList
func PrintList(vnets *[]network.VirtualNetwork) {
	if vnets != nil {
		for _, vnet := range *vnets {
			Print(&vnet)
		}
	}
}

// Print
func PrintWssd(vnet *wssdnetwork.VirtualNetwork) {
	str, err := yaml.Marshal(vnet)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

// PrintList
func PrintListWssd(vnets []*wssdnetwork.VirtualNetwork) {
	if vnets != nil {
		for _, vnet := range vnets {
			PrintWssd(vnet)
		}
	} else {
		fmt.Printf("No vNET to print\n")
	}
}
