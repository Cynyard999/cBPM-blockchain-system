package main

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"time"
)

type CBPMChaincode struct {
	contractapi.Contract
}

type SupplyAsset struct {
	ObjectType        string  `json:"objectType"`
	AssetID           string  `json:"assetID"`
	AssetName         string  `json:"assetName"`
	AssetPrice        float32 `json:"assetPrice"`
	ShippingAddress   string  `json:"shippingAddress"`
	OwnerOrg          string  `json:"ownerOrg"` // 供货商的Org
	PublicDescription string  `json:"publicDescription"`
}

type SupplyOrder struct {
	ObjectType      string  `json:"objectType"`
	TradeID         string  `json:"tradeID"`
	AssetID         string  `json:"assetID"`
	AssetName       string  `json:"assetName"`
	AssetPrice      float32 `json:"assetPrice"`
	ShippingAddress string  `json:"shippingAddress"`
	Quantity        int     `json:"quantity"`
	TotalPrice      float32 `json:"totalPrice"` // 自动生成
	Status          int     `json:"status"`     // 0: 未开始 1: 中间商开始处理 2: 中间商确认已完成 3: 生产商确认已完成
	CreateTime      string  `json:"createTime"` // 自动生成
	UpdateTime      string  `json:"updateTime"` // 自动生成
	OwnerOrg        string  `json:"ownerOrg"`   // 限制修改的权限
	HandlerOrg      string  `json:"handlerOrg"` // 限制修改的权限
	Note            string  `json:"note"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *SupplyOrder `json:"record"`
	TxId      string       `json:"txId"`
	Timestamp time.Time    `json:"timestamp"`
	IsDelete  bool         `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*SupplyOrder `json:"records"`
	FetchedRecordsCount int32          `json:"fetchedRecordsCount"`
	Bookmark            string         `json:"bookmark"`
}

// CreateSupplyAsset 创建asset，需要transient传入assetName，assetPrice，shippingAddress，publicDescription 会生成一个UUID表示asset的名称，最终返回创建的asset对象
func (t *CBPMChaincode) CreateSupplyAsset(ctx contractapi.TransactionContextInterface) (*SupplyAsset, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientAssetJSON, ok := transMap["supplyAsset"]
	if !ok {
		return nil, fmt.Errorf("supplyAsset field not found in the transient map")
	}
	type assetTransientInput struct {
		AssetName         string  `json:"assetName"`
		AssetPrice        float32 `json:"assetPrice"`
		ShippingAddress   string  `json:"shippingAddress"`
		PublicDescription string  `json:"publicDescription"`
	}
	var assetInput assetTransientInput
	err = json.Unmarshal(transientAssetJSON, &assetInput)
	if err != nil {
		return nil, fmt.Errorf("fail to create supplyAsset: fail to unmarshal JSON: %s", err.Error())
	}
	if len(assetInput.ShippingAddress) == 0 {
		return nil, fmt.Errorf("fail to create supplyAsset: shipping address must be a non-empty string")
	}
	if len(assetInput.AssetName) == 0 {
		return nil, fmt.Errorf("fail to create supplyAsset: supplyAsset name must be a non-empty string")
	}
	if assetInput.AssetPrice <= 0 {
		return nil, fmt.Errorf("fail to create supplyAsset: supplyAsset price field must be a positive number")
	}

	exists, err := t.assetNameExists(ctx, assetInput.AssetName)
	if err != nil {
		return nil, fmt.Errorf("fail to create SupplyAsset: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("fail to create SupplyAsset: supplyAsset name already exists")
	}

	assetIDUUID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("fail to create supplyAsset: fail to generate supplyAsset ID: %v", err)
	}
	assetID := assetIDUUID.String()

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("fail to create supplyAsset: %v", err)
	}

	// create supplyAsset
	supplyAsset := &SupplyAsset{
		ObjectType:        "SupplyAsset",
		AssetID:           assetID,
		AssetName:         assetInput.AssetName,
		AssetPrice:        assetInput.AssetPrice,
		ShippingAddress:   assetInput.ShippingAddress,
		OwnerOrg:          clientOrgID,
		PublicDescription: assetInput.PublicDescription,
	}
	assetJSONasBytes, err := json.Marshal(supplyAsset)
	if err != nil {
		return nil, fmt.Errorf("fail to create supplyAsset: %v", err)
	}

	err = ctx.GetStub().PutState(supplyAsset.AssetID, assetJSONasBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create supplyAsset: %v", err)
	}
	return supplyAsset, nil
}

