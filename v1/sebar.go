package sebar

import (
	"github.com/eaciit/knot/knot.v1"
)

type NodeRoleEnum string

const (
	RoleMaster        NodeRoleEnum = "master"
	RoleNode          NodeRoleEnum = "node"
	RoleReplicaMaster NodeRoleEnum = "replicamaster"
	RoleReplicaNode   NodeRoleEnum = "replicanode"
)

func init() {
	app := knot.NewApp("sebar-rest")
	app.DefaultOutputType = knot.OutputJson
	knot.RegisterApp(app)
}
