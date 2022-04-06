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

type DeliveryDetail struct {
	ObjectType string `json:"objectType"`
	TradeID    string `json:"tradeID"` // 生产商下单时生成的ID，用于表明一次流程
	AssetName  string `json:"assetName"`
	StartPlace string `json:"startPlace"`
	EndPlace   string `json:"endPlace"`
	Contact    string `json:"contact"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	Status     int    `json:"status"` // 0: 未处理 1：开始运输 2: 运输已完成
	OwnerOrg   string `json:"ownerOrg"`
	Note       string `json:"note"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *DeliveryDetail `json:"record"`
	TxId      string          `json:"txId"`
	Timestamp time.Time       `json:"timestamp"`
	IsDelete  bool            `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*DeliveryDetail `json:"records"`
	FetchedRecordsCount int32             `json:"fetchedRecordsCount"`
	Bookmark            string            `json:"bookmark"`
}

func (t *CBPMChaincode) CreateDeliveryDetail(ctx contractapi.TransactionContextInterface) error {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientDetailJSON, ok := transMap["detail"]
	if !ok {
		return fmt.Errorf("detail not found in the transient map")
	}
	type detailTransientInput struct {
		TradeID    string `json:"tradeID"`
		AssetName  string `json:"assetName"`
		StartPlace string `json:"startPlace"`
		EndPlace   string `json:"endPlace"`
		Contact    string `json:"contact"`
		Note       string `json:"note"`
	}
	var detailInput detailTransientInput
	err = json.Unmarshal(transientDetailJSON, &detailInput)
	if err != nil {
		return fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	if len(detailInput.TradeID) == 0 {
		return fmt.Errorf("trade ID must be a non-empty string")
	}
	if len(detailInput.AssetName) == 0 {
		return fmt.Errorf("asset name must be a non-empty string")
	}
	if len(detailInput.StartPlace) == 0 {
		return fmt.Errorf("start place must be a non-empty string")
	}
	if len(detailInput.EndPlace) == 0 {
		return fmt.Errorf("end place must be a non-empty string")
	}
	if len(detailInput.Contact) == 0 {
		return fmt.Errorf("contact must be a non-empty string")
	}

	exists, err := t.DeliveryDetailExists(ctx, detailInput.TradeID)
	if err != nil {
		return fmt.Errorf("fail to create delivery detail: %v", err)
	}
	if exists {
		return fmt.Errorf("fail to create delivery detail: detail for trade %s already exists", detailInput.TradeID)
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to create delivery detail: %v", err)
	}

	deliveryDetail := &DeliveryDetail{
		ObjectType: "deliveryDetail",
		TradeID:    detailInput.TradeID,
		AssetName:  detailInput.AssetName,
		StartPlace: detailInput.StartPlace,
		EndPlace:   detailInput.EndPlace,
		Contact:    detailInput.Contact,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     0,
		OwnerOrg:   clientOrgID,
		Note:       detailInput.Note,
	}

	deliveryDetailBytes, err := json.Marshal(deliveryDetail)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(detailInput.TradeID, deliveryDetailBytes)
}

func (t *CBPMChaincode) DeleteDeliveryDetail(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryDetail, err := t.GetDeliveryDetail(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to delete delivery detail: %v", err)
	}
	if deliveryDetail.Status != 0 {
		return fmt.Errorf("fail to delete delivery detail: cannot delete detail in progress")
	}
	return ctx.GetStub().DelState(tradeID)
}

func (t *CBPMChaincode) HandleDeliveryDetail(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryDetail, err := t.GetDeliveryDetail(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to handle delivery detail: %v", err)
	}
	if deliveryDetail.Status != 0 {
		return fmt.Errorf("fail to handle delivery detail: delivery detail(status: %d) has been handled", deliveryDetail.Status)
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to handle delivery detail: %v", err)
	}
	if deliveryDetail.OwnerOrg != clientOrgID {
		return fmt.Errorf("fail to handle delivery detail: only owner can handle")
	}
	deliveryDetail.Status = 1
	deliveryDetail.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	deliveryDetailBytes, err := json.Marshal(deliveryDetail)
	if err != nil {
		return fmt.Errorf("fail to handle delivery detail: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, deliveryDetailBytes)
}

func (t *CBPMChaincode) FinishDeliveryDetail(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryDetail, err := t.GetDeliveryDetail(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to finish delivery detail: %v", err)
	}
	if deliveryDetail.Status == 0 {
		return fmt.Errorf("fail to finish delivery detail: detail for trade #{tradeID} has not been handled")
	}
	if deliveryDetail.Status == 1 {
		return fmt.Errorf("fail to finish delivery detail: detail for trade #{tradeID} has been finished")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to finish delivery detail: %v", err)
	}
	if deliveryDetail.OwnerOrg != clientOrgID {
		return fmt.Errorf("fail to finish delivery detail: only owner can finish")
	}
	// 只能由运输方来修改状态
	deliveryDetail.Status = 2
	deliveryDetail.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	deliveryDetailBytes, err := json.Marshal(deliveryDetail)
	if err != nil {
		return fmt.Errorf("fail to finish delivery detail: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, deliveryDetailBytes)
}

func (t *CBPMChaincode) GetDeliveryDetail(ctx contractapi.TransactionContextInterface, tradeID string) (*DeliveryDetail, error) {
	deliveryDetailBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to get delivery detail for trade %s: %v", tradeID, err)
	}
	if deliveryDetailBytes == nil {
		return nil, fmt.Errorf("delivery detail for trade %s does not exist", tradeID)
	}

	var deliveryDetail DeliveryDetail
	err = json.Unmarshal(deliveryDetailBytes, &deliveryDetail)
	if err != nil {
		return nil, err
	}

	return &deliveryDetail, nil
}

func (t *CBPMChaincode) GetAllDeliveryDetails(ctx contractapi.TransactionContextInterface) ([]*DeliveryDetail, error) {
	queryString := "{\"selector\":{\"objectType\":\"DeliveryDetail\"}}"
	return getQueryResultForQueryString(ctx, queryString)
}

func (t *CBPMChaincode) QueryDeliveryDetails(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryDetail, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryDetail, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*DeliveryDetail, error) {
	var deliveryDetailList []*DeliveryDetail
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var deliveryDetail DeliveryDetail
		err = json.Unmarshal(queryResult.Value, &deliveryDetail)
		if err != nil {
			return nil, err
		}
		deliveryDetailList = append(deliveryDetailList, &deliveryDetail)
	}
	return deliveryDetailList, nil
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	deliveryDetails, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             deliveryDetails,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetDeliveryDetailsByRangeWithPagination performs a range query based on the start and end key,
// page size and a bookmark.
// The number of fetched records will be equal to or lesser than the page size.
// Paginated range queries are only valid for read only transactions.
// Example: Pagination with Range Query
func (t *CBPMChaincode) GetDeliveryDetailsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) ([]*DeliveryDetail, error) {

	resultsIterator, _, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// QueryDetailsWithPagination uses a query string, page size and a bookmark to perform a query
// for assets. Query string matching state database syntax is passed in and executed as is.
// The number of fetched records would be equal to or lesser than the specified page size.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the QueryAssetsForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// Paginated queries are only valid for read only transactions.
// Example: Pagination with Ad hoc Rich Query
func (t *CBPMChaincode) QueryDetailsWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	return getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

// GetDetailHistory returns the chain of custody for an asset since issuance.
func (t *CBPMChaincode) GetDetailHistory(ctx contractapi.TransactionContextInterface, tradeID string) ([]HistoryQueryResult, error) {
	log.Printf("GetDetailHistory for trade ID %v", tradeID)

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

		var deliveryDetail DeliveryDetail
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &deliveryDetail)
			if err != nil {
				return nil, err
			}
		} else {
			deliveryDetail = DeliveryDetail{
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
			Record:    &deliveryDetail,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}
	return records, nil
}

func (t *CBPMChaincode) DeliveryDetailExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	deliveryDetailBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read delivery detail for trade %s from world state. %v", tradeID, err)
	}
	return deliveryDetailBytes != nil, nil

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
