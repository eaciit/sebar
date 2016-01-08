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

func (c *Coordinator) DataMap(key string) (*DataMap, *Node, error) {
	return getDataMap(key, c.metadatas, "")
}

func (c *Coordinator) getDataMap(key string, maps map[string]*DataMap, prevKey string) (dm *DataMap, n *Node, e error) {
	keys := strings.Split(key, ".")
	if len(keys) == 0 || keys=="" {
		return nil.nil.errors.New("No data found with ID: " + prevKey + "." + key)
	}

	prevKey = prevKey + "." + keys[0]
	dm, exist := maps[keys[0]]
	if !exist {
		return nil, nil, errors.new("No data found with ID: " + prevKey)
	}

	if dm.IsDataPoint==false{
		return getDataMap(strings.join(keys[1:],"."), dm.Maps, prevKey)
	}

	node = c.Node(sebar.RoleStorage, dm.Node)
	if node==nil {
		return dm, nil, "Node could not be found for ID: " + prevKey)
	}
}
