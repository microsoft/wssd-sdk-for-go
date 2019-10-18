// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package store

import (
	"io/ioutil"
	log "k8s.io/klog"
	"os"
	"path"
	"reflect"
	"sync"

	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/marshal"
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

func (c *ConfigStore) GetPath() string {
	return c.path
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
				c.Delete(f.Name())
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
		configPath := c.getConfigPath(id)
		os.MkdirAll(configPath, os.ModePerm)
		err := marshal.ToYAMLFile(result, path.Join(configPath, "config"))
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
	restoreData := reflect.New(c.datatype)
	restored := restoreData.Interface()
	err := marshal.FromYAMLFile(path.Join(configPath, "config"), restored)
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

func (c *ConfigStore) ListFilter(filterName, filterValue string) (*[]interface{}, error) {
	if len(filterValue) == 0 {
		return c.List()
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	// Validate if filterName is valid
	field, found := c.datatype.FieldByName(filterName)
	if !found {
		return nil, errors.InvalidFilter
	}
	if field.Type.Kind() != reflect.String {
		return nil, errors.InvalidFilter
	}

	for k := range c.data {
		tmp := c.data[k]
		val := reflect.ValueOf(tmp)
		if val.Kind() == reflect.Ptr {
			fieldValue := reflect.Indirect(val).FieldByName(filterName).Interface()
			if (fieldValue).(string) == filterValue {
				return &[]interface{}{c.data[k]}, nil
			}
		} else if val.Kind() == reflect.Struct {
			fieldValue := (val).FieldByName(filterName).Interface()
			if (fieldValue).(string) == filterValue {
				return &[]interface{}{c.data[k]}, nil
			}
		}

	}

	return nil, errors.NotFound
}

func (c *ConfigStore) ListFilterMany(filterName, filterValue string) (*[]interface{}, error) {
	if len(filterValue) == 0 {
		return c.List()
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	// Validate if filterName is valid
	field, found := c.datatype.FieldByName(filterName)
	if !found {
		return nil, errors.InvalidFilter
	}
	if field.Type.Kind() != reflect.String {
		return nil, errors.InvalidFilter
	}

	filteredVals := []interface{}{}

	for k := range c.data {
		tmp := c.data[k]
		val := reflect.ValueOf(tmp)
		if val.Kind() == reflect.Ptr {
			fieldValue := reflect.Indirect(val).FieldByName(filterName).Interface()
			if (fieldValue).(string) == filterValue {
				filteredVals = append(filteredVals, c.data[k])
			}
		} else if val.Kind() == reflect.Struct {
			fieldValue := (val).FieldByName(filterName).Interface()
			if (fieldValue).(string) == filterValue {
				filteredVals = append(filteredVals, c.data[k])
			}
		}

	}

	return &filteredVals, nil
}

func (c *ConfigStore) Delete(id string) error {
	if len(id) == 0 {
		return errors.Wrap(errors.InvalidInput, "[Store] id should be valid for delete")
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	// Remove all the configuration, ignore the err
	os.RemoveAll(c.getConfigPath(id))
	if _, ok := c.data[id]; ok {
		delete(c.data, id)

		// Delete all the files
		return nil
	}

	return errors.NotFound
}
