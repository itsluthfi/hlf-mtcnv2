package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type GoChaincode struct {
	contractapi.Contract
}

type Asset struct {
	DocType    string `json:"docType"`    //docType is used to distinguish the various types of objects in state database
	UsernameID string `json:"usernameId"` //the field tags are needed to keep case from bouncing around
	Name       string `json:"name"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	CoinAmount int    `json:"coinAmount"`
}

// type TxRecord struct {
// 	Sender     string `json:"sender"`
// 	Receiver   string `json:"receiver"`
// 	CoinAmount int    `json:"coinAmount"`
// }

type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// CreateAsset initializes a new asset in the ledger
func (t *GoChaincode) CreateAsset(ctx contractapi.TransactionContextInterface, usernameID string, name string, password string, phone string, email string, coinAmount int) error {
	exists, err := t.AssetExists(ctx, usernameID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if exists {
		return fmt.Errorf("asset already exists: %s", usernameID)
	}

	asset := &Asset{
		DocType:    "asset",
		UsernameID: usernameID,
		Name:       name,
		Password:   password,
		Phone:      phone,
		Email:      email,
		CoinAmount: coinAmount,
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(usernameID, assetBytes)
	if err != nil {
		return err
	}

	return nil

	// //  Create an index to enable color-based range queries, e.g. return all blue assets.
	// //  An 'index' is a normal key-value entry in the ledger.
	// //  The key is a composite key, with the elements that you want to range query on listed first.
	// //  In our case, the composite key is based on indexName~color~name.
	// //  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	// colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Color, asset.ID})
	// if err != nil {
	// 	return err
	// }
	// //  Save index entry to world state. Only the key name is needed, no need to store a duplicate copy of the asset.
	// //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	// value := []byte{0x00}
	// return ctx.GetStub().PutState(colorNameIndexKey, value)
}

// ReadAsset retrieves an asset from the ledger
func (t *GoChaincode) ReadAsset(ctx contractapi.TransactionContextInterface, usernameID string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(usernameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", usernameID, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", usernameID)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the ledger
func (t *GoChaincode) UpdateAsset(ctx contractapi.TransactionContextInterface, usernameID string, name string, password string, phone string, email string, coinAmount int) error {
	exists, err := t.AssetExists(ctx, usernameID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if !exists {
		return fmt.Errorf("asset is not exists: %s", usernameID)
	}

	asset := &Asset{
		DocType:    "asset",
		UsernameID: usernameID,
		Name:       name,
		Password:   password,
		Phone:      phone,
		Email:      email,
		CoinAmount: coinAmount,
	}
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(usernameID, assetBytes)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAsset removes an asset key-value pair from the ledger
func (t *GoChaincode) DeleteAsset(ctx contractapi.TransactionContextInterface, usernameID string) error {
	_, err := t.ReadAsset(ctx, usernameID)
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelState(usernameID)
	if err != nil {
		return fmt.Errorf("failed to delete asset %s: %v", usernameID, err)
	}

	return nil

	// colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.Color, asset.ID})
	// if err != nil {
	// 	return err
	// }

	// // Delete index entry
	// return ctx.GetStub().DelState(colorNameIndexKey)
}

// TransferAsset transfers an asset by setting an new owner name on the asset
func (t *GoChaincode) TransferAsset(ctx contractapi.TransactionContextInterface, senderID string, receiverID string, amount int) error {
	sender, err := t.ReadAsset(ctx, senderID)
	if err != nil {
		return err
	}

	if sender.CoinAmount <= 0 || sender.CoinAmount < amount {
		return fmt.Errorf("not enough coins to transfer")
	}

	newSender := Asset{
		DocType:    sender.DocType,
		UsernameID: sender.UsernameID,
		Name:       sender.Name,
		Password:   sender.Password,
		Phone:      sender.Phone,
		Email:      sender.Email,
		CoinAmount: sender.CoinAmount - amount,
	}

	senderBytes, err := json.Marshal(newSender)
	if err != nil {
		return err
	}

	receiver, err := t.ReadAsset(ctx, receiverID)
	if err != nil {
		return err
	}

	newReceiver := Asset{
		DocType:    receiver.DocType,
		UsernameID: receiver.UsernameID,
		Name:       receiver.Name,
		Password:   receiver.Password,
		Phone:      receiver.Phone,
		Email:      receiver.Email,
		CoinAmount: receiver.CoinAmount + amount,
	}

	receiverBytes, err := json.Marshal(newReceiver)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(senderID, senderBytes)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(receiverID, receiverBytes)
}

// GetAssetHistory returns the chain of custody for an asset since issuance.
func (t *GoChaincode) GetAssetHistory(ctx contractapi.TransactionContextInterface, usernameID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", usernameID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(usernameID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				UsernameID: usernameID,
			}
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: response.Timestamp.AsTime(),
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

// AssetExists returns true when asset with given ID exists in the ledger.
func (t *GoChaincode) AssetExists(ctx contractapi.TransactionContextInterface, usernameID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(usernameID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", usernameID, err)
	}

	return assetBytes != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&GoChaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