// UpdateSupplyAsset 更新指定的asset，需要args传入assetID，assetName，assetPrice，shippingAddress，desc，只有owner能够更新
func (t *CBPMChaincode) UpdateSupplyAsset(ctx contractapi.TransactionContextInterface, assetID string, assetName string, assetPrice float32, shippingAddress string, desc string) error {
	asset, err := t.GetSupplyAsset(ctx, assetID)
	if err != nil {
		return fmt.Errorf("fail to update asset: %v", err)
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to update asset: %v", err)
	}
	if asset.OwnerOrg != clientOrgID {
		return fmt.Errorf("fail to update SupplyAsset: unauthorized updater %s", clientOrgID)
	}
	asset.ShippingAddress = shippingAddress
	asset.AssetPrice = assetPrice
	asset.AssetName = assetName
	asset.PublicDescription = desc
	newAssetBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("fail to update asset: %v", err)
	}
	return ctx.GetStub().PutState(assetID, newAssetBytes)
}

// DeleteSupplyAsset 删除指定的asset，需要args传入assetID
func (t *CBPMChaincode) DeleteSupplyAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	exist, err := t.assetExists(ctx, assetID)
	if !exist {
		return fmt.Errorf("fail to delete asset: asset does not exist")
	}
	if err != nil {
		return fmt.Errorf("fail to delete asset: %v", err)
	}
	return ctx.GetStub().DelState(assetID)
}

// GetSupplyAsset 获取指定的asset，需要args传入assetID
func (t *CBPMChaincode) GetSupplyAsset(ctx contractapi.TransactionContextInterface, assetID string) (*SupplyAsset, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"SupplyAsset\",\"assetID\":\"%s\"}}", assetID)
	queryResults, err := t.getAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, fmt.Errorf("fail to get asset: %v", err)
	}
	if len(queryResults) == 0 {
		return nil, fmt.Errorf("fail to get asset: %s does not exist", assetID)
	}
	return queryResults[0], nil
}

// GetAllSupplyAssets 获取所有的asset
func (t *CBPMChaincode) GetAllSupplyAssets(ctx contractapi.TransactionContextInterface) ([]*SupplyAsset, error) {

	queryString := "{\"selector\":{\"objectType\":\"SupplyAsset\"}}"

	queryResults, err := t.getAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// QuerySupplyAssets 查询assets，需要args传入查询语句
func (t *CBPMChaincode) QuerySupplyAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyAsset, error) {
	queryResults, err := t.getAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

func (t *CBPMChaincode) assetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"SupplyAsset\",\"assetID\":\"%s\"}}", assetID)
	queryResults, err := t.getAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return false, fmt.Errorf("fail to check whether asset exists: %v", err)
	}
	if len(queryResults) == 0 {
		return false, nil
	}
	return true, nil
}

func (t *CBPMChaincode) assetNameExists(ctx contractapi.TransactionContextInterface, assetName string) (bool, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"SupplyAsset\",\"assetName\":\"%s\"}}", assetName)
	queryResults, err := t.getAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return false, fmt.Errorf("fail to check whether asset name exists: %v", err)
	}
	if len(queryResults) == 0 {
		return false, nil
	}
	return true, nil
}

