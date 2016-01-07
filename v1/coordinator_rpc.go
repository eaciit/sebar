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
	//key := in.GetString("key")
	//owner, table, datakey := ParseKey(key)
	result := toolkit.NewResult()

	nodeIdx, e := getAvailableNode(in.Get("data"))
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

func getAvailableNode(o interface{}) (int, error) {
	return 0, errors.New("Coordinator.getAvailableNode: No node is available to receive data")
	return 0, nil
}
