package zota

import (
	"crypto/sha256"
	"fmt"
)

type OrderStatus string

const (
	// Order was created
	Created OrderStatus = "CREATED"
	// Order is being processed, continue polling until final status
	Processing = "PROCESSING"
	// Order is approved, final status
	Approved = "APPROVED"
	// Order is declined, final status
	Declined = "DECLINED"
	// Order is declined by fraud-prevention system.
	Filtered = "FILTERED"
	// Order is still in processing state, awaiting next step.
	Pending = "PENDING"
	// Order status is unknown, please inform Zota support team. Not a final status.
	Unknown = "UNKNOWN"
	// Order is declined due to a technical error, please inform Zota support team.
	// Final status.
	Error = "ERROR"
)

var FinalStatuses = []OrderStatus{
	Approved,
	Declined,
	Filtered,
	Error,
}

type ZotaOrderStatusRequest struct {
	OrderId         string `json:"orderID"`
	MerchantOrderId string `json:"merchantOrderID"`
	Timestamp       int64  `json:"timestamp"`
	Signature       string `json:"signature"`
}

func NewZotaOrderStatusRequest(orderId, merchantOrderId string) *ZotaOrderStatusRequest {
	return &ZotaOrderStatusRequest{
		OrderId:         orderId,
		MerchantOrderId: merchantOrderId,
		Timestamp:       0,
		Signature:       "",
	}
}

// GenSignature generates the signature required for the Zota
// Order Status request and returns it.
//
// Every request must be signed by the merchant in order to
// be successfully authenticated by Zotapay servers. The signature
// parameter of an Order Status Request must be generated by hashing
// a string of concatenated parameters using SHA-256 in the exact following order:
//
// MerchantID + merchantOrderID + orderID + timestamp + MerchantSecretKey
func (zosb *ZotaOrderStatusRequest) GenSignature(merchantId, secretKey string) string {
	signatureStr := fmt.Sprintf(
		"%s%s%s%d%s",
		merchantId,
		zosb.MerchantOrderId,
		zosb.OrderId,
		zosb.Timestamp,
		secretKey,
	)

	signature := fmt.Sprintf("%x", sha256.Sum256([]byte(signatureStr)))

	return signature
}

type ZotaOrderStatusResponse struct {
	Code    string  `json:"code"`
	Message *string `json:"message"`
	Data    *struct {
		Type                   string      `json:"type"`
		Status                 OrderStatus `json:"status"`
		ErrorMessage           string      `json:"errorMessage"`
		ProcessorTransactionId string      `json:"processorTransactionID"`
		OrderId                string      `json:"orderID"`
		MerchantOrderId        string      `json:"merchantOrderID"`
		Amount                 string      `json:"amount"`
		Currency               string      `json:"currency"`
		CustomerEmail          string      `json:"customerEmail"`
		// CustomParam interface{} `json:"customParam"`
		// ExtraData interface{} `json:"extraData"`
		// Request interface{} `json:"request"`
	} `json:"data"`
}

func (zosr *ZotaOrderStatusResponse) IsInFinalStatus() bool {
	var inFinalStatus = false

	// We're not good if we don't have data...
	if zosr.Data == nil {
		return inFinalStatus
	}

	for _, finalStatus := range FinalStatuses {
		if zosr.Data.Status == finalStatus {
			inFinalStatus = true
			break
		}
	}

	return inFinalStatus
}
