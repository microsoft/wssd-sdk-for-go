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

package virtualnetworkinterface

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/network"
	wssdnetwork "github.com/microsoft/wssdagent/rpc/network"
)

// Load the virtual machine configuration from the specified path
func LoadConfig(path string) (*network.VirtualNetworkInterface, error) {
	log.Infof("[LoadConfig] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vnetInterface := &network.VirtualNetworkInterface{}

	err = yaml.Unmarshal(contents, vnetInterface)
	if err != nil {
		return nil, err
	}

	return vnetInterface, nil
}

// Print
func Print(vnetInterface *network.VirtualNetworkInterface) {
	str, err := yaml.Marshal(vnetInterface)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

// PrintList
func PrintList(vnets *[]network.VirtualNetworkInterface) {
	if vnets != nil {
		for _, vnetInterface := range *vnets {
			Print(&vnetInterface)
		}
	}
}

// Print
func PrintWssd(vnetInterface *wssdnetwork.VirtualNetworkInterface) {
	str, err := yaml.Marshal(vnetInterface)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

// PrintList
func PrintListWssd(vnets []*wssdnetwork.VirtualNetworkInterface) {
	if vnets != nil {
		for _, vnetInterface := range vnets {
			PrintWssd(vnetInterface)
		}
	} else {
		fmt.Printf("No vNET Interface to print\n")
	}
}
