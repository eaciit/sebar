package sebar

import (
	"time"
)

type DataPoint struct {
	ID          string
	Value       []byte
	Created     time.Time
	ExpiryAfter time.Duration
	ExpiredOn   time.Time
}

func (d *DataPoint) setExpiry(ed time.Duration) {
	if ed != 0 {
		d.ExpiryAfter = ed
	}
	d.ExpiredOn = time.Now().Add(d.ExpiryAfter)
}
