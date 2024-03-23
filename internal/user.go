package internal

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
