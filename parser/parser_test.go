package parser

import (
	"github.com/Mansur51-hub/ClubHandler/model"
	"testing"
	"time"
)

func TestParseLine(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		// the table itself
		{"line length should be 3-4", "11:45 3 client4 12 132"},
		{"time format should be HH:mm", "12:33:00 4 client1"},
		{"desk number should be positive", "12:33 2 client1 -3"},
	}
	// The execution loop
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseLine(tt.input)
			if err == nil {
				t.Errorf("should be error thrown")
			}
		})
	}
}

func TestParseLineReturnModel(t *testing.T) {
	incomeTime := "09:48"
	moment, _ := time.Parse(model.TimeFormat, incomeTime)
	var tests = []struct {
		name  string
		input string
		want  model.Event
	}{
		// the table itself
		{"event with id 1",
			incomeTime + " 1 client2",
			model.Event{
				Id:       model.ClientIncomeEventId,
				ClientId: "client2",
				Moment:   moment,
				DescId:   0,
			}},
		{name: "event with id 4",
			input: incomeTime + " 4 client2",
			want: model.Event{
				Id:       model.ClientLeaveEventId,
				ClientId: "client2",
				Moment:   moment,
				DescId:   0,
			}},
		{name: "event with id 2",
			input: incomeTime + " 2 client2 2",
			want: model.Event{
				Id:       model.ClientTookDeskEventId,
				ClientId: "client2",
				Moment:   moment,
				DescId:   2,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, _ := ParseLine(tt.input)
			if ans != tt.want {
				t.Error(ans, tt.want)
			}
		})
	}
}
