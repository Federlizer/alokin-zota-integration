package main

import (
	"os"

	"github.com/federlizer/alokin-zota-integration/api"
	"github.com/federlizer/alokin-zota-integration/internal/storage"
	"github.com/federlizer/alokin-zota-integration/zota"
)

func main() {
	zotaSecretKey := os.Getenv("ZOTA_SECRET_KEY")
	zotaEndpointId := os.Getenv("ZOTA_ENDPOINT_ID")
	zotaMerchantId := os.Getenv("ZOTA_MERCHANT_ID")
	zotaBaseUrl := os.Getenv("ZOTA_BASE_URL")

	zotaApi := zota.NewZotaAPI(
		zotaSecretKey,
		zotaEndpointId,
		zotaMerchantId,
		zotaBaseUrl,
	)

	orderRepo := storage.NewOrderRepo()

	engine := api.SetupApi(zotaApi, *orderRepo)
	addr := ":8080"
	err := engine.Run(addr)
	if err != nil {
		panic(err)
	}
}
