package sebar

type Node struct {
	SebarServer

	Role          NodeRoleEnum
	masterAddress string
	masterSecret  string
}

func (n *Node) SetMaster(url, secret string) {
	n.masterAddress = url
	n.masterSecret = secret
}

func (n *Node) Key() string {
	return string(n.Role) + "_" + n.url
}
