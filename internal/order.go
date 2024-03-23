package internal

import (
	"strconv"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "PENDING"
	PaymentStatusApproved               = "APPROVED"
	PaymentStatusFailed                 = "FAILED"
)

type Order struct {
	Id            uuid.UUID     `json:"id"`
	Description   string        `json:"description"`
	Amount        float64       `json:"amount"`
	User          User          `json:"-"`
	PaymentStatus PaymentStatus `json:"paymentStatus"`
}

func NewOrder(user *User, amount float64, description string) *Order {
	orderId := uuid.New()

	return &Order{
		Id:            orderId,
		Amount:        amount,
		Description:   description,
		User:          *user,
		PaymentStatus: PaymentStatusPending,
	}
}

func (o *Order) AmountStr() string {
	// Need to do it like this if we want to ommit the zeroes and allow
	// amounts that have more than 2 decimals after floating point
	return strconv.FormatFloat(o.Amount, 'f', -1, 64)
}
