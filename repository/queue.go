package repository

import (
	"github.com/Mansur51-hub/ClubHandler/model"
)

type Queue struct {
	arr []model.Client
}

func (q *Queue) Front() model.Client {
	return q.arr[0]
}

func (q *Queue) Empty() bool {
	return len(q.arr) == 0
}

func (q *Queue) Pop() {
	q.arr = q.arr[1:]
}

func (q *Queue) Size() int {
	return len(q.arr)
}

func (q *Queue) Push(client model.Client) {
	q.arr = append(q.arr, client)
}

func NewQueue(deskNumber uint) *Queue {
	return &Queue{arr: make([]model.Client, 0, deskNumber)}
}
