//
// data.go
//
// Created by Frederic DELBOS - fred@hyperboloide.com on Feb  8 2015.
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of this source code package.
//

package sprocess

import (
	"encoding/json"
	"errors"
	"sync"
)

type Data struct {
	sync.RWMutex
	data map[string]interface{}
}

func NewData() *Data {
	return &Data{
		data: make(map[string]interface{}),
	}
}

func NewDataFrom(o map[string]interface{}) *Data {
	return &Data{
		data: o,
	}
}

func (d *Data) Get(key string) (interface{}, error) {
	d.RLock()
	defer d.RUnlock()
	v, exists := d.data[key]
	if exists == false {
		return nil, errors.New("Data '" + key + "' not found")
	}
	return v, nil
}

func (d *Data) Export() map[string]interface{} {
	d.RLock()
	defer d.RUnlock()

	copy := make(map[string]interface{})
	for k, v := range d.data {
		copy[k] = v
	}
	return copy
}

func (d *Data) Set(key string, value interface{}) {
	d.Lock()
	defer d.Unlock()
	d.data[key] = value
}

func (d *Data) Filter() ([]byte, error) {
	d.RLock()
	defer d.RUnlock()

	copy := map[string]interface{}{
		"size":       d.data["size"],
		"identifier": d.data["identifier"],
		"filename":   d.data["filename"],
	}
	return json.Marshal(copy)
}

func (d *Data) ToJson() ([]byte, error) {
	d.RLock()
	defer d.RUnlock()
	return json.Marshal(d.data)
}
