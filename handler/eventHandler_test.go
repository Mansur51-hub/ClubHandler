package handler

import (
	"fmt"
	"github.com/Mansur51-hub/ClubHandler/model"
	"reflect"
	"testing"
	"time"
)

func TestHandleEvent(t *testing.T) {
	time1, _ := time.Parse(model.TimeFormat, "10:00")
	time2, _ := time.Parse(model.TimeFormat, "10:01")
	time3, _ := time.Parse(model.TimeFormat, "14:04")

	data := model.NewInputMetaData()
	data.Cafe.DescNumber = 2
	data.Cafe.CostPerHour = 1
	data.Cafe.OpenMoment = time1
	data.Cafe.ClosureMoment = time3

	clientId1, clientId2 := "client1", "client2"

	events := []model.Event{
		{
			Moment:   time1,
			Id:       model.ClientIncomeEventId,
			ClientId: clientId1,
			DescId:   0,
		},
		{
			Moment:   time2,
			Id:       model.ClientTookDeskEventId,
			ClientId: clientId1,
			DescId:   2,
		},
	}

	handler := NewHandler(data)

	for _, e := range events {
		_, _ = handler.handleEvent(&e)
	}

	var tests = []struct {
		name  string
		event *model.Event
		want  string
	}{{
		name: "should be ClientUnknown error",
		event: &model.Event{
			Id:       model.ClientTookDeskEventId,
			ClientId: clientId2,
			DescId:   2,
			Moment:   time2,
		},
		want: "ClientUnknown",
	}, {
		name: "should be YouShallNotPass error",
		event: &model.Event{
			Id:       model.ClientIncomeEventId,
			ClientId: clientId1,
			DescId:   0,
			Moment:   time2,
		},
		want: "YouShallNotPass",
	}, {
		name: "should be NotOpenYet error",
		event: &model.Event{
			Id:       model.ClientIncomeEventId,
			ClientId: clientId2,
			DescId:   0,
			Moment:   time3.Add(time.Minute),
		},
		want: "NotOpenYet",
	}, {
		name: "should be PlaceIsBusy error",
		event: &model.Event{
			Id:       model.ClientTookDeskEventId,
			ClientId: clientId1,
			DescId:   2,
			Moment:   time2,
		},
		want: "PlaceIsBusy",
	}, {
		name: "should be ICanWaitNoLonger error",
		event: &model.Event{
			Id:       model.ClientWaitEventId,
			ClientId: clientId1,
			DescId:   0,
			Moment:   time3,
		},
		want: "ICanWaitNoLonger",
	}, {
		name: "should be ClientUnknown leave error",
		event: &model.Event{
			Id:       model.ClientLeaveEventId,
			ClientId: clientId2,
			DescId:   0,
			Moment:   time3,
		},
		want: "ClientUnknown",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.handleEvent(tt.event)

			if err == nil || err.Error() != tt.want {
				t.Errorf("error got: %s, want: %s", err, tt.want)
			}
		})
	}
}

func TestHandleInput(t *testing.T) {
	time1, _ := time.Parse(model.TimeFormat, "10:00")
	time2, _ := time.Parse(model.TimeFormat, "10:01")
	time3, _ := time.Parse(model.TimeFormat, "10:02")

	data := model.NewInputMetaData()
	data.Cafe.DescNumber = 1
	data.Cafe.CostPerHour = 1
	data.Cafe.OpenMoment = time1
	data.Cafe.ClosureMoment = time3

	clientId1, clientId2 := "client1", "client2"

	events := []model.Event{
		{
			Moment:   time1,
			Id:       model.ClientIncomeEventId,
			ClientId: clientId1,
			DescId:   0,
		},
		{
			Moment:   time1,
			Id:       model.ClientTookDeskEventId,
			ClientId: clientId1,
			DescId:   1,
		},
		{
			Moment:   time1,
			Id:       model.ClientIncomeEventId,
			ClientId: clientId2,
			DescId:   0,
		},
		{
			Moment:   time1,
			Id:       model.ClientWaitEventId,
			ClientId: clientId2,
			DescId:   0,
		},
		{
			Moment:   time2,
			Id:       model.ClientLeaveEventId,
			ClientId: clientId1,
			DescId:   0,
		},
	}

	data.Events = events

	handler := NewHandler(data)

	want := []string{
		time1.Format(model.TimeFormat),
		events[0].ToString(),
		events[1].ToString(),
		events[2].ToString(),
		events[3].ToString(),
		events[4].ToString(),
	}

	deskId := 1

	want = append(want, fmt.Sprintf("%s %d %s %d",
		time2.Format(model.TimeFormat),
		model.ClientTookDeskOutcomeEventId,
		clientId2,
		deskId))

	want = append(want, fmt.Sprintf("%s %d %s",
		time3.Format(model.TimeFormat),
		model.ClientLeaveOutcomeEventId,
		clientId2))

	want = append(want, time3.Format(model.TimeFormat))

	totalUsedTime, _ := time.Parse(model.TimeFormat, "00:02")
	totalRevenue := 2

	want = append(want, fmt.Sprintf("%d %d %s", deskId, totalRevenue, totalUsedTime.Format(model.TimeFormat)))

	report, _ := handler.HandleEvents()

	if !reflect.DeepEqual(report, want) {
		t.Errorf("got %s\n want %s", report, want)
	}
}
