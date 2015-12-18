package sebar

import (
	"strings"
)

func NewServer(role NodeRoleEnum, url string) IServer {
	if role == RoleMaster {
		m := new(Master)
		m.Address = url
		return m
	} else if role == RoleNode {
		n := new(Node)
		n.Address = url
		return n
	}
	return nil
}

func MakeUrl(serverUrl, role, method string) string {
	return serverUrl + "/" + role + "/" + method
}

func ParseKey(userid string, key string) (string, string, string) {
	keys := strings.Split(key, ":")
	owner := "public"
	table := "common"
	datakey := ""

	if len(keys) >= 3 {
		owner = keys[0]
		table = keys[1]
		datakey = keys[2]
	} else if len(keys) == 2 {
		table = keys[0]
		datakey = keys[1]
	} else {
		datakey = keys[0]
	}
	return owner, table, datakey
}
