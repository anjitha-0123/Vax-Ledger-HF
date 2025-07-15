# 🧊 Vax Ledger - Vaccine Cold Chain Management System
A blockchain-powered system for secure tracking, monitoring, and verification of vaccine batches across the supply chain. It ensures temperature compliance via IoT or manual inputs and safeguards environmental integrity using Hyperledger Fabric and RESTful APIs.

## Problem Statement
Ensuring the efficacy and safety of vaccines during storage and transportation is a critical challenge in the global healthcare supply chain.Vaccines are highly sensitive to temperature fluctuations.
traditional cold chain systems suffer from: Lack of transparency and traceability in shipping records.

## 🚀 Key Features
✅ Vaccine Batch Registration

Manufacturers (Org1) register batches with private data.

📦 Temperature Tracking

Manual or IoT-based updates by Transporter (Org3).

🔍 Batch Verification

Retailers (Org2) query and verify batch temperature & details.

🔐 Tamper-Proof Ledger

Secure tracking using Hyperledger Fabric.

🌐 RESTful APIs with Gin

Simple API access for apps and dashboards.

## Tech Stack
![Hyperledger Fabric](https://img.shields.io/badge/Hyperledger%20Fabric-2C3E50?style=for-the-badge&logo=hyperledger&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/-Docker-2496ED?style=flat-square&logo=Docker&logoColor=white)
![HTML](https://img.shields.io/badge/-HTML-E34F26?style=flat-square&logo=HTML5&logoColor=white)

## To Clone this Repository
```
git@github.com:anjitha-0123/Vax-Ledger-HF.git
```

## To Build the Network
```
cd fabric-samples/test-network
```
```
./network.sh up createChannel -c coldchannel -ca -s couchdb
```
### Adding Org3
```
cd addOrg3
```
```
./addOrg3.sh up -c coldchannel -ca -s couchdb
```
```
cd ..
```
## To deploy the ChainCode
```
./network.sh deployCC -ccn Vax-Ledger -ccp ../../Vax-Ledger-HF/Chaincode/ -ccl go -c coldchannel -cccg ../../Vax-Ledger-HF/Chaincode/collections.json
```
```
./network.sh deployCC -ccn Vax-Ledger -ccp ../../Vax-Ledger-HF/Chaincode/ -ccl go -c coldchannel -ccv 2.0 -ccs 2 -cccg ../../Vax-Ledger-HF/Chaincode/collections.json

```


### Invoke
###  general variables
```
export FABRIC_CFG_PATH=$PWD/../config/

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

export ORG1_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_TLS_ENABLED=true
```
### org1 env variables
```
export CORE_PEER_LOCALMSPID=Org1MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

export CORE_PEER_ADDRESS=localhost:7051
```
### Pvt Data
### Setting the transient data
```
export MANUFACTURER=$(echo -n "Pfizer" | base64)
export VACCINE_TYPE=$(echo -n "mRNA" | base64)
```
```
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile $ORDERER_CA \
  -C coldchannel \
  -n Vax-Ledger \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
  -c '{"Args":["VaxContract:CreateBatch","Batch-01","12-01-2025","12-01-2027","3","5"]}' \      
  --transient "{\"manufacturer\":\"$MANUFACTURER\",\"vaccineType\":\"$VACCINE_TYPE\"}"

```
### Query
```
peer chaincode query -C coldchannel -n Vax-Ledger -c '{"Args":["VaxContract:ReadBatch","Batch-02"]}'
```
```
peer chaincode query -C coldchannel -n Vax-Ledger -c '{"function":"GetAllBatch","Args":[]}'
```
### Query Temp History
```
peer chaincode query   -C coldchannel   -n Vax-Ledger   -c '{"Args":["TempContract:GetTemperatureLogHistory", "Batch-01"]}'
```
### Delete Invoke
```
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile $ORDERER_CA \
  -C coldchannel \
  -n Vax-Ledger \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT \
  -c '{"Args":["VaxContract:DeleteBatch","Batch-01"]}'
```

### Env veriables of Org2
```
export CORE_PEER_LOCALMSPID=Org2MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp

export CORE_PEER_ADDRESS=localhost:9051


```

### Query Temp History
```
peer chaincode query   -C coldchannel   -n Vax-Ledger   -c '{"Args":["TempContract:GetTemperatureLogHistory", "Batch-01"]}'
```
### Verify Temp
```
export ORDERER_CA=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export ORG2_PEER_TLSROOTCERT=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com \
  --tls \
  --cafile $ORDERER_CA \
  -C coldchannel \
  -n Vax-Ledger \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT \
  -c '{"Args":["TempContract:VerifyTemperatureLogs","Batch-01"]}'


```

# UI
```
cd SampleApp
```
## To run the Client App
```
go run .
```

## ALternative Running Cammands for Org 3
### org3 env var
```
export CORE_PEER_LOCALMSPID=Org3MSP

export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org3.example.com/peers/peer0.org3.example.com/tls/ca.crt

export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp

export CORE_PEER_ADDRESS=localhost:11051
```
### invoke for templog org3
```
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile $ORDERER_CA -C coldchannel -n Vax-Ledger --peerAddresses localhost:7051 --tlsRootCertFiles $ORG1_PEER_TLSROOTCERT --peerAddresses localhost:9051 --tlsRootCertFiles $ORG2_PEER_TLSROOTCERT -c '{"Args":["TempContract:AddTemperatureLog","Batch-01","5.3","07-07-2025T15:35:00Z"]}'


```

## To Down the Network
```
./network.sh down
```



