package main

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/federlizer/alokin-zota-integration/api"
	"github.com/federlizer/alokin-zota-integration/zota"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	zotaSecretKey := os.Getenv("ZOTA_SECRET_KEY")
	zotaEndpointId := os.Getenv("ZOTA_ENDPOINT_ID")
	zotaMerchantId := os.Getenv("ZOTA_MERCHANT_ID")
	zotaBaseUrl := os.Getenv("ZOTA_BASE_URL")

	zotaApi := zota.ZotaAPI{
		SecretKey:  zotaSecretKey,
		EndpointId: zotaEndpointId,
		MerchantId: zotaMerchantId,
		BaseUrl:    zotaBaseUrl,
	}

	engine := api.SetupApi(zotaApi)
	err = engine.Run()
	if err != nil {
		panic(err)
	}
}
