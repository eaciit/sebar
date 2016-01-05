package sebar

import (
	"errors"
	"github.com/eaciit/appserver/v1"

	"github.com/eaciit/toolkit"
	//"net/rpc"
)

type IServer interface {
	Start() error
	Stop() error
	Validate(authType, referenceID, secret string) error
}

type SebarServer struct {
	appserver.Server
	//Address,
	//Secret   string
	Coordinator       string
	CoordinatorUserID string
	CoordinatorSecret string
	Role              NodeRoleEnum

	nodes map[string]map[string]*Node
	//_rpcAddress string
	//_urlAddress string
}

func (s *SebarServer) Validate(authType, referenceID, secret string) error {
	return errors.New("Method validate is not yet initialized")
}

func (s *SebarServer) Start() error {
	var e error
	/*
		e = s.RegisterRPCFunctions(s)
		if e != nil {
			return errors.New("Unable to register RPC: " + e.Error())
		}
	*/

	e = s.Server.Start(s.Address)
	if e != nil {
		s.Server.Log.Error("Unable to start server + " + s.Address + ": " + e.Error())
		return e
	}
	//--- get coordinator approval
	if s.Role != RoleCoordinator {
		nodeCoordinator := new(Node)
		nodeCoordinator.Role = RoleCoordinator
		nodeCoordinator.ID = s.Coordinator
		nodeCoordinator.UserID = s.CoordinatorUserID
		nodeCoordinator.Secret = s.CoordinatorSecret
		s.AddNode(s.Coordinator, nodeCoordinator)

		rjoin := nodeCoordinator.Call("requestjoin", toolkit.M{}.
			Set("auth_referenceid", s.CoordinatorUserID).
			Set("auth_secret", s.CoordinatorSecret).
			Set("nodeid", s.Address).
			Set("noderole", string(s.Role)))
		if rjoin.Status != toolkit.Status_OK {
			s.Stop()
			return errors.New(s.Address + " Request to Join Fail: " + rjoin.Message)
		}

		mnode := toolkit.M{}
		rjoin.GetFromBytes(&mnode)
		//s.Secret = mnode.GetString("secret")
		nodeCoordinator.UserID = mnode.GetString("referenceid")
		nodeCoordinator.Secret = mnode.GetString("secret")
	}

	s.Server.Log.Info("Starting server " + s.Address + " [" + string(s.Role) + "]")
	return nil
}

func (s *SebarServer) Stop() error {
	s.Server.Log.Info("Stopping server " + s.Address)
	return nil
}

func (s *SebarServer) AddNode(id string, node *Node) error {
	//--- validate node
	if node.ID == "" {
		return errors.New("Node ID is empty")
	}

	if node.Role == "" {
		return errors.New("Node Role for " + node.ID + " is empty")
	}

	/*
		if node.clientRpc == nil {
			e := node.InitRPC()
			if e != nil {
				return errors.New("Unable to initialize RPC for Node " + node.ID + ": " + e.Error())
			}
		}
	*/

	s.initNodes()
	nodes := s.initNodeType(string(node.Role))
	nodes[id] = node
	s.nodes[string(node.Role)] = nodes
	s.Log.Info("Regitering node " + node.ID + " as [" + string(node.Role) + "] to " + s.Address)
	return nil
}

func (s *SebarServer) initNodes() {
	if s.nodes == nil {
		s.nodes = map[string]map[string]*Node{}
	}
}

func (s *SebarServer) initNodeType(nodeTypeName string) map[string]*Node {
	nodes, nodesExist := s.nodes[nodeTypeName]
	if !nodesExist {
		nodes = map[string]*Node{}
	}
	return nodes
}

func (s *SebarServer) Node(id string) *Node {
	s.initNodes()
	for _, nodes := range s.nodes {
		for _, node := range nodes {
			if node.ID == id {
				return node
			}
		}
	}
	return nil
}

func (s *SebarServer) RemoveNode(id string) {
	if s.nodes == nil {
		return
	}
	//delete(s.nodes, id)
	for nodeTypeName, nodes := range s.nodes {
		for _, node := range nodes {
			if node.ID == id {
				delete(nodes, id)
				s.nodes[nodeTypeName] = nodes
				return
			}
		}
	}
}

/*
func (s *SebarServer) SetURL(rawurl string) *SebarServer {
	u, e := url.Parse(rawurl)
	if e != nil {
		return s
	}
	s.Protocol = u.Scheme
	s.Address = u.Host
	return s
}
*/
