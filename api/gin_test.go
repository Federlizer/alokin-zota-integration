package api

import (
	"net/http"
	"net/http/httptest"
    "testing"

	"github.com/federlizer/alokin-zota-integration/zota"
	"github.com/federlizer/alokin-zota-integration/internal"
	"github.com/federlizer/alokin-zota-integration/internal/storage"
)

type zotaAPIMock struct {}

func (api *zotaAPIMock) SecretKey() string { return "00000000-1111-2222-3333-444444444444" }
func (api *zotaAPIMock) EndpointId() string { return "123456" }
func (api *zotaAPIMock) MerchantId() string { return "COOKIES1337" }
func (api *zotaAPIMock) BaseUrl() string { return "https://federlizer.com/api/" }

func (api *zotaAPIMock) Deposit(req *zota.ZotaDepositRequest) (*zota.ZotaDepositResponse, error) { return nil, nil }
func (api *zotaAPIMock) OrderStatus(req *zota.ZotaOrderStatusRequest) (*zota.ZotaOrderStatusResponse, error) { return nil, nil }
func (api *zotaAPIMock) PollOrderStatus(req *zota.ZotaOrderStatusRequest, order *internal.Order) (*zota.ZotaOrderStatusResponse, error) { return nil, nil }

func createZotaAPIMock() *zotaAPIMock {
    return &zotaAPIMock{}
}

func createOrderRepo() *storage.OrderRepo {
    return storage.NewOrderRepo()
}

func TestPingEndpoint(t *testing.T) {
    engine := SetupApi(createZotaAPIMock(), *createOrderRepo())

    resWriter := httptest.NewRecorder()
    req, err := http.NewRequest("GET", "/ping", nil)
    if err != nil {
        t.Errorf("Failed to init request: %q\n", err)
    }

    engine.ServeHTTP(resWriter, req)

    expectedCode := 200
    if resWriter.Code != expectedCode {
        t.Errorf("Server response %q doesn't equal expected %q", resWriter.Code, expectedCode)
    }

    expectedBody := "{\"message\":\"pong\"}"
    if resWriter.Body.String() != expectedBody {
        t.Errorf("Server response %q doesn't equal expected %q", resWriter.Body.String(), expectedBody)
    }
}
