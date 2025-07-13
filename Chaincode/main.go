package main

import (
	"coldvax/contracts"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	vaxContract := new(contracts.VaxContract)
	tempContract := new(contracts.TempContract)

	chaincode, err := contractapi.NewChaincode(vaxContract,tempContract)

	if err != nil {
		log.Panicf("Could not create chaincode : %v", err)
	}
	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode : %v", err)
	}

}
