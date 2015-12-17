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
	node.nodeRole = NodeRoleEnum(model.Role)
	node.SetURL(model.URL)
	return result
}
