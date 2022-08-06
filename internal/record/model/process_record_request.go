package model

import (
	"github.com/go-openapi/strfmt"
)

type ProcessRecordRequest struct {
	State         State        `json:"state" validate:"required,state"`
	Amount        float64      `json:"amount" validate:"required,min=0"`
	TransactionID strfmt.UUID4 `json:"transaction_id" validate:"required,uuid4"`
	UserID        strfmt.UUID4 `json:"user_id" validate:"required,uuid4"`
}
