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

/*
func (c *Coordinator) findMap(key string, maps map[string]*DataMap, prevKey string) (dm *DataMap, node *Node, e error) {
	keys := strings.Split(key, ".")
	if len(keys) == 0 || key == "" {
		return nil, nil, errors.New("No data found with ID: " + prevKey + "." + key)
	}

	prevKey = prevKey + "." + keys[0]
	dm, exist := maps[keys[0]]
	if !exist {
		return nil, nil, errors.New("No data found with ID: " + prevKey)
	}

	if dm.IsDataPoint == false {
		return c.findMap(strings.Join(keys[1:], "."), dm.Maps, prevKey)
	}

	node = c.Node(RoleStorage, dm.Node)
	if node == nil {
		return dm, nil, errors.New("Node could not be found for ID: " + prevKey)
	}

	return
}
*/
