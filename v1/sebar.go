package sebar

import (
	"github.com/eaciit/knot/knot.v1"
)

type NodeRoleEnum string

const (
	WriteMemory  int = 1
	WriteStorage int = 2

	RoleCoordinator        NodeRoleEnum = "COR"
	RoleCoordinatorReplica NodeRoleEnum = "CORR"
	RoleWeb                NodeRoleEnum = "WEB"
	RoleWorker             NodeRoleEnum = "WRK"
	RoleStorage            NodeRoleEnum = "STR"
	RoleStorageReplica                  = "STRR"
)

func init() {
	app := knot.NewApp("sebar-rest")
	app.DefaultOutputType = knot.OutputJson
	knot.RegisterApp(app)
}
