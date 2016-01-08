package sebar

import (
	"github.com/eaciit/toolkit"
)

func NewServer(role NodeRoleEnum, url string) IServer {
	if role == RoleCoordinator {
		m := new(Coordinator)
		m.Role = role
		m.Address = url
		return m
	} else if role == RoleStorage {
		s := new(Storage)
		s.Role = role
		s.Address = url
		return s
	} else if role == RoleStorageReplica {
		//s := new(StorageReplica)
		//return s
	}
	return nil
}

func MakeUrl(serverUrl, role, method string) string {
	return serverUrl + "/" + role + "/" + method

}

func ParseSize(size float64) string {
	return doParseSize(size, "")
}

func doParseSize(size float64, unit string) string {
	if unit == "" {
		unit = "B"
	}
	ret := ""
	if size > 1024 {
		size = size / 1024
		if unit == "B" {
			unit = "K"
		} else if unit == "K" {
			unit = "M"
		} else if unit == "M" {
			unit = "G"
		} else if unit == "G" {
			unit = "T"
		} else {
			unit = "P"
		}

		if unit != "P" {
			ret = doParseSize(size, unit)
		} else {
			ret = toolkit.Sprintf("%2.2f%s", size, unit)
		}
	} else {

		ret = toolkit.Sprintf("%2.2f%s", size, unit)
	}
	return ret
}
