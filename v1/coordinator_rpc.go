package sebar

import (
	"errors"
	"github.com/eaciit/toolkit"
	"strings"
)

func ParseKey(key string) (owner, table, datakey string) {
	keyparts := strings.Split(key, ".")
	if len(keyparts) == 1 {
		owner = "public"
		table = "common"
		datakey = keyparts[0]
	} else if len(keyparts) == 2 {
		owner = "public"
		table = keyparts[0]
		datakey = keyparts[1]
	} else if len(keyparts) == 3 {
		owner = keyparts[0]
		table = keyparts[1]
		datakey = keyparts[2]
	} else if len(keyparts) > 3 {
		owner = keyparts[0]
		table = strings.Join(keyparts[1:len(keyparts)-1], ".")
		datakey = keyparts[len(keyparts)-1]
	}
	return
}

func MakeKey(owner, table, datakey string) string {
	return owner + "." + table + "." + datakey
}

func (c *Coordinator) Set(in toolkit.M) *toolkit.Result {
	var e error
	if in == nil {
		in = toolkit.M{}
	}
	result := toolkit.NewResult()

	key := in.GetString("key")

	if key == "" {
		return result.SetErrorTxt("Key is empty")
	}

	data := toolkit.ToBytes(in.Get("data"), "")
	if len(data) == 0 {
		return result.SetErrorTxt("Data is not valid")
	}

	nodeIdx, e := c.getAvailableNode(data)
	if e != nil {
		result.SetErrorTxt("Coordinator.Set: " + e.Error())
	}
	node := c.Node(RoleStorage, nodeIdx)

	delete(in, "auth_referenceid")
	delete(in, "auth_secret")
	rw := node.Call("write", in)
	result.Data = rw.Data
	result.Status = rw.Status
	if result.Status != toolkit.Status_OK {
		result.SetErrorTxt("Coordinator.Set: " + rw.Message)
	}
	return result
}

func (c *Coordinator) Get(in toolkit.M) *toolkit.Result {
	result := toolkit.NewResult()
	result.SetErrorTxt("Get command is still under development")
	key := in.GetString("key")
	owner, table, datakey := ParseKey(key)
	key = MakeKey(owner, table, datakey)
	return result
}

func (c *Coordinator) getAvailableNode(data []byte) (nodeIndex int, e error) {
	var currentMax int
	found := false
	dataLength := float64(len(data))
	dataLengthInt := len(data)
	nodes := c.Nodes(RoleStorage)
	for k, n := range nodes {
		resultAvail := n.Call("status", toolkit.M{}.Set("size", dataLengthInt))
		if resultAvail.Status == toolkit.Status_OK {
			m := toolkit.M{}
			resultAvail.GetFromBytes(&m)
			nodeAvailableSize := m.GetInt("available")
			if nodeAvailableSize > dataLengthInt && nodeAvailableSize > currentMax {
				found = true
				currentMax = nodeAvailableSize
				nodeIndex = k
			}
		}
	}

	if !found {
		e = errors.New(toolkit.Sprintf("No node available to hosts %s bytes of data", ParseSize(dataLength)))
	}
	return
}
