package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TempContract struct {
	contractapi.Contract
}

type Temperature struct {
	AssetType string `json:"assettype"`
	BatchID   string `json:"batchID"`
	TempNow   int    `json:"tempnow"`
	TimeStamp string `json:"timestamp"`
}

type DeliveryStatus struct {
	BatchID string `json:"batchID"`
	Status  string `json:"status"` // In-Transit or Delivered
}

// AddTemperatureLog with automatic timestamp
func (t *TempContract) AddTemperatureLog(ctx contractapi.TransactionContextInterface, batchID string, tempNow int) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "Org3MSP" {
		return fmt.Errorf("access denied: AddTemperatureLog can only be invoked by Org3MSP")
	}

	timestamp := time.Now().Format(time.RFC3339)

	log := Temperature{
		AssetType: "TemperatureLog",
		BatchID:   batchID,
		TempNow:   tempNow,
		TimeStamp: timestamp,
	}

	logBytes, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal temperature log: %v", err)
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey("TempLog", []string{batchID, timestamp})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	return ctx.GetStub().PutState(compositeKey, logBytes)
}

// GetTemperatureLogHistory with °C unit
func (t *TempContract) GetTemperatureLogHistory(ctx contractapi.TransactionContextInterface, batchID string) ([]string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "Org2MSP" && clientMSPID != "Org3MSP" && clientMSPID != "Org1MSP" {
		return nil, fmt.Errorf("access denied: GetTemperatureLogHistory is only available to Org1MSP, Org2MSP, and Org3MSP")
	}

	var results []string

	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("TempLog", []string{batchID})
	if err != nil {
		return nil, fmt.Errorf("failed to get composite keys: %v", err)
	}
	defer iterator.Close()

	for iterator.HasNext() {
		kvResult, err := iterator.Next()
		if err != nil {
			return nil, err
		}

		var tempLog Temperature
		err = json.Unmarshal(kvResult.Value, &tempLog)
		if err != nil {
			continue
		}

		display := fmt.Sprintf("Temp: %d°C at %s", tempLog.TempNow, tempLog.TimeStamp)
		results = append(results, display)
	}

	return results, nil
}

// VerifyTemperatureLogs with °C unit, Org1 and Org2 can verify
func (t *TempContract) VerifyTemperatureLogs(ctx contractapi.TransactionContextInterface, batchID string) (string, error) {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get client MSP ID: %v", err)
	}
	if clientMSPID != "Org2MSP" && clientMSPID != "Org1MSP" {
		return "", fmt.Errorf("access denied: only Org1MSP and Org2MSP can verify temperature logs")
	}

	vaccine, err := ReadPrivateState(ctx, batchID)
	if err != nil {
		return "", fmt.Errorf("failed to read vaccine batch: %v", err)
	}

	logs, err := t.GetTemperatureLogHistory(ctx, batchID)
	if err != nil {
		return "", fmt.Errorf("failed to get temperature log history: %v", err)
	}

	for _, entry := range logs {
		var tempValue int
		_, err := fmt.Sscanf(entry, "Temp: %d°C at", &tempValue)
		if err != nil {
			continue
		}

		if tempValue < vaccine.MinTemp || tempValue > vaccine.MaxTemp {
			return fmt.Sprintf("Not Verified: Temp %d°C out of range [%d°C, %d°C]", tempValue, vaccine.MinTemp, vaccine.MaxTemp), nil
		}
	}

	return fmt.Sprintf("Verified: All temperatures within safe range [%d°C, %d°C]", vaccine.MinTemp, vaccine.MaxTemp), nil
}

// StartDelivery sets status to In-Transit
func (t *TempContract) StartDelivery(ctx contractapi.TransactionContextInterface, batchID string) error {
	status := DeliveryStatus{
		BatchID: batchID,
		Status:  "In-Transit",
	}

	statusBytes, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery status: %v", err)
	}

	return ctx.GetStub().PutState("DELIVERY_"+batchID, statusBytes)
}

// CompleteDelivery sets status to Delivered
func (t *TempContract) CompleteDelivery(ctx contractapi.TransactionContextInterface, batchID string) error {
	status := DeliveryStatus{
		BatchID: batchID,
		Status:  "Delivered",
	}

	statusBytes, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery status: %v", err)
	}

	return ctx.GetStub().PutState("DELIVERY_"+batchID, statusBytes)
}

// GetDeliveryStatus returns delivery status
func (t *TempContract) GetDeliveryStatus(ctx contractapi.TransactionContextInterface, batchID string) (string, error) {
	statusBytes, err := ctx.GetStub().GetState("DELIVERY_" + batchID)
	if err != nil {
		return "", fmt.Errorf("failed to get delivery status: %v", err)
	}
	if statusBytes == nil {
		return "No delivery status found", nil
	}

	var status DeliveryStatus
	err = json.Unmarshal(statusBytes, &status)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal delivery status: %v", err)
	}

	return fmt.Sprintf("Batch %s is currently: %s", batchID, status.Status), nil
}
