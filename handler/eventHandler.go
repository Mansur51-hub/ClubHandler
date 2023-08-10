package handler

import (
	"errors"
	"fmt"
	"github.com/Mansur51-hub/ClubHandler/model"
	"github.com/Mansur51-hub/ClubHandler/repository"
	"time"
)

type Handler struct {
	repository *repository.Repository
	data       *model.InputMetaData
}

func (h *Handler) HandleEvents() ([]string, error) {
	if !checkTimesOrder(h.data.Events) {
		return nil, errors.New("input data error. The times are not in ascending order")
	}

	report := make([]string, 0, len(h.data.Events)*2+int(h.data.Cafe.DescNumber)+2)
	report = append(report, h.data.Cafe.OpenMoment.Format(model.TimeFormat))

	for _, event := range h.data.Events {
		report = append(report, event.ToString())
		if r, err := h.handleEvent(&event); err != nil {
			errMessage := fmt.Sprintf("%s %d %s", event.Moment.Format(model.TimeFormat), model.ErrorOutcomeEventId, err)
			report = append(report, errMessage)
		} else {
			if r != "" {
				report = append(report, r)
			}
		}
	}

	remainingClients := h.repository.GetRemainingClients()

	for _, val := range remainingClients {
		r := h.handleClientLeaveAtClosureEvent(&val)
		report = append(report, r)
	}

	report = append(report, h.data.Cafe.ClosureMoment.Format(model.TimeFormat))

	desks := h.repository.GetDesks()

	for _, val := range desks {
		report = append(report, val.ToString())
	}

	return report, nil
}

func NewHandler(data *model.InputMetaData) *Handler {
	return &Handler{repository: repository.NewRepository(data.Cafe.DescNumber), data: data}
}

func (h *Handler) handleEvent(event *model.Event) (string, error) {
	switch event.Id {
	case model.ClientIncomeEventId:
		return h.handleClientIncomeEvent(event)
	case model.ClientTookDeskEventId:
		return h.handleClientTookDeskEvent(event)
	case model.ClientWaitEventId:
		return h.handleClientWaitEvent(event)
	case model.ClientLeaveEventId:
		return h.handleClientLeaveEvent(event)
	default:
		return "", fmt.Errorf("unable to handle event with id %d", event.Id)
	}
}

func (h *Handler) handleClientIncomeEvent(event *model.Event) (string, error) {
	if h.repository.ClientExists(event.ClientId) {
		return "", errors.New("YouShallNotPass")
	}
	if cafeClosed(h.data.Cafe.OpenMoment, h.data.Cafe.ClosureMoment, event.Moment) {
		return "", errors.New("NotOpenYet")
	}
	client := model.Client{Id: event.ClientId}
	h.repository.CreateClient(client.Id, &client)

	return "", nil
}

func (h *Handler) handleClientTookDeskEvent(event *model.Event) (string, error) {
	if !h.repository.ClientExists(event.ClientId) {
		return "", errors.New("ClientUnknown")
	}

	if occupied, err := h.repository.DeskOccupied(event.DescId); err != nil {
		return "", err
	} else {
		if occupied {
			return "", errors.New("PlaceIsBusy")
		}

		client, _ := h.repository.GetClient(event.ClientId)

		if client.TookDesk {
			desk, _ := h.repository.GetDesk(client.DeskId)
			desk.HandleClientLeaveAction(event.Moment, h.data.Cafe.CostPerHour)
			h.repository.UpdateDesk(desk.Id, desk)
		}

		newDesk, _ := h.repository.GetDesk(event.DescId)

		h.handleClientTookDeskAction(client, newDesk, event.Moment)

		return "", nil
	}
}

func (h *Handler) handleClientWaitEvent(event *model.Event) (string, error) {
	if h.repository.HasFreeDesk() {
		return "", errors.New("ICanWaitNoLonger")
	}

	if h.repository.Queue.Size() >= int(h.data.Cafe.DescNumber) {
		return fmt.Sprintf("%s %d %s", event.Moment.Format(model.TimeFormat), model.ClientLeaveOutcomeEventId, event.ClientId), nil
	}

	h.repository.Queue.Push(model.Client{Id: event.ClientId})

	return "", nil
}

func (h *Handler) handleClientLeaveEvent(event *model.Event) (string, error) {
	if client, err := h.repository.GetClient(event.ClientId); err != nil {
		return "", errors.New("ClientUnknown")
	} else {
		defer h.repository.DeleteClient(client.Id)
		if client.TookDesk {
			desk, _ := h.repository.GetDesk(client.DeskId)
			desk.HandleClientLeaveAction(event.Moment, h.data.Cafe.CostPerHour)

			if !h.repository.Queue.Empty() {
				newClient := h.repository.Queue.Front()
				h.repository.Queue.Pop()

				h.handleClientTookDeskAction(&newClient, desk, event.Moment)

				return fmt.Sprintf("%s %d %s %d",
					event.Moment.Format(model.TimeFormat),
					model.ClientTookDeskOutcomeEventId,
					newClient.Id,
					desk.Id), nil
			}

			h.repository.UpdateDesk(desk.Id, desk)

		}

		return "", nil

	}
}

func (h *Handler) handleClientLeaveAtClosureEvent(client *model.Client) string {
	defer h.repository.DeleteClient(client.Id)

	if client.TookDesk {
		desk, _ := h.repository.GetDesk(client.DeskId)
		desk.HandleClientLeaveAction(h.data.Cafe.ClosureMoment, h.data.Cafe.CostPerHour)
		h.repository.UpdateDesk(desk.Id, desk)
	}

	return fmt.Sprintf("%s %d %s", h.data.Cafe.ClosureMoment.Format(model.TimeFormat), model.ClientLeaveOutcomeEventId, client.Id)
}

func (h *Handler) handleClientTookDeskAction(client *model.Client, desk *model.Desk, time time.Time) {
	client.TookDesk = true
	client.DeskId = desk.Id
	desk.Occupied = true
	desk.LastUsageStartMoment = time

	h.repository.UpdateClient(client.Id, client)
	h.repository.UpdateDesk(desk.Id, desk)
}

func cafeClosed(openTime, closureTime, time time.Time) bool {
	return time.Before(openTime) || closureTime.Before(time)
}

func checkTimesOrder(events []model.Event) bool {
	if len(events) < 2 {
		return true
	}

	moment := events[0].Moment

	for _, val := range events {
		if val.Moment.Before(moment) {
			return false
		}

		moment = val.Moment
	}

	return true
}
