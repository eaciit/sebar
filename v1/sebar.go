package sebar

import (
	"github.com/eaciit/knot/knot.v1"
)

const (
	RoleMaster        int = 1
	RoleNode          int = 5
	RoleReplicaMaster int = 101
	RoleReplicaNode   int = 105
)

func init() {
	app := knot.NewApp("sebar-rest")
	app.DefaultOutputType = knot.OutputJson
	knot.RegisterApp(app)
}
