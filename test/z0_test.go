package sebar_test

import (
	//"github.com/eaciit/sebar/client.v1"
	"github.com/eaciit/appserver/v1"
	"github.com/eaciit/crowd"
	"github.com/eaciit/sebar/v1"
	"github.com/eaciit/toolkit"
	//"strconv"
	"fmt"
	"strings"
	"testing"
	//"time"
)

var (
	secretMaster = "master"
	userID       = "admin"
	password     = "admin"
	userSecret   = ""
	masterUrl    = "localhost:6789"

	coordinator *sebar.Coordinator
	nodes       []*sebar.Node
	session     *sebar.Session

	client *appserver.Client
)

func TestMaster(t *testing.T) {
	coordinator = sebar.NewServer(sebar.RoleCoordinator, masterUrl).(*sebar.Coordinator)
	coordinator.AddUser(userID, password)
	coordinator.AllowMultiLogin = true
	e := coordinator.Start()
	if e != nil {
		t.Error(e)
	}
}

func TestStorageNode(t *testing.T) {
	for i := 0; i < 5; i++ {
		var storage *sebar.Storage
		storage = sebar.NewServer(sebar.RoleStorage, fmt.Sprintf("localhost:%d", 9601+i)).(*sebar.Storage)
		storage.Coordinator = coordinator.Address
		storage.CoordinatorUserID = userID
		storage.CoordinatorSecret = password
		e := storage.Start()
		if e != nil {
			t.Error(e)
		}
	}
}

func skipIfClientNil(t *testing.T) {
	if client == nil {
		t.Skip()
	}
}

func TestClient(t *testing.T) {
	client = new(appserver.Client)
	e := client.Connect(coordinator.Address, password, userID)
	if e != nil {
		t.Error(e.Error())
	}
}

var totalInt int

func TestStorageWrite(t *testing.T) {
	skipIfClientNil(t)
	es := []string{}

	toolkit.Printf("Writing Data:\n")
	for i := 0; i < 200; i++ {
		dataku := toolkit.RandInt(1000)
		totalInt += dataku
		//toolkit.Printf("%d ", dataku)

		in := toolkit.M{}.Set("key", fmt.Sprintf("public.dataku.%d", i)).Set("data", dataku)
		//.Set("encode", false).Set("encoderid", "")
		writeResult := client.Call("set", in)
		if writeResult.Status != toolkit.Status_OK {
			es = append(es, toolkit.Sprintf("Fail to write data %d : %d => %s", i, dataku, writeResult.Message))
		}
	}

	if len(es) > 0 {
		errorTxt := ""
		if len(es) <= 10 {
			errorTxt = strings.Join(es, "\n")
		} else {
			errorTxt = strings.Join(es[:10], "\n") + "\n... And others ..."
		}
		t.Errorf("Write data fail.\n%s", errorTxt)
	}
}

func TestStorageGet(t *testing.T) {
	skipIfClientNil(t)

	in := toolkit.M{}
	in.Set("key", "dataku")
	getResult := client.Call("get", in)
	if getResult.Status != toolkit.Status_OK {
		t.Errorf(getResult.Message)
		return
	}

	data := []int{}
	e := toolkit.FromBytes(getResult.Data.([]byte), "", &data)
	if e != nil {
		t.Errorf("Unable to decode: " + e.Error())
	}
	fmt.Println("Result received: %s \n", toolkit.JsonString(data))

	total := int(crowd.From(data).Sum(nil))

	if total != totalInt {
		t.Errorf("Wrong summation expecting %d got %d", totalInt, total)
	}
}

func TestClose(t *testing.T) {
	var e error

	if client != nil {
		client.Close()
	}

	if coordinator != nil {
		e = coordinator.Stop()
		if e != nil {
			t.Error(e)
			return
		}
	}
}
