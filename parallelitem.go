package sebar

import (
	"errors"
	"fmt"
	"github.com/eaciit/toolkit"
	"sync"
	"time"
)

type ParallelManager struct {
	sync.Mutex
	items []*PipeItem
	keys  [][]interface{}

	parm               interface{}
	in                 interface{}
	allKeysHasBeenSent bool
	waiting            bool
	_waitingPeriod     time.Duration
}

func (pm *ParallelManager) SetWaitingPeriod(t time.Duration) {
	pm._waitingPeriod = t
}

func (pm *ParallelManager) WaitingPeriod() time.Duration {
	if pm._waitingPeriod == 0 {
		pm._waitingPeriod = time.Millisecond * 1
	}
	return pm._waitingPeriod
}

func NewParallelManager(partitionNo int, modelItem *PipeItem) (pm *ParallelManager, e error) {
	if modelItem.noParralelism {
		e = errors.New(fmt.Sprintf("crows.NewParallelManager: modelItem to be parallelized does not support for Parallel operation"))
		return
	}

	pm = new(ParallelManager)
	for i := 0; i < partitionNo; i++ {
		pi := new(PipeItem)
		*pi = *modelItem
		pm.items = append(pm.items, pi)
		pm.keys = append(pm.keys, []interface{}{})
	}
	return
}

func copyItem(pi *PipeItem) {
	if pi.nextItem == nil {
		return
	}

	if pi.nextItem.noParralelism {
		return
	}

	nextItem := new(PipeItem)
	*nextItem = *(pi.nextItem)
	copyItem(nextItem)
	pi.nextItem = nextItem
}

func (pm *ParallelManager) SendKey(k interface{}) {
	pm.Lock()
	var least, leastIndex int
	for i := 0; i < len(pm.keys); i++ {
		if i == 0 {
			least = len(pm.keys[i])
			leastIndex = 0
		} else {
			currentLen := len(pm.keys[i])
			if currentLen < least {
				least = currentLen
				leastIndex = i
			}
		}
		if least == 0 {
			break
		}
	}

	pm.keys[leastIndex] = append(pm.keys[leastIndex], k)
	pm.Unlock()

	if !pm.waiting {
		go func() {
			pm.Wait()
		}()
	}

	//fmt.Println("Sending Key", k, pm.keys)
}

func (pm *ParallelManager) Wait() (e error) {
	if pm.waiting {
		return
	}
	defer func() {
		pm.Lock()
		pm.waiting = false
		pm.Unlock()
	}()

	pm.Lock()
	pm.waiting = true
	pm.Unlock()

	fmt.Println("Wait Start")
	for {
		var max, maxIndex int
		for i := 0; i < len(pm.keys); i++ {
			if i == 0 {
				max = len(pm.keys[i])
				maxIndex = 0
			} else {
				currentLen := len(pm.keys[i])
				if currentLen > max {
					max = currentLen
					maxIndex = i
				}
			}
		}

		if max == 0 {
			return
		}

		keys := pm.keys[maxIndex]
		k := keys[len(keys)-1]
		if len(keys) > 1 {
			//fmt.Println("Processing key ", k)
			pm.Lock()
			pm.keys[maxIndex] = pm.keys[maxIndex][:len(keys)-1]
			var wg *sync.WaitGroup
			wg = pm.parm.(toolkit.M).Get("wg", wg).(*sync.WaitGroup)
			go func() {
				pm.items[maxIndex].Set("parm", pm.parm)
				pm.items[maxIndex].Set("in", k)
				erun := pm.items[maxIndex]._Run(nil)
				if erun != nil {
					wgDone(wg)
					fmt.Println("Error", erun.Error())
				} else {
					//fmt.Print("Run")
				}
			}()
			pm.Unlock()
		} else {
			time.Sleep(pm.WaitingPeriod())
		}
	}

	return
}
