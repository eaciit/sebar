package sebar

import (
	//"errors"
	"fmt"
	"github.com/eaciit/appserver/v1"
	"github.com/eaciit/toolkit"
	"strings"
)

type INode interface {
}

type Node struct {
	ID     string
	Role   NodeRoleEnum
	UserID string
	Secret string

	clientRpc *appserver.Client
}

func (n *Node) formatError(pattern string, others ...interface{}) string {
	if len(others) == 0 {
		return "[" + n.ID + "]" + pattern
	}
	return "[" + n.ID + "]" + fmt.Sprintf(pattern, others)
}

func (n *Node) Call(methodName string, in toolkit.M) *toolkit.Result {
	if n.clientRpc == nil {
		n.clientRpc = new(appserver.Client)
		e := n.clientRpc.Connect(n.ID, n.UserID, n.Secret)
		if e != nil {
			n.clientRpc.Close()
			n.clientRpc = nil
			return toolkit.NewResult().SetErrorTxt("Unable to connect to " + n.ID + " : " + e.Error())
		}
		//return toolkit.NewResult().SetErrorTxt(n.formatError("RPC Client is not yet initialized"))
	}

	methodName = strings.ToLower(methodName)
	r := n.clientRpc.Call(methodName, in)
	return r
}

func (n *Node) InitRPC() error {
	/*
		var e error
		n.clientRpc = new(appserver.Client)
		e = n.clientRpc.Connect(n.ID, n.UserID, n.Secret)
		if e != nil {
			n.clientRpc.Close()
			n.clientRpc = nil
			return errors.New("Unable to connect " + n.ID + ": " + e.Error())
		}
		//n.clientRpc.Close()
	*/
	return nil
}
