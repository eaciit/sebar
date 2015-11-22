package sebar

import (
	"fmt"
	"github.com/eaciit/errorlib"
	"github.com/eaciit/toolkit"
	"net/http"
	"time"
)

// ServerType define type of a Sebar server
type ServerType int

const (
	// ServerTypeMaster responsible as main proxy of Sebar request
	ServerTypeMaster ServerType = 1 // ServerTypeNode to write data. Master node is also at same time a node
	ServerTypeNode   ServerType = 100
	// ServerTypeReplicaMaster handle data replica request to each of its replica node
	ServerTypeReplicaMaster ServerType = 1001
	// ServerTypeReplicaNode replica, handle data replica
	ServerTypeReplicaNode ServerType = 1100
)

// String returns string value of ServerType
func (s ServerType) String() string {
	if s == ServerTypeMaster {
		return "Master"
	} else if s == ServerTypeNode {
		return "Node"
	} else if s == ServerTypeReplicaMaster {
		return "Replica Master"
	} else if s == ServerTypeReplicaNode {
		return "Replica Node"
	}
	return ""
}

/*
Server used to handle main operation on Sebar Server
*/
type Server struct {
	ID         string
	ServerType ServerType
	Map        map[string]*DataMap
	Items      map[string]*Data

	MemorySize int

	Nodes         []*Server
	ReplicaMaster []*Server
}

type DataKey struct {
	//-- auth
	Token  string
	UserID string

	//-- key builder
	Owner   string
	Cluster string
	Key     string

	//-- data related
	DataType string
	Expiry   int
}

func (d *DataKey) BuildKey() string {
	if d.Owner == "" {
		d.Owner = "public"
	}
	if d.Cluster == "" {
		d.Cluster = "general"
	}
	return fmt.Sprintf("%s.%s.%s", d.Owner, d.Cluster, d.Key)
}

func (s *Server) Write(datakey DataKey, data []byte) error {
	handleError := func(e error) error {
		return errorlib.Error(packageName, modServer, "Write",
			fmt.Sprintf("[%s.%s.%s] ", datakey.Owner, datakey.Cluster, datakey.Key)+e.Error())
	}

	if s.ServerType == ServerTypeMaster || s.ServerType == ServerTypeReplicaMaster {
		dataSize := len(data)
		snode, e := s.getNode(ServerTypeNode, dataSize)
		if e != nil {
			return handleError(e)
		}

		snodes := []*Server{snode}
		srepma, e := s.getNode(ServerTypeReplicaMaster, dataSize)
		if srepma != nil {
			snodes = append(snodes, srepma)
		}

		for _, sn := range snodes {
			mDataKey, e := toolkit.ToM(datakey)
			if e != nil {
				return handleError(e)
			}
			_, e = sn.call(opWrite, "GET", mDataKey, data)
			if e != nil {
				return handleError(e)
			}

			e = s.writeDataMap(datakey, dataSize, sn)
			if e != nil {
				return handleError(e)
			}
		}
	} else {
		key := datakey.BuildKey()
		d := new(Data)
		d.Key = key
		d.Type = datakey.DataType
		if datakey.Expiry == 0 {
			d.Expiry = time.Now().Add(100 * 365 * 24 * 60 * time.Minute)
		} else {
			d.Expiry = time.Now().Add(time.Duration(datakey.Expiry))
		}
		d.Value = data
		s.Items[key] = d
	}
	return nil
}

func (s *Server) call(op string, calltype string, qs toolkit.M, payload []byte) (*http.Response, error) {
	url := s.Id
	url += op
	q := ""
	for k, v := range qs {
		if q == "" {
			q = "?" + k + "=" + v.(string)
		} else {
			q += "&" + k + "=" + v.(string)
		}
	}
	url += q

	//rdr := bytes.NewReader(payload)
	r, e := toolkit.HttpCall(url, calltype, payload, nil)
	return r, e
}

func (s *Server) writeDataMap(datakey DataKey, dataSize int, snode *Server) error {
	return nil
}

func (s *Server) AvailableMemory() int {
	return 0
}

func (s *Server) MemoryUsage() int {
	return 0
}

func (s *Server) getNode(serverType ServerType, dataSize int) (*Server, error) {
	var ss []*Server
	if serverType == ServerTypeNode {
		ss = s.Nodes
	} else if serverType == ServerTypeReplicaMaster || serverType == ServerTypeReplicaNode {
		ss = s.ReplicaMaster
	}

	var r *Server
	for _, s := range ss {
		if r == nil {
			r = s
		} else {
			if (s.AvailableMemory() >= dataSize || s.MemorySize == 0) &&
				(s.MemoryUsage() < r.MemoryUsage()) {
				r = s
			}
		}
	}

	if r == nil {
		return nil,
			fmt.Errorf("No node %s available to write %d bytes of data",
				string(serverType), dataSize)
	}

	return s, nil
}
