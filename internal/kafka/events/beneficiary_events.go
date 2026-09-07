package events

type BeneficiaryCreatedEvent struct {
	AccountID     string `json:"account_id"`
	BeneficiaryID string `json:"beneficiary_id"`
	NickName      string `json:"nick_name,omitempty"`

	CommonEvent
}
