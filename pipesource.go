package sebar

import (
	"github.com/eaciit/toolkit"
)

type SeekFromEnum int

const (
	SeekFromStart   SeekFromEnum = 0
	SeekFromCurrent SeekFromEnum = 1
)

type IPipeSource interface {
	Data() interface{}
	SetData(interface{}) IPipeSource
	CurrentPos() int
	Len() int
	First() interface{}
	Next() (interface{}, bool)
	Seek(int, SeekFromEnum) interface{}
	IsAsync() bool
	SetAsync(b bool)
	Done() bool
}

type PipeSource struct {
	done       bool
	data       interface{}
	currentPos int
	dataLen    int
	isAsync    bool
}

func (p *PipeSource) IsAsync() bool {
	return p.isAsync
}

func (p *PipeSource) SetAsync(b bool) {
	p.isAsync = b
}

func (p *PipeSource) Data() interface{} {
	return p.data
}

func (p *PipeSource) SetData(data interface{}) IPipeSource {
	if toolkit.IsPointer(data) == false {
		p.data = nil
		return p
	}

	if toolkit.IsSlice(data) == false {
		p.data = nil
		return p
	}

	p.data = data
	p.dataLen = toolkit.SliceLen(p.data)
	p.currentPos = 0
	//toolkit.Println("Data length: ", toolkit.SliceLen(p.data), " Data sample: ", toolkit.SliceSubset(p.data, 0, 20))
	return p
}

func (p *PipeSource) SetDone(d bool) {
	p.done = d
}

func (p *PipeSource) Done() bool {
	return p.done
}

func (p *PipeSource) Len() int {
	return toolkit.SliceLen(p.data)
}

func (p *PipeSource) CurrentPos() int {
	return p.currentPos
}

func (p *PipeSource) First() interface{} {
	p.done = false
	p.currentPos = 0
	if p.Len() == 0 {
		return nil
	}
	return toolkit.SliceItem(p.data, 0)
}

func (p *PipeSource) Next() (interface{}, bool) {
	if p.currentPos >= p.Len()-1 {
		p.done = true
		return nil, false
	}
	p.currentPos++
	return toolkit.SliceItem(p.data, p.currentPos), true
}

func (p *PipeSource) Seek(index int, seekFrom SeekFromEnum) interface{} {
	var newpos int
	if seekFrom == SeekFromStart {
		newpos = index
	} else {
		newpos = p.currentPos + index
	}
	if newpos >= p.Len() {
		return nil
	} else {
		p.currentPos = newpos
		return toolkit.SliceItem(p.data, newpos)
	}
}
