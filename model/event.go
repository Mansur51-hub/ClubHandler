package model

import (
	"fmt"
	"time"
)

const TimeFormat = "15:04"

type EventId int

const (
	ClientIncomeEventId EventId = iota + 1
	ClientTookDeskEventId
	ClientWaitEventId
	ClientLeaveEventId
)

const (
	ClientLeaveOutcomeEventId EventId = iota + 11
	ClientTookDeskOutcomeEventId
	ErrorOutcomeEventId
)

type Event struct {
	Moment   time.Time
	Id       EventId
	ClientId string
	DescId   uint
}

func (e *Event) ToString() string {
	s := fmt.Sprintf("%s %d %s", e.Moment.Format(TimeFormat), e.Id, e.ClientId)
	if e.Id == ClientTookDeskEventId {
		return fmt.Sprintf("%s %d", s, e.DescId)
	}

	return s
}
