// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package store

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	log "k8s.io/klog"
	"os"
	"path"
	"reflect"
	"sync"

	"github.com/microsoft/wssdagent/pkg/wssdagent/errors"
)

type ConfigStore struct {
	datatype reflect.Type
	path     string
	data     map[string]interface{}
	mux      sync.Mutex
}

func NewConfigStore(basepath string, dtype reflect.Type) *ConfigStore {
	cs := &ConfigStore{
		path:     basepath,
		datatype: dtype,
		data:     map[string]interface{}{},
	}
	err := cs.loadStore()
	if err != nil {
		log.Errorf("Error Loading Store [%v]", err)
	}
	return cs

}

func (c *ConfigStore) getConfigPath(id string) string {
	return path.Join(c.path, id)
}

// LoadStore loads the ConfigStore from the file system
func (c *ConfigStore) loadStore() error {
	log.Infof("LoadStore - Path [%s], Type [%s]", c.path, c.datatype)
	files, err := ioutil.ReadDir(c.path)
	if err != nil {
		return err

	}

	for _, f := range files {
		if f.IsDir() {
			err = c.restore(f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil

}

func (c *ConfigStore) Add(id string, val interface{}) error {
	if len(id) == 0 {
		return errors.Wrap(errors.InvalidInput, "[Store] id should be valid for add")
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	// TODO: Check overrite
	c.data[id] = val
	return c.save(id)
}

// Save would persist the data from config store to the file system
func (c *ConfigStore) save(id string) error {
	if len(id) == 0 {
		return errors.Wrap(errors.InvalidInput, "[Store] id should be valid for save")
	}
	if result, ok := c.data[id]; ok {
		enc, err := yaml.Marshal(result)
		if err != nil {
			return err

		}
		configPath := c.getConfigPath(id)
		os.MkdirAll(configPath, os.ModePerm)
		err = ioutil.WriteFile(path.Join(configPath, "config"), enc, 0644)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.NotFound
}

// restore would restore the data back to the config store
func (c *ConfigStore) restore(id string) error {
	configPath := c.getConfigPath(id)
	log.Infof("restore - Path [%s], Type [%s]", configPath, c.datatype)
	// restore config
	data, err := ioutil.ReadFile(path.Join(configPath, "config"))
	if err != nil {
		return err
	}
	restoreData := reflect.New(c.datatype)
	restored := restoreData.Interface()
	err = yaml.Unmarshal(data, restored)
	if err != nil {
		return err
	}
	c.data[id] = restored
	return nil
}

func (c *ConfigStore) Get(id string) (interface{}, error) {
	if len(id) == 0 {
		return nil, errors.Wrap(errors.InvalidInput, "[Store] id should be valid for get")
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	if val, ok := c.data[id]; ok {
		return val, nil
	}
	return nil, errors.NotFound
}

func (c *ConfigStore) List() (*[]interface{}, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	vals := []interface{}{}
	for k := range c.data {
		vals = append(vals, c.data[k])
	}
	return &vals, nil
}

func (c *ConfigStore) Delete(id string) error {
	if len(id) == 0 {
		return errors.Wrap(errors.InvalidInput, "[Store] id should be valid for delete")
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.data[id]; ok {
		delete(c.data, id)

		// Delete all the files
		os.RemoveAll(c.getConfigPath(id))
		return nil
	}
	return errors.NotFound
}
