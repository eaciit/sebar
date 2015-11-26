package main

import (
	"github.com/eaciit/knot"
	//"net/http"
)

type ServerController struct {
}

func (s *ServerController) Status(r *knot.Request) interface{} {
	return "Sebar Server v 0.9 (c) EACIIT"
}

func (s *ServerController) Stop(r *knot.Request) interface{} {
	r.Server().Stop()
	return nil
}
