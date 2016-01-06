package sebar

import (
	"github.com/eaciit/toolkit"
)

func ParseKey(key string) (owner, table, datakey string) {
	return
}

func MakeKey(owner, datakey, dataindex string) string {
	return owner + "." + datakey + "." + dataindex
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
