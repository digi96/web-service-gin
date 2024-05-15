package controllers

import (
	"context"
	"database/sql"
	db "example/web-service-gin/db/sqlc"
	"example/web-service-gin/schemas"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RideOrderController struct {
	db  *db.Queries
	ctx context.Context
}

func NewRideOrderController(db *db.Queries, ctx context.Context) *RideOrderController {
	return &RideOrderController{db, ctx}
}

func (roc *RideOrderController) CreatRideOrder(ctx *gin.Context) {
	var payload *schemas.CreateRideOrder

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	args := &db.CreateOrderParams{
		ContactID:   uuid.MustParse(payload.ContactId),
		RiderName:   payload.RiderName,
		RiderPhone:  payload.RiderPhone,
		Destination: payload.Destinaion,
	}

	rideOrder, err := roc.db.CreateOrder(ctx, *args)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving order", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully created order", "riderorder": rideOrder})

}

// Get a single handler
func (roc *RideOrderController) GetOrderById(ctx *gin.Context) {
	rideOrderId, _ := strconv.ParseInt(ctx.Param("orderId"), 10, 64)

	rideOrder, err := roc.db.GetOrderById(ctx, int32(rideOrderId))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve order with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving order", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully retrived id", "order": rideOrder})
}

func (roc *RideOrderController) UpdateRideOrder(ctx *gin.Context) {
	var payload *schemas.UpdateRideOrder
	rideOrderId, _ := strconv.ParseInt(ctx.Param("orderId"), 10, 64)

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "Failed payload", "error": err.Error()})
		return
	}

	now := time.Now()
	fmt.Println("PickUpAt:", payload.PickUpAt)
	pickUpTime, err := time.Parse("2006-01-02 15:04:05", payload.PickUpAt)
	if err != nil {
		fmt.Println("format time failed,", err.Error())
	}

	args := &db.UpdateOrderParams{
		RiderorderID: int32(rideOrderId),
		PickupAt:     sql.NullTime{Time: pickUpTime, Valid: payload.PickUpAt != ""},
		UpdatedAt:    sql.NullTime{Time: now, Valid: true},
	}

	rideOrder, err := roc.db.UpdateOrder(ctx, *args)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": "Failed to retrieve rideoder with this ID"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "Failed retrieving rideorder", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "successfully updated rideorder", "order": rideOrder})

}
