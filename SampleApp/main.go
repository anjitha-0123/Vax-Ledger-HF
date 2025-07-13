package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TempLog struct {
	AssetType string `json:"assettype"`
	BatchID   string `json:"batchID"`
	TempNow   int    `json:"tempnow"`
	TimeStamp string `json:"timestamp"`
}

func main() {
	router := gin.Default()
	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Add Temperature Log (Org3MSP only)
	router.POST("/api/temp", func(ctx *gin.Context) {
		var req struct {
			BatchID string `json:"batchID"`
			TempNow int    `json:"tempnow"`
		}
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
			return
		}
		fmt.Printf("Temperature log input: %+v\n", req)

		result := submitTxnFn(
			"org3",
			"coldchannel",
			"Vax-Ledger",
			"TempContract",
			"invoke",
			make(map[string][]byte),
			"AddTemperatureLog",
			req.BatchID,
			strconv.Itoa(req.TempNow),
		)

		fmt.Println("Transaction Result:", result)
		ctx.JSON(http.StatusOK, gin.H{"message": "Temperature log added", "batchID": req.BatchID})
	})

	// Get Temperature Log History (Org1MSP, Org2MSP, Org3MSP)
	router.GET("/api/temp/:batchID", func(ctx *gin.Context) {
		batchID := ctx.Param("batchID")
		result := submitTxnFn(
			"org2",
			"coldchannel",
			"Vax-Ledger",
			"TempContract",
			"query",
			make(map[string][]byte),
			"GetTemperatureLogHistory",
			batchID,
		)

		var logs []string
		json.Unmarshal([]byte(result), &logs)

		ctx.JSON(http.StatusOK, gin.H{
			"batchID": batchID,
			"logs":    logs,
		})
	})

	// Verify Temperature Logs (Org1 or Org2)
	router.GET("/api/verify/:batchID", func(ctx *gin.Context) {
		batchID := ctx.Param("batchID")
		result := submitTxnFn(
			"org1",
			"coldchannel",
			"Vax-Ledger",
			"TempContract",
			"query",
			make(map[string][]byte),
			"VerifyTemperatureLogs",
			batchID,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"batchID": batchID,
			"result":  result,
		})
	})

	// Start Delivery (Org3MSP)
	router.POST("/api/delivery/start/:batchID", func(ctx *gin.Context) {
		batchID := ctx.Param("batchID")
		result := submitTxnFn(
			"org3",
			"coldchannel",
			"Vax-Ledger",
			"TempContract",
			"invoke",
			make(map[string][]byte),
			"StartDelivery",
			batchID,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Delivery started",
			"result":  result,
		})
	})

	// Complete Delivery (Org3MSP)
	router.POST("/api/delivery/complete/:batchID", func(ctx *gin.Context) {
		batchID := ctx.Param("batchID")
		result := submitTxnFn(
			"org3",
			"coldchannel",
			"Vax-Ledger",
			"TempContract",
			"invoke",
			make(map[string][]byte),
			"CompleteDelivery",
			batchID,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Delivery completed",
			"result":  result,
		})
	})

	// Get Delivery Status (Org1MSP or Org2MSP)
	router.GET("/api/delivery/status/:batchID", func(ctx *gin.Context) {
		batchID := ctx.Param("batchID")
		result := submitTxnFn(
			"org1",
			"coldchannel",
			"Vax-Ledger",
			"TempContract",
			"query",
			make(map[string][]byte),
			"GetDeliveryStatus",
			batchID,
		)

		ctx.JSON(http.StatusOK, gin.H{
			"batchID": batchID,
			"status":  result,
		})
	})

	router.Run(":3002")
}
