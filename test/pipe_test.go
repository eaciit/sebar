package v1dev_test

import (
	"github.com/eaciit/sebar"
	"github.com/eaciit/toolkit"
	"testing"
	"time"
)

var dataCount int = 200
var pipe *sebar.Pipe
var dataPipe []int
var outs []int

type DataOut struct {
	Group int
	X     int
}

func TestPrepareData(t *testing.T) {
	for i := 0; i < dataCount; i++ {
		dataPipe = append(dataPipe, toolkit.RandInt(600)+1)
	}
	if len(dataPipe) != dataCount {
		t.Fatalf("Error: want %d data got %d", dataCount, len(dataPipe))
	}
	toolkit.Println("Data (20 samples): ", toolkit.JsonString(dataPipe[:20]))
}

func TestLoad(t *testing.T) {
	ds := new(sebar.PipeSource).SetData(&dataPipe)
	pipe1 := new(sebar.Pipe).From(ds).SetOutput(&outs)
	e := pipe1.Exec(nil)
	if e != nil {
		t.Fatalf("Error load: " + e.Error())
	}
	if len(outs) != len(dataPipe) {
		t.Fatalf("Error: want %d data got %d", len(dataPipe), len(outs))
	}
	for idx, val := range dataPipe {
		if val != outs[idx] {
			t.Fatalf("Data %d is not same. Expect %d got %d",
				idx, val, outs[idx])
		}
	}
	toolkit.Printfn("Data: " + toolkit.JsonString(outs[0:20]))
}

func TestWhereMap(t *testing.T) {
	var outsmap []struct {
		X int
		Y int
	}
	pipe1 := new(sebar.Pipe).From(new(sebar.PipeSource).SetData(&dataPipe))
	pipe1.Where(func(x int) bool {
		return x <= 300
	})
	pipe1.Map(func(x int) struct {
		X int
		Y int
	} {
		return struct {
			X int
			Y int
		}{x, x * 2}
	})

	pipe1.SetOutput(&outsmap)
	if pipe1.ErrorTxt() != "" {
		t.Fatalf("Error: %s", pipe1.ErrorTxt())
	}

	eExec := pipe1.Exec(toolkit.M{}.Set("verbose", false))
	if eExec != nil {
		t.Fatalf("Exec: %s", eExec.Error())
	}
	if len(outsmap) == 0 {
		t.Fatalf("No record returned")
	}

	for idx, v := range outsmap {
		if v.X > 300 {
			t.Fatalf("Data index %d, %d > 100", idx, v.X)
		}
	}
	toolkit.Printfn("Data: " + toolkit.JsonString(outsmap))
}

func TestWhereMapReduce(t *testing.T) {
	type xy struct {
		X int
		Y int
	}

	var total1, total2 int
	var ints []int

	pipe1 := new(sebar.Pipe).From(new(sebar.PipeSource).SetData(&dataPipe))
	pipe1.Where(func(x int) bool {
		return x <= 300
	}).Map(func(x int) xy {
		time.Sleep(10 * time.Millisecond)
		return xy{x, x * 2}
	}).Reduce(func(m xy, b int) int {
		ints = append(ints, m.Y)
		return b + m.Y
	}).SetOutput(&total1)
	if pipe1.ErrorTxt() != "" {
		t.Fatalf("Error: %s", pipe1.ErrorTxt())
	}
	eExec := pipe1.Exec(toolkit.M{}.Set("verbose", false))
	if eExec != nil {
		t.Fatalf("Exec: %s", eExec.Error())
	}
	if len(ints) == 0 {
		t.Fatalf("No record returned")
	}

	for _, v := range ints {
		total2 += v
	}
	if total1 != total2 {
		t.Fatalf("Summation error. Expect %d got %d", total2, total1)
	}
	toolkit.Printfn("Total: %d Data: %s", total1, toolkit.JsonString(ints))
}

func TestWhereMapReducePartition(t *testing.T) {
	type xy struct {
		X int
		Y int
	}

	var total1, total2 int
	var ints []int

	pipe1 := new(sebar.Pipe).From(new(sebar.PipeSource).SetData(&dataPipe))
	pipe1.Parallel(5).Where(func(x int) bool {
		return x <= 300
	}).Map(func(x int) xy {
		time.Sleep(10 * time.Millisecond)
		return xy{x, x * 2}
	}).Reduce(func(m xy, b int) int {
		ints = append(ints, m.Y)
		return b + m.Y
	}).SetOutput(&total1)
	if pipe1.ErrorTxt() != "" {
		t.Fatalf("Error: %s", pipe1.ErrorTxt())
	}
	eExec := pipe1.Exec(toolkit.M{}.Set("verbose", true))
	if eExec != nil {
		t.Fatalf("Exec: %s", eExec.Error())
	}

	eWait := pipe1.Wait()
	if eWait != nil {
		t.Fatalf("Wait: %s:", eWait.Error())
	}

	if len(ints) == 0 {
		t.Fatalf("No record returned")
	}

	for _, v := range ints {
		total2 += v
	}
	if total1 != total2 {
		t.Fatalf("Summation error. Expect %d got %d", total2, total1)
	}
	toolkit.Printfn("Total: %d Data: %s", total1, toolkit.JsonString(ints))
}

/*
func TestPipe(t *testing.T) {
	t.Skip()
	pipe1 := new(sebar.Pipe).From(nil).Map(func(x int) DataOut {
		return DataOut{x / 100, x}
	}).Sort(func(x DataOut) int {
		return x.Group
	})

	pipe2 := new(sebar.Pipe).From(nil)

	pipe3 := new(sebar.Pipe).Join(pipe1, pipe3, func(x DataOut, y int) bool {
		return x.Group == y
	}, func(x DataOut, y int) DataOut {
		return x.Group
	}).Reduce(func(x DataOut, prev int) (int, int) {
		return x.Group, prev + int
	})

	pipe3.ParseAndExec(nil)
	if pipe3.Error != nil {
		t.Fatalf("Error: %s", pipe3.Error.Error())
	}
	t.Logf("P1:\n%s\n"+
		"P2:\n%s\n"+
		"P3:\n%s\n",
		toolkit.JsonString(pipe1.Data),
		toolkit.JsonString(pipe2.Data),
		toolkit.JsonString(pipe3.Data))
i}
*/
