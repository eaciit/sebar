package sebar

import (
	"github.com/eaciit/toolkit"
	"reflect"
	"strings"
	"sync"
	"time"
)

type PipeItemState string

const (
	PipeItemReady   PipeItemState = "Ready"
	PipeItemRunning PipeItemState = "Running"
	PipeItemDone    PipeItemState = "Done"
)

type PipeItem struct {
	sync.Mutex
	attributes         toolkit.M
	nextItem           *PipeItem
	noParralelism      bool
	keyCount           int
	allKeysHasBeenSent bool

	state           PipeItemState
	parallelManager *ParallelManager
	reduceTemp      interface{}
	waiting         bool
	wg              *sync.WaitGroup
}

func (p *PipeItem) initAttributes() {
	if p.attributes == nil {
		p.attributes = toolkit.M{}
	}
}

func (p *PipeItem) AllKeysHasBeenSent() {
	p.allKeysHasBeenSent = true
}

func (p *PipeItem) Set(k string, v interface{}) {
	p.initAttributes()
	p.attributes.Set(k, v)
}

func (p *PipeItem) Get(k string, def interface{}) interface{} {
	p.initAttributes()
	return p.attributes.Get(k, def)
}

func (p *PipeItem) SetError(err string) error {
	return nil
}

func (p *PipeItem) reset() {
	p.keyCount = 0
	p.wg = nil
	p.state = PipeItemReady

	if p.nextItem != nil {
		p.nextItem.reset()
	}
}

func (p *PipeItem) Wait() error {
	//--- one wait should only run once
	if p.waiting {
		return nil
	}

	//--- tell that waiting has been done
	defer func() {
		p.Lock()
		p.waiting = false
		p.Unlock()
	}()

	for {
		if p.allKeysHasBeenSent {
			break
		} else {
			time.Sleep(p.Get("waitduration", 1*time.Second).(time.Duration))
		}
	}

	if p.wg != nil {
		p.wg.Wait()
	}

	p.state = PipeItemDone
	return nil
}

func (p *PipeItem) verbose(txt string) {
	//parm := p.parm()
	toolkit.Printf("[%d] %s \n", p.keyCount, txt)
}

func (p *PipeItem) send(k interface{}) {
	if p.wg == nil {
		p.wg = new(sync.WaitGroup)
	}

	parm := p.parm()

	p.keyCount++
	edata := toolkit.M{}
	edata.Set("data", k)
	edata.Set("dataindex", p.keyCount)

	if parm.Get("verbose", false).(bool) {
		p.verbose(toolkit.Sprintf("p.send: %s", toolkit.JsonString(edata)))
	}

	p.wg.Add(1)
	p.execute(edata)
}

func (p *PipeItem) parm() toolkit.M {
	parm := p.Get("parm", nil)
	if parm == nil {
		return toolkit.M{}
	} else {
		return parm.(toolkit.M)
	}
}

/*
func wgDone(wg *sync.WaitGroup) {
	if wg != nil {
		//toolkit.Println("Done 1 elem of WaitGroup")
		wg.Done()
	}
}
*/

func (p *PipeItem) execute(executeData toolkit.M) {
	defer p.wg.Done()
	if p.state != PipeItemRunning {
		p.state = PipeItemRunning
	}
	p.sendToNext(executeData)
	return
}

func (p *PipeItem) sendToNext(executeData toolkit.M) {
	if p.nextItem == nil {
		return
	} else {
		p.nextItem.Set("parm", p.parm())
		if executeData == nil {
			executeData = toolkit.M{}
		}
		in := executeData.Get("data", nil)
		p.nextItem.send(in)
	}
}

