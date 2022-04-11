/*
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"time"
)

// CBPMChaincode implements the fabric-contract-api-go programming model
type CBPMChaincode struct {
	contractapi.Contract
}

type DeliveryArrangement struct {
	ObjectType string  `json:"objectType"`
	TradeID    string  `json:"tradeID"`
	AssetName  string  `json:"assetName"`
	Quantity   int     `json:"quantity"`
	StartPlace string  `json:"startPlace"`
	EndPlace   string  `json:"endPlace"`
	Fee        float32 `json:"fee"`
	CreateTime string  `json:"createTime"`
	UpdateTime string  `json:"updateTime"`
	OwnerOrg   string  `json:"ownerOrg"`
	HandlerOrg string  `json:"handler"`
	Status     int     `json:"status"` // 0: 未处理 1：已接单 2: 已完成 3: 确认已完成
	Note       string  `json:"note"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *DeliveryArrangement `json:"record"`
	TxId      string               `json:"txId"`
	Timestamp time.Time            `json:"timestamp"`
	IsDelete  bool                 `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*DeliveryArrangement `json:"records"`
	FetchedRecordsCount int32                  `json:"fetchedRecordsCount"`
	Bookmark            string                 `json:"bookmark"`
}

// CreateDeliveryArrangement 创建DeliveryArrangement，需要transient传入tradeID，assetName，quantity，startPlace，endPlace，fee，note，返回创建好的DeliveryArrangement
func (t *CBPMChaincode) CreateDeliveryArrangement(ctx contractapi.TransactionContextInterface) (*DeliveryArrangement, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientArrangementJSON, ok := transMap["arrangement"]
	if !ok {
		return nil, fmt.Errorf("arrangement not found in the transient map")
	}
	type arrangementTransientInput struct {
		TradeID    string  `json:"tradeID"`
		AssetName  string  `json:"assetName"`
		Quantity   int     `json:"quantity"`
		StartPlace string  `json:"startPlace"`
		EndPlace   string  `json:"endPlace"`
		Fee        float32 `json:"fee"`
		Note       string  `json:"note"`
	}
	var arrangementInput arrangementTransientInput
	err = json.Unmarshal(transientArrangementJSON, &arrangementInput)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	if len(arrangementInput.TradeID) == 0 {
		return nil, fmt.Errorf("trade ID must be a non-empty string")
	}
	if len(arrangementInput.AssetName) == 0 {
		return nil, fmt.Errorf("asset name must be a non-empty string")
	}
	if len(arrangementInput.StartPlace) == 0 {
		return nil, fmt.Errorf("start place must be a non-empty string")
	}
	if len(arrangementInput.EndPlace) == 0 {
		return nil, fmt.Errorf("end place must be a non-empty string")
	}
	if arrangementInput.Quantity <= 0 {
		return nil, fmt.Errorf("asset quantity must be a positive number")
	}
	if arrangementInput.Fee <= 0 {
		return nil, fmt.Errorf("fee must be a positive number")
	}

	exists, err := t.deliveryArrangementExists(ctx, arrangementInput.TradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery arrangement: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("fail to create delivery arrangement: arrangement for trade %s already exists", arrangementInput.TradeID)
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery arrangement: %v", err)
	}

	deliveryArrangement := &DeliveryArrangement{
		ObjectType: "DeliveryArrangement",
		TradeID:    arrangementInput.TradeID,
		AssetName:  arrangementInput.AssetName,
		Quantity:   arrangementInput.Quantity,
		StartPlace: arrangementInput.StartPlace,
		EndPlace:   arrangementInput.EndPlace,
		Fee:        arrangementInput.Fee,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     0,
		OwnerOrg:   clientOrgID,
		HandlerOrg: "",
		Note: 		arrangementInput.Note,
	}

	deliveryArrangementBytes, err := json.Marshal(deliveryArrangement)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery arrangement: %v", err)
	}
	err = ctx.GetStub().PutState(arrangementInput.TradeID, deliveryArrangementBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery arrangement: %v", err)
	}
	return deliveryArrangement, nil
}

// DeleteDeliveryArrangement 删除DeliveryArrangement，args传入tradeID
func (t *CBPMChaincode) DeleteDeliveryArrangement(ctx contractapi.TransactionContextInterface, tradeID string) error {
	// TODO 只能拥有者删除，并且不能在进行过程中删除
	deliveryArrangement, err := t.GetDeliveryArrangement(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to delete delivery arrangement: %v", err)
	}
	if deliveryArrangement.Status != 0 {
		return fmt.Errorf("fail to delete delivery arrangement: cannot delete arrangement in progress")
	}
	return ctx.GetStub().DelState(tradeID)
}

// HandleDeliveryArrangement 非owner的一方处理DeliveryArrangement，agrs传入tradeID
func (t *CBPMChaincode) HandleDeliveryArrangement(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryArrangement, err := t.GetDeliveryArrangement(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to handle delivery arrangement: %v", err)
	}
	if deliveryArrangement.HandlerOrg != "" {
		return fmt.Errorf("fail to handle delivery arrangement: delivery arrangement for trade #{tradeID} has been handled")
	}
	if deliveryArrangement.Status != 0 {
		return fmt.Errorf("fail to handle delivery arrangement: delivery arrangement(status: %d) for trade #{tradeID} has been handled", deliveryArrangement.Status)
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to handle delivery arrangement: %v", err)
	}
	if deliveryArrangement.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to handle delivery arrangement: cannot handle as owner")
	}
	deliveryArrangement.HandlerOrg = clientOrgID
	deliveryArrangement.Status = 1
	deliveryArrangement.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	deliveryArrangementBytes, err := json.Marshal(deliveryArrangement)
	if err != nil {
		return fmt.Errorf("fail to handle delivery arrangement: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, deliveryArrangementBytes)
}

// FinishDeliveryArrangement handler完成DeliveryArrangement，args传入tradeID
func (t *CBPMChaincode) FinishDeliveryArrangement(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryArrangement, err := t.GetDeliveryArrangement(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to finish delivery arrangement: %v", err)
	}
	if deliveryArrangement.Status == 0 {
		return fmt.Errorf("fail to finish delivery arrangement: arrangement for trade #{tradeID} has not been handled")
	}
	if deliveryArrangement.Status == 2 {
		return fmt.Errorf("fail to finish delivery arrangement: arrangement for trade #{tradeID} has been finished")
	}
	if deliveryArrangement.HandlerOrg == "" {
		return fmt.Errorf("fail to finish delivery arrangement: no handler is specified")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to finish delivery arrangement: %v", err)
	}
	if deliveryArrangement.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to finish delivery arrangement: cannot finish as owner")
	}
	if deliveryArrangement.HandlerOrg != clientOrgID {
		return fmt.Errorf("fail to finish delivery arrangement: unauthorized handler #{clientOrgID}}")
	}
	deliveryArrangement.Status = 2
	deliveryArrangement.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	deliveryArrangementBytes, err := json.Marshal(deliveryArrangement)
	if err != nil {
		return fmt.Errorf("fail to finish delivery arrangement: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, deliveryArrangementBytes)
}

// ConfirmFinishDeliveryArrangement owner确定handler完成deliveryArrangement，args传入tradeID，在例子中需要manufacturer收到货物过后，middleman才能确认carrier完成
func (t *CBPMChaincode) ConfirmFinishDeliveryArrangement(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryArrangement, err := t.GetDeliveryArrangement(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to confirm finish delivery arrangement: %v", err)
	}
	if deliveryArrangement.Status == 0 {
		return fmt.Errorf("fail to confirm finish delivery arrangement: arrangement has not been handled")
	}
	if deliveryArrangement.Status == 1 {
		return fmt.Errorf("fail to confirm finish delivery arrangement: arrangement has not been finished")
	}
	if deliveryArrangement.Status == 3 {
		return fmt.Errorf("fail to confirm finish delivery arrangement: arrangement has been confirmed finished")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to confirm finish delivery arrangement: %v", err)

	}
	// 只能由运输方来修改状态
	if deliveryArrangement.HandlerOrg == clientOrgID {
		return fmt.Errorf("fail to confirm finish delivery arrangement: only owner can comfirm finish arrangement")
	}
	deliveryArrangement.Status = 3
	deliveryArrangement.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	deliveryArrangementBytes, err := json.Marshal(deliveryArrangement)
	if err != nil {
		return fmt.Errorf("fail to confirm finish delivery arrangement: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, deliveryArrangementBytes)
}

// GetDeliveryArrangement 获取DeliveryArrangement，args传入tradeID
func (t *CBPMChaincode) GetDeliveryArrangement(ctx contractapi.TransactionContextInterface, tradeID string) (*DeliveryArrangement, error) {
	deliveryArrangementBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to get delivery arrangement for trade %s: %v", tradeID, err)
	}
	if deliveryArrangementBytes == nil {
		return nil, fmt.Errorf("delivery arrangement for trade %s does not exist", tradeID)
	}

	var deliveryArrangement DeliveryArrangement
	err = json.Unmarshal(deliveryArrangementBytes, &deliveryArrangement)
	if err != nil {
		return nil, err
	}

	return &deliveryArrangement, nil
}

// GetAllDeliveryArrangements 获取所有DeliveryArrangement
func (t *CBPMChaincode) GetAllDeliveryArrangements(ctx contractapi.TransactionContextInterface) ([]*DeliveryArrangement, error) {
	queryString := "{\"selector\":{\"objectType\":\"DeliveryArrangement\"}}"
	return getQueryResultForQueryString(ctx, queryString)
}

// QueryDeliveryArrangements 查询满足条件的DeliveryArrangements，args传入查询语句
func (t *CBPMChaincode) QueryDeliveryArrangements(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryArrangement, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryArrangement, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*DeliveryArrangement, error) {
	var deliveryArrangementList []*DeliveryArrangement
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var deliveryArrangement DeliveryArrangement
		err = json.Unmarshal(queryResult.Value, &deliveryArrangement)
		if err != nil {
			return nil, err
		}
		deliveryArrangementList = append(deliveryArrangementList, &deliveryArrangement)
	}
	return deliveryArrangementList, nil
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	deliveryArrangements, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             deliveryArrangements,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetDeliveryArrangementsByRangeWithPagination performs a range query based on the start and end key,
// page size and a bookmark.
// The number of fetched records will be equal to or lesser than the page size.
// Paginated range queries are only valid for read only transactions.
// Example: Pagination with Range Query
func (t *CBPMChaincode) GetDeliveryArrangementsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) ([]*DeliveryArrangement, error) {

	resultsIterator, _, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// QueryOrdersWithPagination uses a query string, page size and a bookmark to perform a query
// for assets. Query string matching state database syntax is passed in and executed as is.
// The number of fetched records would be equal to or lesser than the specified page size.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the QueryAssetsForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// Paginated queries are only valid for read only transactions.
// Example: Pagination with Ad hoc Rich Query
func (t *CBPMChaincode) QueryOrdersWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	return getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

// GetOrderHistory returns the chain of custody for an asset since issuance.
func (t *CBPMChaincode) GetOrderHistory(ctx contractapi.TransactionContextInterface, tradeID string) ([]HistoryQueryResult, error) {
	log.Printf("GetOrderHistory for trade ID %v", tradeID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(tradeID)
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

		var deliveryArrangement DeliveryArrangement
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &deliveryArrangement)
			if err != nil {
				return nil, err
			}
		} else {
			deliveryArrangement = DeliveryArrangement{
				TradeID: tradeID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &deliveryArrangement,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}
	return records, nil
}

func (t *CBPMChaincode) deliveryArrangementExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	deliveryArrangementBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read delivery arrangement for trade %s from world state. %v", tradeID, err)
	}
	return deliveryArrangementBytes != nil, nil

}
func getClientOrgID(ctx contractapi.TransactionContextInterface, verifyOrg bool) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("fail getting client's orgID: %v", err)
	}
	if clientOrgID == "" {
		return "", fmt.Errorf("client ID is not set")
	}
	if verifyOrg {
		err = verifyClientOrgMatchesPeerOrg(clientOrgID)
		if err != nil {
			return "", err
		}
	}

	return clientOrgID, nil
}

// verifyClientOrgMatchesPeerOrg checks the client org id matches the peer org id.
func verifyClientOrgMatchesPeerOrg(clientOrgID string) error {
	peerOrgID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("fail getting peer's orgID: %v", err)
	}

	if clientOrgID != peerOrgID {
		return fmt.Errorf("client from org %s is not authorized to read or write private data from an org %s peer",
			clientOrgID,
			peerOrgID,
		)
	}
	return nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&CBPMChaincode{})
	if err != nil {
		log.Panicf("Error creating scchaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting scchaincode: %v", err)
	}
}
