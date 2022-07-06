/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Transaction struct {
	IDBill   string `json:"id_bill"`
	IDProduct   string `json:"id_product"`
	Quantity  string `json:"quantity"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Transaction
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	transactions := []Transaction{
		Transaction{IDBill: "62c27c668772d1ba23884fa1", IDProduct: "62a60b1f717c6989b0165e55", Quantity: "2"},
		Transaction{IDBill: "62c27c668772d1ba23884fa1", IDProduct: "62a60b42717c6989b0165e58", Quantity: "2"},
		Transaction{IDBill: "62c2a2cf9b6915dea0039b5b", IDProduct: "62a61170717c6989b0165e73", Quantity: "3"},
		Transaction{IDBill: "62c2a3ae9b6915dea0039b5d", IDProduct: "62a611e4717c6989b0165e7a", Quantity: "3"},
	}

	for i, transaction := range transactions {
		transactionAsBytes, _ := json.Marshal(transaction)
		err := ctx.GetStub().PutState("TRANSACTION"+strconv.Itoa(i), transactionAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, transactionNumber string, idBill string, idProduct string, quantity string) error {
	transaction := Transaction{
		IDBill:   idBill,
		IDProduct: idProduct,
		Quantity:  quantity,
	}

	transactionAsBytes, _ := json.Marshal(transaction)

	return ctx.GetStub().PutState(transactionNumber, transactionAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryTransaction(ctx contractapi.TransactionContextInterface, transactionNumber string) (*Transaction, error) {
	transactionAsBytes, err := ctx.GetStub().GetState(transactionNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if transactionAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", transactionNumber)
	}

	transaction := new(Transaction)
	_ = json.Unmarshal(transactionAsBytes, transaction)

	return transaction, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllTransactions(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "TRANSACTION0"
	endKey := "TRANSACTION999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		transaction := new(Transaction)
		_ = json.Unmarshal(queryResponse.Value, transaction)

		queryResult := QueryResult{Key: queryResponse.Key, Record: transaction}
		results = append(results, queryResult)
	}

	return results, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
