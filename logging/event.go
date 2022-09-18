package logging

type EventType byte

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

const (
	EventDelete EventType = iota + 1
	EventPut
)
