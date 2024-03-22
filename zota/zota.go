package zota

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ZotaAPI struct {
	SecretKey  string
	EndpointId string
	MerchantId string
	BaseUrl    string
}

func (api *ZotaAPI) Deposit(request *ZotaDepositRequest) (*ZotaDepositResponse, error) {
	endpointUrl := fmt.Sprintf("/api/v1/deposit/request/%s/", api.EndpointId)
	url := fmt.Sprintf("%s%s", api.BaseUrl, endpointUrl)

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	fmt.Printf("url: %s\n", url)
	fmt.Printf("body: %v\n", string(jsonBody))

	response, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %s\n", string(responseBody))
	zotaDepositResponse := ZotaDepositResponse{}
	err = json.Unmarshal(responseBody, &zotaDepositResponse)
	if err != nil {
		return nil, err
	}

	// Handle failed deposits
	if zotaDepositResponse.Code != "200" {
		errorMsg := ""
		// The zota API has returned an error to us, report that same error up the chain
		if zotaDepositResponse.Message != nil {
			errorMsg = "Received non-OK response from Zota API with no error message"
		} else {
			errorMsg = fmt.Sprintf("Received non-OK response from Zota API: %s", *zotaDepositResponse.Message)
		}

		return nil, errors.New(errorMsg)
	}

	fmt.Printf("Unmarshalled response: %v\n", zotaDepositResponse)

	return &zotaDepositResponse, nil
}

func (api *ZotaAPI) OrderStatus(request *ZotaOrderStatusRequest) (*ZotaOrderStatusResponse, error) {
	// First ensure the timestamp and signautre are correct
	ts := time.Now().Unix()
	request.Timestamp = ts
	request.Signature = request.GenSignature(api.MerchantId, api.SecretKey)

	endpointUrl := "/api/v1/query/order-status/"
	params := fmt.Sprintf(
		"?merchantID=%s&orderID=%s&merchantOrderID=%s&timestamp=%d&signature=%s",
		api.MerchantId,
		request.OrderId,
		request.MerchantOrderId,
		request.Timestamp,
		request.Signature,
	)
	url := fmt.Sprintf("%s%s%s", api.BaseUrl, endpointUrl, params)

	fmt.Printf("url: %s\n", url)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Response: %s\n", string(responseBody))
	zotaOrderStatusResponse := ZotaOrderStatusResponse{}
	err = json.Unmarshal(responseBody, &zotaOrderStatusResponse)
	if err != nil {
		return nil, err
	}

	return &zotaOrderStatusResponse, nil
}

// PollOrderStatus will continuously poll the Zota API server for an order status.
// This function can be called as an external goroutine.
func (api *ZotaAPI) PollOrderStatus(request *ZotaOrderStatusRequest) (*ZotaOrderStatusResponse, error) {
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
				close(quitChan)
				continue
			}

			zosr, err := api.OrderStatus(request)
			if err != nil {
				fmt.Println("Received error when checking order status: %v", err)
				continue
			}

			if zosr.Code != "200" {
				fmt.Printf("Received non-OK response from Zota: %v", zosr.Message)
				continue
			}

			if zosr.Data == nil {
				fmt.Printf("Received OK response from Zota, but no data field: %v", zosr)
				continue
			}

			if zosr.IsInFinalStatus() {
				fmt.Printf("We've received a final status for order %v\n", zosr.Data.MerchantOrderId)
				close(quitChan)
				return zosr, nil
			}

		case <-quitChan:
			ticker.Stop()
			return nil, errors.New("Couldn't get response ")
		}
	}
}
