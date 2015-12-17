package sebar

func NewServer(role NodeRoleEnum, url string) IServer {
	if role == RoleMaster {
		m := new(Master)
		m.Address = url
		return m
	} else if role == RoleNode {
		n := new(Node)
		n.Address = url
		return n
	}
	return nil
}

func MakeUrl(serverUrl, role, method string) string {
	return serverUrl + "/" + role + "/" + method
}
