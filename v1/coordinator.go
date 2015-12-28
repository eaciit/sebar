package sebar

import (
//"github.com/eaciit/knot/knot.v1"
//"github.com/eaciit/toolkit"
//"strings"
)

type Coordinator struct {
	SebarServer

	sessions map[string]*usersession
	nodes    map[string]INode
}

func (c *Coordinator) AddNode(id string, node INode) {
	if c.nodes == nil {
		c.nodes = map[string]INode{}
	}
	c.nodes[id] = node
}

func (c *Coordinator) RemoveNode(id string) {
	if c.nodes == nil {
		return
	}
	delete(c.nodes, id)
}

/*
func (m *Coordinator) Login(k *knot.WebContext) interface{} {
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

func (m *Coordinator) Logout(k *knot.WebContext) interface{} {
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

func (m *Coordinator) validate(validateType, secret, reference1 string) bool {
	ret := false
	validateType = strings.ToLower(validateType)
	if validateType == "user" {
		return true
	} else {
		return true
	}
	return ret
}

func (m *Coordinator) AddNode(k *knot.WebContext) interface{} {
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

func (m *Coordinator) RemoveNode(k *knot.WebContext) interface{} {
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

func (m *Coordinator) Write(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	type model struct {
		UserID, Secret, Key string
		Data                interface{}
		Mode                int
	}
	k.GetPayload(model)
	if m.validate("user", model.Secret, model.UserID) == false {
		return result.SetErrorTxt("User is not authorised")
	}
	return result
}

func (m *Coordinator) Read(k *knot.WebContext) interface{} {
	result := toolkit.NewResult()
	return result
}
*/
