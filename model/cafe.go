package model

import "time"

type Cafe struct {
	OpenMoment    time.Time
	ClosureMoment time.Time
	DescNumber    uint
	CostPerHour   uint
}
