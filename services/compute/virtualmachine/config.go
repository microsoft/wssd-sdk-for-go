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

package virtualmachine

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	wssdcompute "github.com/microsoft/wssdagent/rpc/compute"
)

// Load the virtual machine configuration from the specified path
func LoadConfig(path string) (*compute.VirtualMachine, error) {
	log.Infof("[LoadConfig] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vm := &compute.VirtualMachine{}

	err = yaml.Unmarshal(contents, vm)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

// Print
func PrintYAML(vm *compute.VirtualMachine) {
	str, err := yaml.Marshal(vm)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

// PrintList
func PrintYAMLList(vms *[]compute.VirtualMachine) {
	if vms != nil {
		for _, vm := range *vms {
			PrintYAML(&vm)
		}
	}
}

// Print
func PrintWssd(vnet *wssdcompute.VirtualMachine) {
	str, err := yaml.Marshal(vnet)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

// PrintList
func PrintListWssd(vms []*wssdcompute.VirtualMachine) {
	if vms != nil {
		for _, vm := range vms {
			PrintWssd(vm)
		}
	} else {
		fmt.Printf("No VM to print\n")
	}
}
