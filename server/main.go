package main

import (
	"fmt"
	"github.com/eaciit/kingpin"
	"github.com/eaciit/knot"
	//"net/http"
)

var (
	port = kingpin.Flag("port", "Port to be used to listen for request").Default("13000").String()

	e error
)

func main() {
	kingpin.Parse()

	ks := new(knot.Server)
	ks.Address = ":" + *port
	ks.Route("/", Index, nil)
	/*
		ks.Route("/s/status", Status, nil)
		ks.Route("/s/stop", Stop, nil)
	*/
	e = ks.Register(new(ServerController), "", nil)
	if e != nil {
		fmt.Println("Error: " + e.Error())
		return
	}

	ks.Listen()
}

func Index(r *knot.Request) interface{} {
	return "Welcome to Sebar Server"
}
