package parser

import (
	"fmt"
	"github.com/Mansur51-hub/ClubHandler/model"
	"strconv"
	"strings"
	"time"
)

const (
	lineDefaultLength = 3
	lineMaxLength     = 4
)

func ParseLine(line string) (model.Event, error) {
	parts := strings.Fields(line)
	event := model.Event{}

	if !(len(parts) >= lineDefaultLength && len(parts) <= lineMaxLength) {
		err := fmt.Errorf("invalid event format: expected words count in range %d-%d", lineDefaultLength, lineMaxLength)
		return event, err
	}

	if moment, err := time.Parse(model.TimeFormat, parts[0]); err != nil {
		return event, err
	} else {
		event.Moment = moment
	}

	if id, err := strconv.Atoi(parts[1]); err != nil {
		return event, err
	} else {
		event.Id = model.EventId(id)
	}

	event.ClientId = parts[2]

	if len(parts) == lineMaxLength {
		if event.Id != model.ClientTookDeskEventId {
			err := fmt.Errorf("invalid event format: expected %d words, got: %d", lineDefaultLength, lineMaxLength)
			return event, err
		} else {
			if descId, err := strconv.ParseUint(parts[3], 10, 32); err != nil {
				return event, err
			} else {
				id := uint(descId)
				event.DescId = id
			}
		}
	}

	return event, nil
}
