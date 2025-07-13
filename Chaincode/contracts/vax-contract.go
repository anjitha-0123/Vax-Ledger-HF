package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type VaxContract struct {
	contractapi.Contract
}
type VaccineBatch struct {
	DocType           string `json:"docType"`
	BatchID           string `json:"batchID"`
	Manufacturer      string `json:"manufacturer"`
	VaccineType       string `json:"vaccineType"`
	ManufactureDate   string `json:"manufactureDate"`
	ExpiryDate        string `json:"expiryDate"`
	MinTemp           int    `json:"minTemp"`
	MaxTemp           int    `json:"maxTemp"`
	Status            string `json:"status"`
	CreationTimestamp string `json:"creationTimestamp"`
}

func getCollectionName() string {
	collectionName := "BatchCollection"
	return collectionName
}

// Org1
func (v *VaxContract) VaccineExists(ctx contractapi.TransactionContextInterface, batchID string) (bool, error) {
	collectionName := getCollectionName()
	data, err := ctx.GetStub().GetPrivateData(collectionName, batchID)
	if err != nil {
		return false, fmt.Errorf("failed to read from private data collection: %v", err)
	}
	return data != nil, nil
}

// Org1
func (v *VaxContract) CreateBatch(ctx contractapi.TransactionContextInterface, batchID string, manufactureDate string, expiryDate string, minTemp int, maxTemp int) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID != "Org1MSP" {
		return "", fmt.Errorf("VaccineBatch cannot be created by organisation with MSPID %v", clientOrgID)
	}

	exists, err := v.VaccineExists(ctx, batchID)
	if err != nil {
		return "", fmt.Errorf("Could not read from world state: %s", err)
	}
	if exists {
		return "", fmt.Errorf("The asset %s already exists", batchID)
	}

	transientData, err := ctx.GetStub().GetTransient()
	if err != nil {
		return "", err
	}
	if len(transientData) == 0 {
		return "", fmt.Errorf("Please provide the private data of manufacturer and vaccineType")
	}

	manufacturer, exists := transientData["manufacturer"]
	if !exists {
		return "", fmt.Errorf("The manufacturer was not specified in transient data. Please try again")
	}

	vaccineType, exists := transientData["vaccineType"]
	if !exists {
		return "", fmt.Errorf("The vaccineType was not specified in transient data. Please try again")
	}
	timestamp, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to get transaction timestamp: %v", err)
	}
	creationTime := time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).UTC().Format(time.RFC3339)

	vaccine := &VaccineBatch{
		DocType:           "vaccineBatch",
		BatchID:           batchID,
		Manufacturer:      string(manufacturer),
		VaccineType:       string(vaccineType),
		ManufactureDate:   manufactureDate,
		ExpiryDate:        expiryDate,
		MinTemp:           minTemp,
		MaxTemp:           maxTemp,
		Status:            "Created",
		CreationTimestamp: creationTime,
	}

	vaccineJSON, err := json.Marshal(vaccine)
	if err != nil {
		return "", err
	}

	collectionName := getCollectionName()

	err = ctx.GetStub().PutPrivateData(collectionName, batchID, vaccineJSON)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("VaccineBatch with ID %v added successfully", batchID), nil
}

func (v *VaxContract) ReadBatch(ctx contractapi.TransactionContextInterface, batchID string) (string, error) {
	vaccine, err := ReadPrivateState(ctx, batchID)
	if err != nil {
		return "", err
	}

	display := fmt.Sprintf(
		"Batch ID: %s\nManufacturer: %s\nVaccine Type: %s\nManufacture Date: %s\nExpiry Date: %s\nTemperature Range: %d°C to %d°C\nStatus: %s",
		vaccine.BatchID,
		vaccine.Manufacturer,
		vaccine.VaccineType,
		vaccine.ManufactureDate,
		vaccine.ExpiryDate,
		vaccine.MinTemp,
		vaccine.MaxTemp,
		vaccine.Status,
	)

	return display, nil
}

func ReadPrivateState(ctx contractapi.TransactionContextInterface, batchID string) (*VaccineBatch, error) {
	collectionName := getCollectionName()

	bytes, err := ctx.GetStub().GetPrivateData(collectionName, batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to read private data: %v", err)
	}
	if bytes == nil {
		return nil, fmt.Errorf("no private data found for batch ID %s", batchID)
	}

	vaccine := new(VaccineBatch)

	err = json.Unmarshal(bytes, vaccine)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal private data collection data to type VaccineBatch")
	}

	return vaccine, nil
}

func (v *VaxContract) DeleteBatch(ctx contractapi.TransactionContextInterface, batchID string) error {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}
	if clientOrgID == "Org1MSP" {
		exists, err := v.VaccineExists(ctx, batchID)
		if err != nil {
			return fmt.Errorf("Could not read from world state. %s", err)
		} else if !exists {
			return fmt.Errorf("The asset %s does not exist", batchID)
		}
		collectionName := getCollectionName()
		return ctx.GetStub().DelPrivateData(collectionName, batchID)
	} else {
		return fmt.Errorf("Organisation with %v cannot delete the order", clientOrgID)
	}
}

func (v *VaxContract) GetAllBatch(ctx contractapi.TransactionContextInterface) ([]*VaccineBatch, error) {
	collectionName := getCollectionName()
	queryString := `{"selector": {"docType": "vaccineBatch"}}`
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult(collectionName, queryString)
	if err != nil {

		return nil, err
	}
	defer resultsIterator.Close()
	return orderResultIteratorFunction(resultsIterator)
}

// iterator function
func orderResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*VaccineBatch, error) {

	var vaccines []*VaccineBatch

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var vaccine VaccineBatch
		err = json.Unmarshal(queryResult.Value, &vaccine)
		if err != nil {
			return nil, err
		}
		vaccines = append(vaccines, &vaccine)
	}
	return vaccines, nil
}
func (v *VaxContract) DeliverToTransporter(ctx contractapi.TransactionContextInterface, batchID string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID != "Org1MSP" {
		return "", fmt.Errorf("Only Org1 can deliver the vaccine batch")
	}

	vaccine, err := ReadPrivateState(ctx, batchID)
	if err != nil {
		return "", err
	}

	if vaccine.Status != "Created" {
		return "", fmt.Errorf("Cannot deliver batch %s as it is in %s state", batchID, vaccine.Status)
	}

	vaccine.Status = "In-transit"

	vaccineJSON, err := json.Marshal(vaccine)
	if err != nil {
		return "", err
	}

	collectionName := getCollectionName()
	err = ctx.GetStub().PutPrivateData(collectionName, batchID, vaccineJSON)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Vaccine batch %s has been delivered to transporter", batchID), nil
}