// CreateSupplyOrder 创建supplyOrder，需要transient传入tradeID,assetID,quantity,note，返回创建好的对象
func (t *CBPMChaincode) CreateSupplyOrder(ctx contractapi.TransactionContextInterface) (*SupplyOrder, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientOrderJSON, ok := transMap["order"]
	if !ok {
		return nil, fmt.Errorf("order not found in the transient map")
	}
	type orderTransientInput struct {
		TradeID  string `json:"tradeID"`
		AssetID  string `json:"assetID"`
		Quantity int    `json:"quantity"`
		Note     string `json:"note"`
	}
	var orderInput orderTransientInput
	err = json.Unmarshal(transientOrderJSON, &orderInput)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	// check input
	if len(orderInput.TradeID) == 0 {
		return nil, fmt.Errorf("trade ID must be a non-empty string")
	}
	if len(orderInput.AssetID) == 0 {
		return nil, fmt.Errorf("asset ID must be a non-empty string")
	}
	if orderInput.Quantity <= 0 {
		return nil, fmt.Errorf("asset quantity field must be a positive number")
	}
	exist, err := t.supplyOrderExists(ctx, orderInput.TradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to create supply order: %v", err)
	}
	if exist {
		return nil, fmt.Errorf("fail to create supply order: order for trade %s already exists", orderInput.TradeID)
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("fail to create supply order: %v", err)
	}
	asset, err := t.GetSupplyAsset(ctx, orderInput.AssetID)
	if err != nil {
		return nil, fmt.Errorf("fail to create supply order: %v", err)
	}
	// create order
	order := &SupplyOrder{
		ObjectType:      "SupplyOrder",
		TradeID:         orderInput.TradeID,
		AssetID:         orderInput.AssetID,
		AssetName:       asset.AssetName,
		AssetPrice:      asset.AssetPrice,
		ShippingAddress: asset.ShippingAddress,
		Quantity:        orderInput.Quantity,
		TotalPrice:      asset.AssetPrice * (float32(orderInput.Quantity)),
		Status:          0,
		CreateTime:      time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime:      time.Now().Format("2006-01-02 15:04:05"),
		HandlerOrg:      "",
		OwnerOrg:        clientOrgID,
		Note:            orderInput.Note,
	}
	orderJSONasBytes, err := json.Marshal(order)
	err = ctx.GetStub().PutState(order.TradeID, orderJSONasBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create supply order: %s", err.Error())
	}
	return order, nil
}

// GetSupplyOrder 获取supplyOrder，需要args传入tradeID
func (t *CBPMChaincode) GetSupplyOrder(ctx contractapi.TransactionContextInterface, tradeID string) (*SupplyOrder, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"SupplyOrder\",\"tradeID\":\"%s\"}}", tradeID)
	queryResults, err := t.getOrderQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	if len(queryResults) == 0 {
		return nil, fmt.Errorf("fail to get supply order for tradeID: %s does not exist", tradeID)
	}
	return queryResults[0], nil
}

// GetAllSupplyOrders 获取所有supplyOrders
func (t *CBPMChaincode) GetAllSupplyOrders(ctx contractapi.TransactionContextInterface) ([]*SupplyOrder, error) {
	queryString := "{\"selector\":{\"objectType\":\"SupplyOrder\"}}"

	queryResults, err := t.getOrderQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// QuerySupplyOrders 查询满足条件的supplyOrders，需要args传入查询条件
func (t *CBPMChaincode) QuerySupplyOrders(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyOrder, error) {
	queryResults, err := t.getOrderQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// DeleteSupplyOrder 删除supplyOrder，只允许owner删除
func (t *CBPMChaincode) DeleteSupplyOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	exists, err := t.supplyOrderExists(ctx, tradeID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("fail to delete supply order: order for trade #{tradeID} does not exist")
	}
	// TODO 只允许owner删除
	return ctx.GetStub().DelState(tradeID)
}

// HandleSupplyOrder 供应商接单，需要args传入tradeID，Owner并不能进行这一步操作
func (t *CBPMChaincode) HandleSupplyOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	order, err := t.GetSupplyOrder(ctx, tradeID)
	if err != nil {
		return err
	}
	if order.Status != 0 {
		return fmt.Errorf("fail to handle supply order: order(status: %d) for trade #{tradeID} has been handled", order.Status)
	}
	if order.HandlerOrg != "" {
		return fmt.Errorf("fail to handle supply order: order for trade #{tradeID} has been handled by some org")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to handle supply order: %v", err)
	}
	if order.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to handle supply order: cannot handle as owner")
	}
	order.HandlerOrg = clientOrgID
	order.Status = 1
	order.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("fail to handle supply order: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, orderBytes)
}

