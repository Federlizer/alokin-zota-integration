package internal

import (
	"strconv"

	"github.com/google/uuid"
)

type Order struct {
	Id     uuid.UUID
	Amount float64
	User   User
}

func (o *Order) AmountStr() string {
	// Need to do it like this if we want to ommit the zeroes and allow
	// amounts that have more than 2 decimals after floating point
	return strconv.FormatFloat(o.Amount, 'f', -1, 64)
}
