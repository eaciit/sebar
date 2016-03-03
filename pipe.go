package sebar

import (
	"errors"
	"github.com/eaciit/crowd"
	"github.com/eaciit/toolkit"
	"reflect"
	"sync"
	"time"
)

type ApplyScope string

const (
	ScopeLocal   ApplyScope = "local"
	ScopeGlobal  ApplyScope = "global"
	ScopeCluster ApplyScope = "cluster"
)

type Pipe struct {
	c *crowd.Crowd

	Items []*PipeItem

	source IPipeSource
	parsed bool
	err    error
	output interface{}

	waitGroup          *sync.WaitGroup
	_waitingPeriod     time.Duration
	allKeysHasBeenSent bool
}

func (p *Pipe) SetWaitingPeriod(t time.Duration) {
	p._waitingPeriod = t
}

func (p *Pipe) AllKeysHasBeenSent() {
	p.allKeysHasBeenSent = true
	if len(p.Items) > 0 {
		p.Items[0].AllKeysHasBeenSent()
	}
}

func (p *Pipe) WaitingPeriod() time.Duration {
	if p._waitingPeriod == 0 {
		p._waitingPeriod = time.Millisecond * 1
	}
	return p._waitingPeriod
}

func (p *Pipe) SetError(s string) {
	p.err = errors.New(s)
}

func (p *Pipe) Error() error {
	return p.err
}

func (p *Pipe) ErrorTxt() string {
	if p.err == nil {
		return toolkit.Sprintf("")
	}
	return p.err.Error()
}

func (p *Pipe) Parsed() bool {
	return p.parsed
}

func (p *Pipe) Parse() error {
	p.err = nil
	p.parsed = true
	return p.err
}

func (p *Pipe) Exec(parm toolkit.M) error {
	if p.source == nil {
		return errors.New("Pipe.Exec: Source is invalid")
	}

	if len(p.Items) == 0 {
		if p.output != nil {
			e := toolkit.Serde(p.source.Data(), p.output, "json")
			if e != nil {
				return errors.New("Pipe.Exec: unable to serde the result " + e.Error())
			}
		}
		return nil
	}

	if parm == nil {
		parm = toolkit.M{}
	}

	parm.Set("verbose", true)
	p.Items[0].Set("parm", parm)
	running := true
	dataIndex := -1
	p.source.First()
	p.Items[0].reset()
	p.allKeysHasBeenSent = false
	go func() {
		for running {
			dataItem, hasData := p.source.Next()
			if hasData {
				dataIndex++
				p.Items[0].send(dataItem)
			} else {
				running = false
			}
		}
	}()
	return nil
}

func (p *Pipe) Wait() error {
	if len(p.Items) == 0 {
		return nil
	}

	ewait := p.Items[0].Wait()
	if ewait != nil {
		p.SetError("pipe.Wait" + ewait.Error())
	}

	return nil
}

/*
func (p *Pipe) ParseAndExec(inputs interface{}, reparse bool) {
	if reparse || p.parsed == false {
		p.Parse()
	}
	if p.Error() != nil {
		return
	}
	p.Exec(inputs)
}
*/

func (p *Pipe) Parallel(i int) *Pipe {
	//p.Set("partition", i)
	pi := new(PipeItem)
	pi.Set("op", "parallel")
	pi.Set("parallel", i)
	eadd := p.addItem(pi)
	if eadd != nil {
		p.SetError(eadd.Error())
		return p
	}
	return p
}

func (p *Pipe) SetOutput(o interface{}) *Pipe {
	pi := new(PipeItem)
	pi.noParralelism = true
	pi.Set("op", "setoutput")
	pi.Set("fn", func(x interface{}) {
		if toolkit.IsSlice(o) {
			toolkit.AppendSlice(o, x)
		} else {
			reflect.ValueOf(o).Elem().Set(reflect.ValueOf(x))
		}
	})
	eadd := p.addItem(pi)
	if eadd != nil {
		p.SetError(eadd.Error())
		return p
	}
	p.output = o
	return p
}

func (p *Pipe) Join(p1 *Pipe, p2 *Pipe, fnJoin interface{}) *Pipe {
	return p
}

func (p *Pipe) From(s IPipeSource) *Pipe {
	p.source = s
	return p
}

func (p *Pipe) Where(fn interface{}) *Pipe {
	pi := new(PipeItem)
	pi.Set("op", "where")
	pi.Set("fn", fn)
	p.addItem(pi)
	return p
}

func (p *Pipe) Map(fn interface{}) *Pipe {
	pi := new(PipeItem)
	pi.Set("op", "map")
	pi.Set("fn", fn)
	p.addItem(pi)
	return p
}

func (p *Pipe) Sort(fn interface{}) *Pipe {
	return p
}

func (p *Pipe) Reduce(fn interface{}) *Pipe {
	pi := new(PipeItem)
	pi.noParralelism = true
	pi.Set("op", "mapreduce")
	pi.Set("fn", fn)
	_ = p.addItem(pi)
	/*
		if eadd != nil {
			toolkit.Println("AddReduce:", eadd.Error())
		}
	*/
	return p
}

func (p *Pipe) addItem(pi *PipeItem) error {
	if p.ErrorTxt() != "" {
		return errors.New("Pipe.addPipeItem: " + p.ErrorTxt())
	}

	if pi == nil {
		return errors.New("Pipe.addPipeItem: PipeItem is nil")
	}

	if len(p.Items) > 0 {
		lastpi := p.Items[len(p.Items)-1]
		if lastpi.Get("op", "") == "setoutput" {
			return errors.New("Pipe.addPipeItem: Last PipeItem is SetOutput. No more PipeItem can't be inserted after SetOutput")
		}
		lastpi.nextItem = pi
	}

	pi.Set("index", len(p.Items))
	p.Items = append(p.Items, pi)

	return nil
}
