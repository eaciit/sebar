package main

import (
	"fmt"
	"github.com/eaciit/knot"
	//"net/http"
)

type ServerController struct {
}

func (s *ServerController) Status(kr *knot.Request) interface{} {
	ret := "Sebar Server v 0.9 (c) EACIIT \n"
	keys := kr.QueryKeys()
	for _, k := range keys {
		ret += fmt.Sprintf("%s = %v \n", k, kr.Query(k))
	}
	return ret
}

func (s *ServerController) Student(kr *knot.Request) interface{} {
	kr.RouteConfig().OutputType = knot.OutputJson
	return struct {
		ID    string
		Name  string
		Score int
	}{"Student01", "Nama Student", 100}
}

func (s *ServerController) Stop(r *knot.Request) interface{} {
	r.Server().Stop()
	return nil
}
