package event

type eventType string

const (
	stop    eventType = "stop"
	pause   eventType = "pause"
	resume  eventType = "resume"
	restart eventType = "restart"
)

type Event struct {
	T eventType
}

func (e Event) IsStop() bool {
	return e.T == stop
}

func (e Event) IsPause() bool {
	return e.T == pause
}

func (e Event) IsResume() bool {
	return e.T == resume
}

func (e Event) IsRestart() bool {
	return e.T == restart
}