func (p *PipeItem) _Run(dataRun toolkit.M) error {
	op := strings.ToLower(p.Get("op", "").(string))
	parm := p.Get("parm", toolkit.M{}).(toolkit.M)
	verbose := parm.Get("verbose", false).(bool)
	pIn := p.Get("in", nil)
	var wg *sync.WaitGroup
	wg = parm.Get("wg", wg).(*sync.WaitGroup)

	if op == "" {
		//p.Set("error", "OP is mandatory")
		//wgDone(wg)
		return p.SetError("OP is mandatory")
	}

	if op == "parallel" {
		if p.nextItem == nil {
			return p.SetError("NextItem is nil. Parallel should be following with another PipeItem")
		} else {
			/*
				p.nextItem.Set("parm", p.Get("parm", nil))
				p.nextItem.Set("in", p.Get("in", nil))
				return p.nextItem.Run()
			*/
			if p.parallelManager == nil {
				p.parallelManager, _ = NewParallelManager(p.Get("parallel", 1).(int), p.nextItem)
				p.parallelManager.parm = parm
				p.parallelManager.Wait()
			}

			p.parallelManager.SendKey(pIn)
			p.parallelManager.allKeysHasBeenSent = p.allKeysHasBeenSent
			return nil
		}
	}

	//fn := p.Get("fn_"+op, nil)
	fn := p.Get("fn", nil)
	if fn == nil {
		//wgDone(wg)
		return p.SetError(toolkit.Sprintf("Function %s is not available", op))
	}

	vfn := reflect.Indirect(reflect.ValueOf(fn))
	if vfn.Kind() != reflect.Func {
		return p.SetError(toolkit.Sprintf("Function %s is not a function", op))
	}

	var ins []reflect.Value
	var outs []reflect.Value

	if !toolkit.IsSlice(pIn) {
		ins = append(ins, reflect.ValueOf(pIn))
	} else {
		pLen := toolkit.SliceLen(pIn)
		for pIndex := 0; pIndex < pLen; pIndex++ {
			ins = append(ins, reflect.ValueOf(toolkit.SliceItem(pIn, pIndex)))
		}
	}

	//toolkit.Println(toolkit.JsonString(ins))
	tfn := vfn.Type()
	lenIn := tfn.NumIn()
	if len(ins) < lenIn {
		for i := len(ins); i < lenIn; i++ {
			fnin := reflect.New(tfn.In(i)).Elem()
			ins = append(ins, fnin)
		}
	}

	if op == "mapreduce" {
		if len(ins) > 0 {
			if p.reduceTemp == nil {
				p.reduceTemp = ins[len(ins)-1].Interface()
			} else {
				ins[len(ins)-1] = reflect.ValueOf(p.reduceTemp)
			}
			//toolkit.Println("mapreduce set reduceTemp", p.reduceTemp)
		}
	}

	/*
		if verbose {
			toolkit.Printf("Data %d Pipe %d %s: %s",
				p.Get("parm", nil).(toolkit.M).Get("dataindex", 0).(int),
				p.Get("index", 0).(int), op,
				toolkit.JsonString(pIn))
		}
	*/

	outs = vfn.Call(ins)

	var iouts []interface{}
	for _, out := range outs {
		iouts = append(iouts, out.Interface())
	}

	if verbose {
		toolkit.Printf("Data %d Pipe %d %s: %s => %s\n",
			p.Get("parm", nil).(toolkit.M).Get("dataindex", 0).(int),
			p.Get("index", 0).(int), op,
			toolkit.JsonString(pIn),
			toolkit.JsonString(iouts))
	}

	if op == "where" && iouts[0] == false {
		//wgDone(wg)
		return nil
	} else if op == "where" && iouts[0] == true {
		iouts = []interface{}{}
		for _, in := range ins {
			iouts = append(iouts, in.Interface())
		}
	}

	if op == "mapreduce" && len(iouts) > 0 {
		p.reduceTemp = iouts[0]
	}

	//p.Set("output", outs)
	if p.nextItem != nil {
		p.nextItem.allKeysHasBeenSent = p.allKeysHasBeenSent
		p.nextItem.Set("parm", parm)
		p.nextItem.Set("in", iouts)
		return p.nextItem._Run(nil)
	} else {
		if wg != nil {
			wg.Done()
		}
		p.Set("output", iouts)
	}

	return nil
}
