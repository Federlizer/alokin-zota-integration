package main

import (
	"fmt"
    "os"

    "github.com/joho/godotenv"

    "github.com/federlizer/alokin-zota-integration/zota"
    "github.com/federlizer/alokin-zota-integration/internal"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        panic(err)
    }

    zotaSecretKey := os.Getenv("ZOTA_SECRET_KEY")
    zotaEndpointId := os.Getenv("ZOTA_ENDPOINT_ID")
    zotaBaseUrl := os.Getenv("ZOTA_BASE_URL")

    zotaApi := zota.ZotaAPI {
        SecretKey: zotaSecretKey,
        EndpointId: zotaEndpointId,
        BaseUrl: zotaBaseUrl,
    }

    // User creation
	userAddress := internal.UserAddress{
		AddressLine: "My lovely home address line",
		CountryCode: "DK",
		City:        "Aalborg",
		ZipCode:     "9000",
	}

	user := internal.User{
		Email:     "federlizer@protonmail.com",
		FirstName: "Nikola",
		LastName:  "Velichkov",
		IpAddr:    "146.70.188.231",
		Phone:     "+4550331329",
		Address:   userAddress,
	}

	order := user.PlaceOrder(13.37)

    zotaDepositRequest := zota.FromOrder(order, zotaApi.EndpointId, zotaApi.SecretKey)

    response, err := zotaApi.Deposit(zotaDepositRequest)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Successfully started deposit flow for order %s\n", response.Data.OrderId)
}
