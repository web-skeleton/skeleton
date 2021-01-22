package event

import "time"

// SystemUpDownEvent 系统启停事件
type SystemUpDownEvent struct {
	Up        bool
	CreatedAt time.Time
	Sync      bool
}

func (evt SystemUpDownEvent) Async() bool {
	return !evt.Sync
}

func (evt SystemUpDownEvent) EvtType() string {
	upOrDown := "down"
	if evt.Up {
		upOrDown = "up"
	}

	return upOrDown
}