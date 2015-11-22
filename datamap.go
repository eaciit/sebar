package sebar

type DataMap struct {
	//--- Key
	Key string

	//--- Chain
	IsParentCluster bool
	Items           map[string]DataMap

	//--- Position
	Node    *Server
	Replica []*Server
}
