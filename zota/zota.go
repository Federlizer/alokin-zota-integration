package zota

import (
    "bytes"
    "errors"
    "fmt"
    "encoding/json"
    "net/http"
    "io"
)

type ZotaAPI struct {
	SecretKey  string
	EndpointId string
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
