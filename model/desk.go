package model

import (
	"fmt"
	"time"
)

type Desk struct {
	Id                   uint
	Occupied             bool
	LastUsageStartMoment time.Time
	TotalUsedTime        time.Time
	TotalRevenue         uint
}

func (desk *Desk) HandleClientLeaveAction(time time.Time, costPerHour uint) {
	desk.Occupied = false
	timeWasted := time.Sub(desk.LastUsageStartMoment)
	desk.TotalUsedTime = desk.TotalUsedTime.Add(timeWasted)
	desk.TotalRevenue += uint((timeWasted.Minutes()+59)/60) * costPerHour
}

func (desk *Desk) ToString() string {
	return fmt.Sprintf("%d %d %s", desk.Id, desk.TotalRevenue, desk.TotalUsedTime.Format(TimeFormat))
}
