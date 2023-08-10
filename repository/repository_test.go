package repository

import (
	"github.com/Mansur51-hub/ClubHandler/model"
	"testing"
)

func TestRepository(t *testing.T) {
	rep := NewRepository(5)

	_, err := rep.GetDesk(12)

	if err == nil {
		t.Errorf("desk not exists")
	}

	_, err = rep.GetDesk(3)

	if err != nil {
		t.Errorf("desk exists")
	}

	_, err = rep.DeskOccupied(12)

	if err == nil {
		t.Error("desk not exists")
	}

	client := model.Client{Id: "client1", TookDesk: false, DeskId: 0}
	rep.CreateClient(client.Id, &client)
	got, _ := rep.GetClient(client.Id)

	if *got != client {
		t.Errorf("client should be the same")
	}

	rep.DeleteClient(client.Id)

	_, err = rep.GetClient(client.Id)

	if err == nil {
		t.Errorf("client was delted")
	}
}
