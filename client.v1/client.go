package sebarclient

import (
	//"github.com/eaciit/toolkit"
	"github.com/eaciit/appserver/v1"
	"github.com/eaciit/errorlib"
)

const (
	objectName = "sebarclient"
	modClient  = "Client"
)

var e error

type Client struct {
	Host, UserID, Secret string

	clientRpc *appserver.Client
}

func (c *Client) Connect() error {
	c.clientRpc = new(appserver.Client)
	e = c.clientRpc.Connect(c.Host, c.Secret, c.UserID)
	if e != nil {
		return errorlib.Error(objectName, modClient, "Connect", e.Error())
	}
	return nil
}

func (c *Client) Close() {
	c.clientRpc.Close()
}
