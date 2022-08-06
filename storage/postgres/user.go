package postgres

import (
	"github.com/go-openapi/strfmt"
)

type User struct {
	ID      strfmt.UUID4 `db:"id"`
	Balance float64      `db:"balance"`
}
