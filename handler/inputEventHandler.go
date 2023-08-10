package handler

import (
	"bufio"
	"fmt"
	"github.com/Mansur51-hub/ClubHandler/model"
	"github.com/Mansur51-hub/ClubHandler/parser"
	"os"
	"strconv"
	"strings"
	"time"
)

func HandleInput(path string) (*model.InputMetaData, error) {

	data := model.NewInputMetaData()

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return data, err
	}

	scanner := bufio.NewScanner(file)

	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) < 3 {
		err := fmt.Errorf("invalid input data. Expected deskNumber, open(close)moment, events number param at least")
		return data, err
	}

	if deskNumber, err := strconv.ParseUint(lines[0], 10, 32); err != nil {
		return data, err
	} else {
		number := uint(deskNumber)
		data.Cafe.DescNumber = number
	}

	times := strings.Fields(lines[1])

	if len(times) != 2 {
		err := fmt.Errorf("invalid input data. Expected two params in open(close)moment line")
		return data, err
	}

	if openMoment, err := time.Parse(model.TimeFormat, times[0]); err != nil {
		return data, err
	} else {
		data.Cafe.OpenMoment = openMoment
	}

	if closureMoment, err := time.Parse(model.TimeFormat, times[1]); err != nil {
		return data, err
	} else {
		data.Cafe.ClosureMoment = closureMoment
	}

	if costPerHour, err := strconv.ParseUint(lines[2], 10, 32); err != nil {
		return data, err
	} else {
		cost := uint(costPerHour)
		data.Cafe.CostPerHour = cost
	}

	for i := 3; i < len(lines); i++ {
		if event, err := parser.ParseLine(lines[i]); err != nil {
			err := fmt.Errorf("invalid input data. At line %d: %w", i+1, err)
			return data, err
		} else {
			data.Events = append(data.Events, event)
		}
	}

	return data, nil
}
