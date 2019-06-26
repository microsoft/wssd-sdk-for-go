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

package virtualmachinescaleset

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/compute"
)

// Load the virtual machine configuration from the specified path
func LoadConfig(path string) (*compute.VirtualMachineScaleSet, error) {
	log.Infof("[LoadConfig] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vm := &compute.VirtualMachineScaleSet{}

	err = yaml.Unmarshal(contents, vm)
	if err != nil {
		return nil, err
	}

	return vm, nil
}

// Print - FixMe: Try to pass in interface{}
func PrintList(vmss *[]compute.VirtualMachineScaleSet) {
	if vmss == nil {
		return
	}
	for _, vms := range *vmss {
		str, err := yaml.Marshal(vms)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		fmt.Printf("%s", string(str))
	}
}
