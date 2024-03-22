package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/federlizer/alokin-zota-integration/internal"
	"github.com/federlizer/alokin-zota-integration/zota"
)

func SetupApi(zotaApi zota.ZotaAPI) *gin.Engine {
	engine := gin.Default()

	engine.Use(func(c *gin.Context) {
		c.Set("zotaApi", zotaApi)
	})

	engine.GET("/ping", pingHandler)
	engine.GET("/order", orderHandler)

	return engine
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func orderHandler(c *gin.Context) {
	// TODO - is there a better way to have this?
	zotaApi := c.MustGet("zotaApi").(zota.ZotaAPI)

	// Init user (ideally, this would be somehow fetched from a DB
	// based on authentication credentials provided by the user)
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

	// Create order for that user (also ideally, this would use the DB to persist
	// the data)
	order := user.PlaceOrder(13.37)
	zotaDepositRequest := zota.FromOrder(order, zotaApi.EndpointId, zotaApi.SecretKey)

	// Make actual request
	response, err := zotaApi.Deposit(zotaDepositRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	zotaOrderStatusRequest := zota.NewZotaOrderStatusRequest(response.Data.OrderId, response.Data.MerchantOrderID)

	// Start order status polling
	go zotaApi.PollOrderStatus(zotaOrderStatusRequest)

	// Redirect user to deposit page
	c.Redirect(http.StatusFound, response.Data.DepositUrl)
}
