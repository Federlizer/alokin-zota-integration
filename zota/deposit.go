package zota

import (
	"crypto/sha256"
	"fmt"

	"github.com/federlizer/alokin-zota-integration/internal"
)

// ZotaDepositRequest represents the request body that's required by Zota's
// deposit request
type ZotaDepositRequest struct {
	MerchantOrderID   string `json:"merchantOrderId"`
	MerchantOrderDesc string `json:"merchantOrderDesc"`
	OrderAmount       string `json:"orderAmount"`
	OrderCurrency     string `json:"orderCurrency"`

	CustomerEmail     string `json:"customerEmail"`
	CustomerFirstName string `json:"customerFirstName"`
	CustomerLastName  string `json:"customerLastName"`
	CustomerIP        string `json:"customerIP"`
	CustomerPhone     string `json:"customerPhone"`

	CustomerAddress     string `json:"customerAddress"`
	CustomerCountryCode string `json:"customerCountryCode"`
	CustomerCity        string `json:"customerCity"`
	CustomerZipCode     string `json:"customerZipCode"`
	// CustomerState       string `json:"customerState"`

	RedirectUrl string `json:"redirectUrl"`
	CheckoutUrl string `json:"checkoutUrl"`
	Signature   string `json:"signature"`
}

// FromOrder creates a new ZotaDepositRequest struct based on the order, endpointId
// and secretKey passed. This function automatically generates a signature for the request
// and assigns it to the ZotaDepositRequest returned.
func FromOrder(order *internal.Order, endpointId, secretKey string) *ZotaDepositRequest {
	// Create request body
	zdr := ZotaDepositRequest{
		MerchantOrderID:   order.Id.String(),
		MerchantOrderDesc: order.Description,
		OrderAmount:       order.AmountStr(),
		OrderCurrency:     "USD",

		CustomerEmail:     order.User.Email,
		CustomerFirstName: order.User.FirstName,
		CustomerLastName:  order.User.LastName,
		CustomerIP:        order.User.IpAddr,
		CustomerPhone:     order.User.Phone,

		CustomerAddress:     order.User.Address.AddressLine,
		CustomerCountryCode: order.User.Address.CountryCode,
		CustomerCity:        order.User.Address.City,
		CustomerZipCode:     order.User.Address.ZipCode,

		RedirectUrl: "https://federlizer.com/deposit-completed",
		CheckoutUrl: "https://federlizer.com/checkout",
		Signature:   "",
	}

	// Generate request signature
	signature := zdr.GenSignature(endpointId, secretKey)
	zdr.Signature = signature

	return &zdr
}

// GenSignature generates the signature required for the
// Zota Deposit request and returns it
//
// Every request must be signed by the merchant in order to be
// successfully authenticated by Zotapay servers.
//
// EndpointID + merchantOrderID + orderAmount + customerEmail + MerchantSecretKey
func (zdr *ZotaDepositRequest) GenSignature(endpointId, secretKey string) string {
	signatureStr := fmt.Sprintf(
		"%s%s%s%s%s",
		endpointId,
		zdr.MerchantOrderID,
		zdr.OrderAmount,
		zdr.CustomerEmail,
		secretKey,
	)

	signature := fmt.Sprintf("%x", sha256.Sum256([]byte(signatureStr)))

	return signature
}

// ZotaDepositResponse represents the response that's received by Zota's
// deposit request
type ZotaDepositResponse struct {
	Code string `json:"code"`
	// Success data
	Data *struct {
		MerchantOrderID string `json:"merchantOrderID"`
		DepositUrl      string `json:"depositUrl"`
		OrderId         string `json:"orderID"`
	} `json:"data"`
	// Error message
	Message *string
}
