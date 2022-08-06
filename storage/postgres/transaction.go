package postgres

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/sylph4/entain-task/internal/record/model"
)

type Transaction struct {
	ID          strfmt.UUID4 `db:"id"`
	UserID      strfmt.UUID4 `db:"user_id"`
	State       model.State  `db:"state"`
	Amount      float64      `db:"amount"`
	ProcessedAt time.Time    `db:"processed_at"`
	IsCanceled  bool         `db:"is_canceled"`
}
