package sebar

import (
	//"github.com/eaciit/knot/knot.v1"
	"errors"
	"github.com/eaciit/toolkit"
	//"strings"
	//"fmt"
)

type Coordinator struct {
	SebarServer

	sessions map[string]*usersession
}

func (c *Coordinator) RequestJoin(in toolkit.M) *toolkit.Result {
	var e error
	r := toolkit.NewResult()
	referenceID := in.GetString("auth_referenceid")
	//toolkit.Println("User ID Request Join : " + referenceID)
	secret := in.GetString("auth_secret")
	nodeid := in.GetString("nodeid")
	noderole := NodeRoleEnum(in.GetString("noderole"))

	//--- init  node
	node := new(Node)
	node.ID = nodeid
	node.Role = noderole
	node.UserID = referenceID
	node.Secret = secret
	//node.InitRPC()
	e = c.AddNode(node)
	if e != nil {
		r.SetErrorTxt(e.Error())
	}
	//fmt.Printf("Nodes now:\n%s\n", toolkit.JsonString(c.nodes))

	r.Data = node.Secret
	return r
}

func (c *Coordinator) Start() error {
	errorPrefix := "Starting coordinator server fail: "
	c.SebarServer.Server.RegisterRPCFunctions(c)
	c.SebarServer.Server.Fn("requestjoin").AuthType = ""

	e := c.SebarServer.Start()
	if e != nil {
		return errors.New(errorPrefix + e.Error())
	}
	return nil
}

func (c *Coordinator) Stop() error {
	var e error
	es := []string{}
	errorPrefix := "Stop server " + c.Address + " fail: "

	for _, nodes := range c.nodes {
		for _, node := range nodes {
			//toolkit.Println("Calling stopserver from node " + node.ID)
			r := node.Call("stopserver", nil)
			if r.Status != toolkit.Status_OK {
				es = append(es, "["+node.ID+"] "+r.Message)
			}
		}
	}

	if len(es) > 0 {
		return errors.New(errorPrefix + "\n" + func() string {
			s := ""
			for _, e := range es {
				s += e + "\n"
			}
			return s
		}())
	}

	e = c.SebarServer.Stop()
	if e != nil {
		return errors.New(errorPrefix + e.Error())
	}
	return nil
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
