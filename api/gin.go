package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/federlizer/alokin-zota-integration/internal"
	"github.com/federlizer/alokin-zota-integration/internal/storage"
	"github.com/federlizer/alokin-zota-integration/zota"
)

func SetupApi(zotaApi zota.IZotaAPI, orderRepo storage.OrderRepo) *gin.Engine {
	engine := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	engine.Use(cors.New(corsConfig))

	engine.Use(func(c *gin.Context) {
		c.Set("zotaApi", zotaApi)
		c.Set("orderRepo", orderRepo)
	})

	engine.GET("/ping", pingHandler)

	engine.GET("/order", getOrdersHandler)
	engine.POST("/order", orderHandler)

	return engine
}

func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func getOrdersHandler(c *gin.Context) {
	type response struct {
		Orders []*internal.Order `json:"orders"`
	}

	orderRepo := c.MustGet("orderRepo").(storage.OrderRepo)
	orders := orderRepo.GetAll()
	resp := response{Orders: orders}

	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	c.Data(http.StatusOK, "application/json", data)
}

type OrderHandlerParams struct {
	Description string  `json:"description" form:"description" binding:"required,max=128"`
	Amount      float64 `json:"amount" form:"amount" binding:"required"`
}

func orderHandler(c *gin.Context) {
	// TODO - is there a better way to do this?
	zotaApi := c.MustGet("zotaApi").(zota.IZotaAPI)
	orderRepo := c.MustGet("orderRepo").(storage.OrderRepo)

	var params OrderHandlerParams
	err := c.Bind(&params)
	if err != nil {
		log.Printf("Couldn't parse form parameters: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
		})
		return
	}

	// Ensure amount is a positive, non-zero number
	if params.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "amount needs to be a positive number",
		})
		return
	}

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

	order := internal.NewOrder(&user, params.Amount, params.Description)

	addedToRepo := false
	retries := 0
	for !addedToRepo && retries < 10 {
		// This call can only fail if there is a duplicate ID for an order
		err := orderRepo.AddOrder(order)
		retries += 1

		if err != nil {
			log.Printf("Couldn't add new order to order repo. Likely a key error: %v\n", err)
			continue
		}

		addedToRepo = true
	}

	// But let's make sure that we've added it anyways
	if !addedToRepo {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to add new order to repo",
		})
		return
	}

	zotaDepositRequest := zota.FromOrder(order, zotaApi.EndpointId(), zotaApi.SecretKey())

	// Make request to Zota API
	response, err := zotaApi.Deposit(zotaDepositRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	zotaOrderStatusRequest := zota.NewZotaOrderStatusRequest(response.Data.OrderId, response.Data.MerchantOrderID)

	// Start Order Status polling
	go zotaApi.PollOrderStatus(zotaOrderStatusRequest, order)

	// Redirect user to deposit page
	c.Redirect(http.StatusFound, response.Data.DepositUrl)
}
