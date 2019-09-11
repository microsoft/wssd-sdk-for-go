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

package virtualharddisk

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"

	log "k8s.io/klog"

	"github.com/microsoft/wssd-sdk-for-go/services/storage"
	wssdstorage "github.com/microsoft/wssdagent/rpc/storage"
)

// Load the virtual hard disk configuration from the specified path
func LoadConfig(path string) (*storage.VirtualHardDisk, error) {
	log.Infof("[LoadConfig] [%s]", path)
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	vhd := &storage.VirtualHardDisk{}

	err = yaml.Unmarshal(contents, vhd)
	if err != nil {
		return nil, err
	}
	
	return vhd, nil
}

func Print(vhd *storage.VirtualHardDisk) {
	str, err := yaml.Marshal(vhd)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

func PrintList(vhdList *[]storage.VirtualHardDisk) {
	if vhdList != nil {
		for _, vhd := range *vhdList {
			Print(&vhd)
		}
	}
}


func PrintWssd(vhd *wssdstorage.VirtualHardDisk) {
	str, err := yaml.Marshal(vhd)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("%s", string(str))
}

func PrintListWssd(vhdList []*wssdstorage.VirtualHardDisk) {
	for _, vhd := range vhdList {
		PrintWssd(vhd)
	}
}
