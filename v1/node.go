package sebar

type Node struct {
	SebarServer

	masterAddress string
	masterSecret  string
}

func (n *Node) SetMaster(url, secret string) {
	n.masterAddress = url
	n.masterSecret = secret
}
