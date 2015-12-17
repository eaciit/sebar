package sebar

type Node struct {
	SebarServer

	nodeRole      NodeRoleEnum
	masterAddress string
	masterSecret  string
}

func (n *Node) SetMaster(url, secret string) {
	n.masterAddress = url
	n.masterSecret = secret
}
