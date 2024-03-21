package internal

import "github.com/google/uuid"

type User struct {
	Email     string
	FirstName string
	LastName  string
	IpAddr    string
	Phone     string
	Address   UserAddress
}

type UserAddress struct {
	AddressLine string
	CountryCode string
	City        string
	// state string
	ZipCode string
}

func (u *User) PlaceOrder(amount float64) *Order {
	orderId := uuid.New()

	return &Order{
		Id:     orderId,
		Amount: amount,
		User:   *u,
	}
}
