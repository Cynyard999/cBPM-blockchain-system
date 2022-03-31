/*
 SPDX-License-Identifier: Apache-2.0
*/

/*
====CHAINCODE EXECUTION SAMPLES (CLI) ==================

==== Invoke assets ====
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["CreateAsset","asset1","blue","5","tom","35"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["CreateAsset","asset2","red","4","tom","50"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["CreateAsset","asset3","blue","6","tom","70"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["TransferAsset","asset2","jerry"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["TransferAssetByColor","blue","jerry"]}'
peer chaincode invoke -C myc1 -n asset_transfer -c '{"Args":["DeleteAsset","asset1"]}'

==== Query assets ====
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["ReadAsset","asset1"]}'
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["GetAssetsByRange","asset1","asset3"]}'
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["GetOrderHistory","asset1"]}'

Rich Query (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["QueryAssetsByOwner","tom"]}'
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["queryDeliveryOrders","{\"selector\":{\"owner\":\"tom\"}}"]}'

Rich Query with Pagination (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["QueryOrdersWithPagination","{\"selector\":{\"owner\":\"tom\"}}","3",""]}'

INDEXES TO SUPPORT COUCHDB RICH QUERIES

Indexes in CouchDB are required in order to make JSON queries efficient and are required for
any JSON query with a sort. Indexes may be packaged alongside
chaincode in a META-INF/statedb/couchdb/indexes directory. Each index must be defined in its own
text file with extension *.json with the index definition formatted in JSON following the
CouchDB index JSON syntax as documented at:
http://docs.couchdb.org/en/2.3.1/api/database/find.html#db-index

This asset transfer ledger example chaincode demonstrates a packaged
index which you can find in META-INF/statedb/couchdb/indexes/indexOwner.json.

If you have access to the your peer's CouchDB state database in a development environment,
you may want to iteratively test various indexes in support of your chaincode queries.  You
can use the CouchDB Fauxton interface or a command line curl utility to create and update
indexes. Then once you finalize an index, include the index definition alongside your
chaincode in the META-INF/statedb/couchdb/indexes directory, for packaging and deployment
to managed environments.

In the examples below you can find index definitions that support asset transfer ledger
chaincode queries, along with the syntax that you can use in development environments
to create the indexes in the CouchDB Fauxton interface or a curl command line utility.


Index for docType, owner.

Example curl command line to define index in the CouchDB channel_chaincode database
curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[\"docType\",\"owner\"]},\"name\":\"indexOwner\",\"ddoc\":\"indexOwnerDoc\",\"type\":\"json\"}" http://hostname:port/myc1_assets/_index


Index for docType, owner, size (descending order).

Example curl command line to define index in the CouchDB channel_chaincode database:
curl -i -X POST -H "Content-Type: application/json" -d "{\"index\":{\"fields\":[{\"size\":\"desc\"},{\"docType\":\"desc\"},{\"owner\":\"desc\"}]},\"ddoc\":\"indexSizeSortDoc\", \"name\":\"indexSizeSortDesc\",\"type\":\"json\"}" http://hostname:port/myc1_assets/_index

Rich Query with index design doc and index name specified (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["queryDeliveryOrders","{\"selector\":{\"docType\":\"asset\",\"owner\":\"tom\"}, \"use_index\":[\"_design/indexOwnerDoc\", \"indexOwner\"]}"]}'

Rich Query with index design doc specified only (Only supported if CouchDB is used as state database):
peer chaincode query -C myc1 -n asset_transfer -c '{"Args":["queryDeliveryOrders","{\"selector\":{\"docType\":{\"$eq\":\"asset\"},\"owner\":{\"$eq\":\"tom\"},\"size\":{\"$gt\":0}},\"fields\":[\"docType\",\"owner\",\"size\"],\"sort\":[{\"size\":\"desc\"}],\"use_index\":\"_design/indexSizeSortDoc\"}"]}'
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

const index = "tradeID~orderID"

// SimpleChaincode implements the fabric-contract-api-go programming model
type SimpleChaincode struct {
	contractapi.Contract
}

type DeliveryOrder struct {
	ObjectType string `json:"docType"`
	TradeID    string `json:"tradeID"` // 生产商下单时生成的ID，用于表明一次流程
	OrderID    string `json:"orderID"` // 自动生成
	AssetID    string `json:"assetID"`
	Quantity   int    `json:"quantity"`
	StartPlace string `json:"startPlace"`
	EndPlace   string `json:"endPlace"`
	Status     int    `json:"status"` // 0: 未处理 1：开始运输 2: 已完成
	OwnerOrg   string `json:"ownerOrg"`
	HandleOrg  string `json:"handlerOrg"`
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

//deleteDeliveryOrder // Status为0的时候才能删除
//queryDeliveryOrder
//queryAllDeliveryOrder

// CreateAsset initializes a new asset in the ledger
func (t *SimpleChaincode) createDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string, assetID string, quantity int, startPlace string, endPlace string) error {
	exists, err := t.deliveryOrderExists(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if exists {
		return fmt.Errorf("asset already exists: %s", assetID)
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	deliveryOrder := &DeliveryOrder{
		ObjectType: "deliveryOrder",
		TradeID:    tradeID,
		OrderID:    "order_" + tradeID,
		AssetID:    assetID,
		Quantity:   quantity,
		StartPlace: startPlace,
		EndPlace:   endPlace,
		Status:     0,
		OwnerOrg:   clientOrgID,
		HandleOrg:  "",
	}

	deliveryOrderBytes, err := json.Marshal(deliveryOrder)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(tradeID, deliveryOrderBytes)
	if err != nil {
		return err
	}
	colorNameIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{deliveryOrder.TradeID, deliveryOrder.OrderID})
	if err != nil {
		return err
	}
	//  Save index entry to world state. Only the key name is needed, no need to store a duplicate copy of the asset.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	return ctx.GetStub().PutState(colorNameIndexKey, value)
}

func (t *SimpleChaincode) DeleteDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryOrder, err := t.getDeliveryOrder(ctx, tradeID)
	if err != nil {
		return err
	}

	err = ctx.GetStub().DelState(tradeID)
	if err != nil {
		return fmt.Errorf("failed to delete delivery order for trade %s: %v", tradeID, err)
	}

	tradeIDOrderIDIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{deliveryOrder.TradeID, deliveryOrder.OrderID})
	if err != nil {
		return err
	}
	// Delete index entry
	return ctx.GetStub().DelState(tradeIDOrderIDIndexKey)
}

func (t *SimpleChaincode) pickDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryOrder, err := t.getDeliveryOrder(ctx, tradeID)
	if err != nil {
		return err
	}
	if deliveryOrder.HandleOrg != "" {
		return fmt.Errorf("delivery order for trade #{tradeID} has been picked by %s", deliveryOrder.HandleOrg)
	}
	if deliveryOrder.Status != 0 {
		return fmt.Errorf("delivery order for trade #{tradeID} has been picked by %s and current status is %d", deliveryOrder.HandleOrg, deliveryOrder.Status)
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	deliveryOrder.HandleOrg = clientOrgID
	deliveryOrder.Status = 1
	deliveryOrderBytes, err := json.Marshal(deliveryOrder)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(tradeID, deliveryOrderBytes)
}

func (t *SimpleChaincode) updateDeliveryOrderStatus(ctx contractapi.TransactionContextInterface, tradeID string, status int) error {
	deliveryOrder, err := t.getDeliveryOrder(ctx, tradeID)
	if err != nil {
		return err
	}
	if deliveryOrder.HandleOrg == "" {
		return fmt.Errorf("delivery order for trade #{tradeID} has not been picked by any organizations")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return err
	}
	// 只能由运输方来修改状态
	if deliveryOrder.HandleOrg != clientOrgID {
		return fmt.Errorf("unauthorized handler #{clientOrgID}")
	}
	deliveryOrder.Status = status
	deliveryOrderBytes, err := json.Marshal(deliveryOrder)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(tradeID, deliveryOrderBytes)
}

func (t *SimpleChaincode) getDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) (*DeliveryOrder, error) {
	deliveryOrderBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery order for trade %s: %v", tradeID, err)
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

func (t *SimpleChaincode) queryDeliveryOrders(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryOrder, error) {
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

	assets, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             assets,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// GetAssetsByRangeWithPagination performs a range query based on the start and end key,
// page size and a bookmark.
// The number of fetched records will be equal to or lesser than the page size.
// Paginated range queries are only valid for read only transactions.
// Example: Pagination with Range Query
func (t *SimpleChaincode) GetAssetsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int, bookmark string) ([]*DeliveryOrder, error) {

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
func (t *SimpleChaincode) QueryOrdersWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int, bookmark string) (*PaginatedQueryResult, error) {

	return getQueryResultForQueryStringWithPagination(ctx, queryString, int32(pageSize), bookmark)
}

// GetOrderHistory returns the chain of custody for an asset since issuance.
func (t *SimpleChaincode) GetOrderHistory(ctx contractapi.TransactionContextInterface, tradeID string) ([]HistoryQueryResult, error) {
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

func (t *SimpleChaincode) deliveryOrderExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	deliveryOrderBytes, err := ctx.GetStub().GetState(tradeID)
	if err != nil {
		return false, fmt.Errorf("failed to read asset %s from world state. %v", tradeID, err)
	}
	return deliveryOrderBytes != nil, nil

}
func getClientOrgID(ctx contractapi.TransactionContextInterface, verifyOrg bool) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed getting client's orgID: %v", err)
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
		return fmt.Errorf("failed getting peer's orgID: %v", err)
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
	chaincode, err := contractapi.NewChaincode(&SimpleChaincode{})
	if err != nil {
		log.Panicf("Error creating asset chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting asset chaincode: %v", err)
	}
}
