package model

type InputMetaData struct {
	Cafe   Cafe
	Events []Event
}

func NewInputMetaData() *InputMetaData {
	return &InputMetaData{Cafe: Cafe{}, Events: make([]Event, 0)}
}
