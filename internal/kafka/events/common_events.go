package events

import "time"

type CommonEvent struct {
	EventID    string    `json:"event_id"`
	EventType  string    `json:"event_type"`
	CreatedBy  string    `json:"created_by"`
	OccurredAt time.Time `json:"occurred_at"`
}
