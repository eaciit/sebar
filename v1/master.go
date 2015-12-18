package sebar

import (
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"strings"
)

type Master struct {
	SebarServer

	sessions map[string]*usersession
	nodes    map[string]*Node
}

func (m *Master) Login(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	var model struct {
		UserId   string
		Password string
	}
	k.GetPayload(&model)
	if model.UserId != "arief" || model.Password != "darmawan" {
		return result.SetErrorTxt("Invalid Authorisation")
	}
	sess := new(usersession)
	sess.UserID = model.UserId
	sess.Secret = toolkit.GenerateRandomString("", 32)
	if _, exist := m.sessions[sess.ID()]; exist {
		return result.SetErrorTxt("Duplicate session created")
	}
	m.sessions[sess.ID()] = sess
	result.Data = sess.Secret
	return result
}

func (m *Master) Logout(k *knot.WebContext) interface{} {
	var model struct {
		UserId string
		Secret string
	}
	k.GetPayload(&model)
	sess := new(usersession)
	sess.UserID = model.UserId
	sess.Secret = model.Secret
	delete(m.sessions, sess.ID())
	return toolkit.NewResult()
}

func (m *Master) validate(validateType, secret, reference1 string) bool {
	ret := false
	validateType = strings.ToLower(validateType)
	if validateType == "user" {
		return true
	} else {
		return true
	}
	return ret
}

func (m *Master) AddNode(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	var model struct {
		URL, Secret string
		Role        string
	}
	k.GetPayload(&model)
	if model.Secret != m.Secret {
		return result.SetErrorTxt("Not authorised")
	}
	node := new(Node)
	node.Role = NodeRoleEnum(model.Role)
	node.Secret = toolkit.GenerateRandomString("", 32)
	node.SetURL(model.URL)
	nodeKey := node.Key()
	if _, exist := m.nodes[nodeKey]; exist {
		return result.SetErrorTxt("Node already exist: " + node.url + " [" + string(node.Role) + "]")
	}
	m.nodes[nodeKey] = node
	result.Data = node.Secret
	return result
}

func (m *Master) RemoveNode(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	var model struct {
		URL, Secret string
		Role        NodeRoleEnum
	}
	k.GetPayload(&model)
	if !m.validate("node", model.Secret, model.URL) {
		return result.SetErrorTxt("Invalid node call authorisation")
	}
	node := Node{}
	node.Role = model.Role
	node.SetURL(model.URL)
	delete(m.nodes, node.Key())
	return result
}

func (m *Master) Write(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	return result
}

func (m *Master) Read(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	return result
}
