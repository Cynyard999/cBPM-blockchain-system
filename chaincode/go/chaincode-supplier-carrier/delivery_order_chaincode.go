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

type DeliveryOrder struct {
	ObjectType string `json:"objectType"`
	TradeID    string `json:"tradeID"` // 生产商下单时生成的ID，用于表明一次流程
	AssetName  string `json:"assetName"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	Status     int    `json:"status"` // 0: 未处理 1：开始运输 2: 运输商已完成
	OwnerOrg   string `json:"ownerOrg"`
	HandlerOrg string `json:"handlerOrg"`
	Note       string `json:"note"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *DeliveryOrder `json:"record"`
	TxId      string         `json:"txId"`
	Timestamp time.Time      `json:"timestamp"`
	IsDelete  bool           `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*DeliveryOrder `json:"records"`
	FetchedRecordsCount int32            `json:"fetchedRecordsCount"`
	Bookmark            string           `json:"bookmark"`
}

func (t *CBPMChaincode) CreateDeliveryOrder(ctx contractapi.TransactionContextInterface) error {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientOrderJSON, ok := transMap["order"]
	if !ok {
		return fmt.Errorf("order not found in the transient map")
	}
	type orderTransientInput struct {
		TradeID   string `json:"tradeID"`
		AssetName string `json:"assetName"`
		Note      string `json:"note"`
	}
	var orderInput orderTransientInput
	err = json.Unmarshal(transientOrderJSON, &orderInput)
	if err != nil {
		return fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	if len(orderInput.TradeID) == 0 {
		return fmt.Errorf("trade ID must be a non-empty string")
	}
	if len(orderInput.AssetName) == 0 {
		return fmt.Errorf("asset name must be a non-empty string")
	}
	exists, err := t.deliveryOrderExists(ctx, orderInput.TradeID)
	if err != nil {
		return fmt.Errorf("fail to create delivery order: %v", err)
	}
	if exists {
		return fmt.Errorf("fail to create delivery order: order for trade %s already exists", orderInput.TradeID)
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to create delivery order: %v", err)
	}

	deliveryOrder := &DeliveryOrder{
		ObjectType: "deliveryOrder",
		TradeID:    orderInput.TradeID,
		AssetName:  orderInput.AssetName,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     0,
		OwnerOrg:   clientOrgID,
		HandlerOrg: "",
		Note:       orderInput.Note,
	}

	deliveryOrderBytes, err := json.Marshal(deliveryOrder)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(orderInput.TradeID, deliveryOrderBytes)
}

func (t *CBPMChaincode) DeleteDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryOrder, err := t.GetDeliveryOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to delete delivery order: %v", err)
	}
	if deliveryOrder.Status != 0 {
		return fmt.Errorf("fail to delete delivery order: cannot delete order in progress")
	}
	return ctx.GetStub().DelState(tradeID)
}

func (t *CBPMChaincode) HandleDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryOrder, err := t.GetDeliveryOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to handle delivery order: %v", err)
	}
	if deliveryOrder.HandlerOrg != "" {
		return fmt.Errorf("fail to handle delivery order: delivery order for trade #{tradeID} has been handled")
	}
	if deliveryOrder.Status != 0 {
		return fmt.Errorf("fail to handle delivery order: delivery order(status: %d) for trade #{tradeID} has been handled", deliveryOrder.Status)
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to handle delivery order: %v", err)
	}
	if deliveryOrder.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to handle delivery order: cannot handle as owner")
	}
	deliveryOrder.HandlerOrg = clientOrgID
	deliveryOrder.Status = 1
	deliveryOrder.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	deliveryOrderBytes, err := json.Marshal(deliveryOrder)
	if err != nil {
		return fmt.Errorf("fail to handle delivery order: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, deliveryOrderBytes)
}

func (t *CBPMChaincode) FinishDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryOrder, err := t.GetDeliveryOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to finish delivery order: %v", err)
	}
	if deliveryOrder.Status == 0 {
		return fmt.Errorf("fail to finish delivery order: order for trade #{tradeID} has not been handled")
	}
	if deliveryOrder.Status == 1 {
		return fmt.Errorf("fail to finish delivery order: order for trade #{tradeID} has been finished")
	}
	if deliveryOrder.HandlerOrg == "" {
		return fmt.Errorf("fail to finish delivery order: no handler is specified")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to finish delivery order: %v", err)
	}
	if deliveryOrder.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to finish delivery order: cannot finish as owner")
	}
	// 只能由运输方来修改状态
	if deliveryOrder.HandlerOrg != clientOrgID {
		return fmt.Errorf("fail to finish delivery order: unauthorized handler #{clientOrgID}}")
	}
	deliveryOrder.Status = 2
	deliveryOrder.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	deliveryOrderBytes, err := json.Marshal(deliveryOrder)
	if err != nil {
		return fmt.Errorf("fail to finish delivery order: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, deliveryOrderBytes)
}

//func (t *CBPMChaincode) ConfirmFinishDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
//	deliveryOrder, err := t.GetDeliveryOrder(ctx, tradeID)
//	if err != nil {
//		return fmt.Errorf("fail to confirm finish delivery order: %v", err)
//	}
//	if deliveryOrder.Status == 0 {
//		return fmt.Errorf("fail to confirm finish delivery order: order has not been handled")
//	}
//	if deliveryOrder.Status == 1 {
//		return fmt.Errorf("fail to confirm finish delivery order: order has not been finished")
//	}
//	if deliveryOrder.Status == 3 {
//		return fmt.Errorf("fail to confirm finish delivery order: order has been confirmed finished")
//	}
//	clientOrgID, err := getClientOrgID(ctx, false)
//	if err != nil {
//		return fmt.Errorf("fail to confirm finish delivery order: %v", err)
//
//	}
//	// 只能由运输方来修改状态
//	if deliveryOrder.HandlerOrg == clientOrgID {
//		return fmt.Errorf("fail to confirm finish delivery order: only owner can comfirm finish order")
//	}
//	deliveryOrder.Status = 3
//	deliveryOrder.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
//	deliveryOrderBytes, err := json.Marshal(deliveryOrder)
//	if err != nil {
//		return fmt.Errorf("fail to confirm finish delivery order: %v", err)
//	}
//	return ctx.GetStub().PutState(tradeID, deliveryOrderBytes)
//}

func (t *CBPMChaincode) GetDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) (*DeliveryOrder, error) {
	deliveryOrderBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to get delivery order for trade %s: %v", tradeID, err)
	}
	if deliveryOrderBytes == nil {
		return nil, fmt.Errorf("delivery order for trade %s does not exist", tradeID)
	}

	var deliveryOrder DeliveryOrder
	err = json.Unmarshal(deliveryOrderBytes, &deliveryOrder)
	if err != nil {
		return nil, err
	}

	return &deliveryOrder, nil
}

func (t *CBPMChaincode) GetAllDeliveryOrders(ctx contractapi.TransactionContextInterface) ([]*DeliveryOrder, error) {
	queryString := "{\"selector\":{\"objectType\":\"DeliveryOrder\"}}"
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *CBPMChaincode) QueryDeliveryOrders(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryOrder, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryOrder, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*DeliveryOrder, error) {
	var deliveryOrderList []*DeliveryOrder
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var deliveryOrder DeliveryOrder
		err = json.Unmarshal(queryResult.Value, &deliveryOrder)
		if err != nil {
			return nil, err
		}
		deliveryOrderList = append(deliveryOrderList, &deliveryOrder)
	}
	return deliveryOrderList, nil
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	deliveryOrders, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             deliveryOrders,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetDeliveryOrdersByRangeWithPagination performs a range query based on the start and end key,
// page size and a bookmark.
// The number of fetched records will be equal to or lesser than the page size.
// Paginated range queries are only valid for read only transactions.
// Example: Pagination with Range Query
func (t *CBPMChaincode) GetDeliveryOrdersByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) ([]*DeliveryOrder, error) {

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

		var deliveryOrder DeliveryOrder
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &deliveryOrder)
			if err != nil {
				return nil, err
			}
		} else {
			deliveryOrder = DeliveryOrder{
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
			Record:    &deliveryOrder,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}
	return records, nil
}

func (t *CBPMChaincode) deliveryOrderExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	deliveryOrderBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read delivery order for trade %s from world state. %v", tradeID, err)
	}
	return deliveryOrderBytes != nil, nil

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
