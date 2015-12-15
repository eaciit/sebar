package test

import (
	"fmt"
	"github.com/eaciit/toolkit"
	//"strings"
	"testing"
)

const (
	restPrefix = "http://localhost:12000/r/"
)

var token string
var e error

func getToken() error {
	surl := restPrefix + "auth/init"
	b, e := toolkit.EncodeByte(struct {
		Key    string
		Secret string
	}{"ariefdarmawan", "Password.1"})
	if e != nil {
		return fmt.Errorf("Unable to connect: %s", e.Error())
	}
	r, e := toolkit.HttpCall(surl, "get", b, nil)

	if e != nil {
		return fmt.Errorf("Unable to connect: %s", e.Error())
	}

	token := toolkit.HttpContentM(r).Get("Token", "").(string)
	if token == "" {
		return fmt.Errorf("Wrong token returned")
	}

	return nil
}

func he(t *testing.T, err error) {
	t.Errorf(err.Error())
}

func createRandomString(randomLength int) string {
	randomTxt := ""
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890!@#$%_-"
	for x := 0; x < randomLength; x++ {
		ic := toolkit.RandInt(len(chars) - 1)
		c := chars[ic]
		randomTxt += string(c)
	}
	return randomTxt
}

func TestPut(t *testing.T) {
	if e := getToken(); e != nil {
		//he(t, e)
		//return
	}

	surl := restPrefix + "/put"
	for i := 1; i <= 1000; i++ {
		randomTxt := createRandomString(toolkit.RandInt(22) + 10)
		data := struct {
			Token string
			Data  string
		}{
			token,
			randomTxt,
		}

		fmt.Printf("Saving %d value %s", i, data.Data)

		r, e := toolkit.HttpCall(surl, "XPUT", toolkit.GetEncodeByte(data), nil)
		//e = nil
		if e != nil {
			fmt.Printf("... Fail: %s \n", e.Error())
		} else {
			if r.StatusCode != 200 {
				fmt.Printf("... Fail: %d %s \n", r.StatusCode, r.Status)
			} else {
				fmt.Println("...Done")
			}
			fmt.Println("...Done")
		}
	}
}