// FinishSupplyOrder Handler完成订单，需要args传入TradeID，Owner并不能进行这一步操作
func (t *CBPMChaincode) FinishSupplyOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	order, err := t.GetSupplyOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to finish supply order: %v", err)
	}
	if order.Status == 0 {
		return fmt.Errorf("fail to finish supply order: order for trade #{tradeID} has not been handled")
	}
	if order.Status == 2 {
		return fmt.Errorf("fail to finish supply order: order for trade #{tradeID} has been finished")
	}
	if order.HandlerOrg == "" {
		return fmt.Errorf("fail to finish supply order: no handler is specified")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to finish supply order: %v", err)
	}
	if order.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to finish supply order: cannot finish as owner")
	}
	if order.HandlerOrg != clientOrgID {
		return fmt.Errorf("fail to finish supply order: cannot finish by other org: %s", clientOrgID)
	}
	order.Status = 2
	order.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("fail to finish supply order: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, orderBytes)
}

// ConfirmFinishSupplyOrder Owner确认供应商完成订单，需要args传入tradeID，
func (t *CBPMChaincode) ConfirmFinishSupplyOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	order, err := t.GetSupplyOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to confirm finish supply order: %v", err)
	}
	if order.Status == 0 {
		return fmt.Errorf("fail to confirm finish supply order: order for trade #{tradeID} has not been handled")
	}
	if order.Status == 1 {
		return fmt.Errorf("fail to confirm finish supply order: order for trade #{tradeID} has not been finished")
	}
	if order.Status == 3 {
		return fmt.Errorf("fail to confirm finish supply order: order for trade #{tradeID} has been confirmed finished")
	}
	if order.HandlerOrg == "" {
		return fmt.Errorf("fail to confirm finish supply order: no handler is specified")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to confirm finish supply order: %v", err)
	}
	if order.OwnerOrg != clientOrgID {
		return fmt.Errorf("fail to confirm finish supply order: only owner can comfirm finish order")
	}
	order.Status = 3
	order.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("fail to confirm finish supply order: %v", err)
	}
	return ctx.GetStub().PutState(tradeID, orderBytes)
}
func (t *CBPMChaincode) supplyOrderExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"SupplyOrder\",\"tradeID\":\"%s\"}}", tradeID)
	queryResults, err := t.getOrderQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return false, fmt.Errorf("fail to check whether order for trade %s exists: %v", tradeID, err)
	}
	if len(queryResults) == 0 {
		return false, nil
	}
	return true, nil
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
		return fmt.Errorf("fail to get peer's orgID: %v", err)
	}

	if clientOrgID != peerOrgID {
		return fmt.Errorf("client from org %s is not authorized to read or write private data from an org %s peer",
			clientOrgID,
			peerOrgID,
		)
	}
	return nil
}

func (s *CBPMChaincode) getAssetQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyAsset, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var results []*SupplyAsset
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		newAsset := new(SupplyAsset)
		err = json.Unmarshal(response.Value, newAsset)
		if err != nil {
			return nil, err
		}
		results = append(results, newAsset)
	}
	return results, nil
}

func (s *CBPMChaincode) getOrderQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyOrder, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var results []*SupplyOrder
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		newOrder := new(SupplyOrder)
		err = json.Unmarshal(response.Value, newOrder)
		if err != nil {
			return nil, err
		}
		results = append(results, newOrder)
	}
	return results, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(&CBPMChaincode{})
	if err != nil {
		log.Panicf("Error creating mischaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting mischaincode: %v", err)
	}
}
