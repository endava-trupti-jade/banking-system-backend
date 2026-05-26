package events

import "banking-system-backend/constants"

type NomineeCreatedEvent struct {
	AccountID         string                 `json:"account_id"`
	NomineeID         string                 `json:"nominee_id"`
	Relation          constants.RelationType `json:"relation"`
	IsPrimary         bool                   `json:"is_primary"`
	NomineePercentage int64                  `json:"nominee_percentage"`

	CommonEvent
}
