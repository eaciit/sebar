package sebar

import (
	"errors"
	"strings"
)

type DataMap struct {
	ID          string
	IsDataPoint bool
	Maps        map[string]*DataMap
	Node        int
}

func NewDataMap(id string) *DataMap {
	dm := new(DataMap)
	dm.Maps = map[string]*DataMap{}
	return dm
}

func (c *Coordinator) DataMap(key string) (*DataMap, error) {
	dm, e := c.findOrCreateMap(key, nil, false)
	return dm, e
}

func (c *Coordinator) findOrCreateMap(key string, maps map[string]*DataMap, createIfNil bool) (dm *DataMap, e error) {
	if key == "" {
		e = errors.New("Key is empty")
		return
	}

	if maps == nil {
		maps = c.metadatas
	}

	keys := strings.Split(key, ".")
	dm, exist := maps[keys[0]]
	if !exist && createIfNil {
		dm = NewDataMap(keys[0])
		maps[keys[0]] = dm
	} else if !exist && !createIfNil {
		e = errors.New("Data with ID " + key + " is not exist")
	}
	if len(keys) > 1 {
		dm, e = c.findOrCreateMap(strings.Join(keys[1:], "."), dm.Maps, createIfNil)
	}
	return
}
