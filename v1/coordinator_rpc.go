package sebar

import (
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
	//key := in.GetString("key")
	//owner, table, datakey := ParseKey(key)

	result := toolkit.NewResult()
	result.SetErrorTxt("Set command is stil under development")
	return result
}

func (c *Coordinator) Get(in toolkit.M) *toolkit.Result {
	/*
		key := in.GetString("key")
		owner, table, datakey := ParseKey(key)
		key = MakeKey(owner, table, datakey)
	*/
	result := toolkit.NewResult()
	result.SetErrorTxt("Get command is still under development")
	return result
}

func getSmallestNode(o interface{}) int {
	return 0
}
