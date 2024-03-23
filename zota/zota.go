package zota

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/federlizer/alokin-zota-integration/internal"
)

type IZotaAPI interface {
    SecretKey() string
    EndpointId() string
    MerchantId() string
    BaseUrl() string

    Deposit(request *ZotaDepositRequest) (*ZotaDepositResponse, error)
    OrderStatus(request *ZotaOrderStatusRequest) (*ZotaOrderStatusResponse, error)
    PollOrderStatus(request *ZotaOrderStatusRequest, order *internal.Order) (*ZotaOrderStatusResponse, error)
}

type ZotaAPI struct {
	secretKey  string
	endpointId string
	merchantId string
	baseUrl    string
}

func NewZotaAPI(secretKey, endpointId, merchantId, baseUrl string) *ZotaAPI {
    return &ZotaAPI {
        secretKey,
        endpointId,
        merchantId,
        baseUrl,
    }
}

func (api *ZotaAPI) SecretKey() string {
    return api.secretKey
}

func (api *ZotaAPI) EndpointId() string {
    return api.endpointId
}

func (api *ZotaAPI) MerchantId() string {
    return api.merchantId
}

func (api *ZotaAPI) BaseUrl() string {
    return api.baseUrl
}

func (api *ZotaAPI) Deposit(request *ZotaDepositRequest) (*ZotaDepositResponse, error) {
	endpointUrl := fmt.Sprintf("/api/v1/deposit/request/%s/", api.EndpointId())
	url := fmt.Sprintf("%s%s", api.BaseUrl(), endpointUrl)

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	zotaDepositResponse := ZotaDepositResponse{}
	err = json.Unmarshal(responseBody, &zotaDepositResponse)
	if err != nil {
		return nil, err
	}

	// Handle failed deposits
	if zotaDepositResponse.Code != "200" {
		errorMsg := ""
		// The zota API has returned an error to us, report that same error up the chain
		if zotaDepositResponse.Message == nil {
			errorMsg = "Received non-OK response from Zota API with no error message"
		} else {
			errorMsg = fmt.Sprintf("Received non-OK response from Zota API: %s", *zotaDepositResponse.Message)
		}

		return nil, errors.New(errorMsg)
	}

	return &zotaDepositResponse, nil
}

func (api *ZotaAPI) OrderStatus(request *ZotaOrderStatusRequest) (*ZotaOrderStatusResponse, error) {
	// First ensure the timestamp and signautre are correct
	ts := time.Now().Unix()
	request.Timestamp = ts
	request.Signature = request.GenSignature(api.MerchantId(), api.SecretKey())

	endpointUrl := "/api/v1/query/order-status/"
	params := fmt.Sprintf(
		"?merchantID=%s&orderID=%s&merchantOrderID=%s&timestamp=%d&signature=%s",
		api.MerchantId(),
		request.OrderId,
		request.MerchantOrderId,
		request.Timestamp,
		request.Signature,
	)
	url := fmt.Sprintf("%s%s%s", api.BaseUrl(), endpointUrl, params)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	zotaOrderStatusResponse := ZotaOrderStatusResponse{}
	err = json.Unmarshal(responseBody, &zotaOrderStatusResponse)
	if err != nil {
		return nil, err
	}

	return &zotaOrderStatusResponse, nil
}

// PollOrderStatus will continuously poll the Zota API server for an order status.
// This function is intended to work as a goroutine.
func (api *ZotaAPI) PollOrderStatus(request *ZotaOrderStatusRequest, order *internal.Order) (*ZotaOrderStatusResponse, error) {
	ticker := time.NewTicker(10 * time.Second)
	quitChan := make(chan bool)
	maxAttempts := 20
	attempts := 0

	for {
		select {
		case <-ticker.C:
			// Stop querying after we've reached max attempts
			if attempts >= maxAttempts {
				fmt.Printf("Reached maximum retry attempts (%v) before receiving a final order status\n", maxAttempts)
				order.PaymentStatus = internal.PaymentStatusFailed
				close(quitChan)
				continue
			}

			zosr, err := api.OrderStatus(request)
			if err != nil {
				fmt.Printf("Received error when checking order status: %v\n", err)
				continue
			}

			if zosr.Code != "200" {
				fmt.Printf("Received non-OK response from Zota: %v\n", zosr.Message)
				continue
			}

			if zosr.Data == nil {
				fmt.Printf("Received OK response from Zota, but no data field: %v\n", zosr)
				continue
			}

			if zosr.IsInFinalStatus() {
				fmt.Printf("We've received a final status for order %v\n", zosr.Data.MerchantOrderId)

				// Set payment status
				if zosr.Data.Status == Approved {
					order.PaymentStatus = internal.PaymentStatusApproved
				} else {
					order.PaymentStatus = internal.PaymentStatusFailed
				}

				close(quitChan)
				return zosr, nil
			}

		case <-quitChan:
			ticker.Stop()
			return nil, errors.New("Couldn't get response")
		}
	}
}
