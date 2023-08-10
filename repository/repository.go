package repository

import (
	"fmt"
	"github.com/Mansur51-hub/ClubHandler/model"
	"sort"
	"time"
)

type Repository struct {
	desk   map[uint]model.Desk
	client map[string]model.Client
	Queue  *Queue
}

func (r *Repository) ClientExists(id string) bool {
	_, found := r.client[id]
	return found
}

func (r *Repository) UpdateClient(id string, client *model.Client) {
	r.client[id] = *client
}

func (r *Repository) CreateClient(id string, client *model.Client) {
	r.client[id] = *client
}

func (r *Repository) DeleteClient(id string) {
	delete(r.client, id)
}

func (r *Repository) GetClient(id string) (*model.Client, error) {
	if client, found := r.client[id]; !found {
		return nil, fmt.Errorf("error: Client with id %s not found", id)
	} else {
		return &client, nil
	}
}

func (r *Repository) UpdateDesk(id uint, desk *model.Desk) {
	r.desk[id] = *desk
}

func (r *Repository) GetDesk(id uint) (*model.Desk, error) {
	if desk, found := r.desk[id]; !found {
		return nil, fmt.Errorf("error: Desk with id %d not found", id)
	} else {
		return &desk, nil
	}
}

func (r *Repository) DeskOccupied(id uint) (bool, error) {
	desk, found := r.desk[id]

	if !found {
		err := fmt.Errorf("error: Desk with id :%d not found", id)
		return false, err
	}

	return desk.Occupied, nil
}

func (r *Repository) HasFreeDesk() bool {
	for _, val := range r.desk {
		if !val.Occupied {
			return true
		}
	}

	return false
}

func (r *Repository) GetRemainingClients() []model.Client {
	var keys []string
	for k := range r.client {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	clients := make([]model.Client, 0, len(r.client))

	for _, k := range keys {
		clients = append(clients, r.client[k])
	}

	return clients
}

func (r *Repository) GetDesks() []model.Desk {
	desks := make([]model.Desk, 0, len(r.desk))

	for i := 1; i <= len(r.desk); i++ {
		desks = append(desks, r.desk[uint(i)])
	}

	return desks
}

func NewRepository(deskNumber uint) *Repository {
	rep := Repository{desk: make(map[uint]model.Desk), client: make(map[string]model.Client), Queue: NewQueue(deskNumber)}

	var zeroTime, _ = time.Parse(model.TimeFormat, "00:00")
	i := uint(1)
	for ; i <= deskNumber; i++ {
		rep.desk[i] = model.Desk{Id: i, Occupied: false, LastUsageStartMoment: zeroTime, TotalUsedTime: zeroTime, TotalRevenue: 0}
	}

	return &rep
}
