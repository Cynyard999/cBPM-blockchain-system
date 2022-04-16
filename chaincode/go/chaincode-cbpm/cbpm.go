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

// =========== Manufacturer-Middleman-Private-Data ==================

type Asset struct {
	ObjectType        string  `json:"objectType"`
	AssetID           string  `json:"assetID"`
	AssetName         string  `json:"assetName"`
	AssetPrice        float32 `json:"assetPrice"`
	ShippingAddress   string  `json:"shippingAddress"`
	OwnerOrg          string  `json:"ownerOrg"` // 链码中生成，就是中间商的org
	PublicDescription string  `json:"publicDescription"`
}

type Order struct {
	ObjectType       string  `json:"objectType"`
	TradeID          string  `json:"tradeID"`
	AssetID          string  `json:"assetID"`
	AssetName        string  `json:"assetName"`
	AssetPrice       float32 `json:"assetPrice"`
	ShippingAddress  string  `json:"shippingAddress"`
	Quantity         int     `json:"quantity"`
	TotalPrice       float32 `json:"totalPrice"` // 自动生成
	ReceivingAddress string  `json:"receivingAddress"`
	Status           int     `json:"status"`     // 0: 未开始 1: 中间商开始处理 2: 中间商确认已完成 3: 生产商确认已完成
	CreateTime       string  `json:"createTime"` // 自动生成
	UpdateTime       string  `json:"updateTime"` // 自动生成
	OwnerOrg         string  `json:"ownerOrg"`   // 限制修改的权限
	HandlerOrg       string  `json:"handlerOrg"` // 限制修改的权限
	Note             string  `json:"note"`
}

// CreateAsset 创建Asset，transient传入assetID，assetName，assetPrice，shippingAddress，publicDescription，中间商根据供应商的asset创建自己能提供的asset，价格名称等可以变化
func (t *CBPMChaincode) CreateAsset(ctx contractapi.TransactionContextInterface) (*Asset, error) {
	// TODO 查询另一通道的信息，与供货商提供的Asset信息校验
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientAssetJSON, ok := transMap["asset"]
	if !ok {
		return nil, fmt.Errorf("asset field not found in the transient map")
	}
	type assetTransientInput struct {
		AssetID           string  `json:"assetID"`
		AssetName         string  `json:"assetName"`
		AssetPrice        float32 `json:"assetPrice"`
		ShippingAddress   string  `json:"shippingAddress"`
		PublicDescription string  `json:"publicDescription"`
	}
	var assetInput assetTransientInput
	err = json.Unmarshal(transientAssetJSON, &assetInput)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	// check input
	if len(assetInput.AssetID) == 0 {
		return nil, fmt.Errorf("asset ID must be a non-empty string")
	}
	if len(assetInput.AssetName) == 0 {
		return nil, fmt.Errorf("asset name must be a non-empty string")
	}
	if len(assetInput.AssetName) == 0 {
		return nil, fmt.Errorf("shipping address must be a non-empty string")
	}
	if assetInput.AssetPrice <= 0 {
		return nil, fmt.Errorf("asset price field must be a positive number")
	}
	exists, err := t.assetExists(ctx, assetInput.AssetID)
	if err != nil {
		return nil, fmt.Errorf("fail to create Asset: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("fail to create Asset: asset already exists")
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("fail to get verified OrgID: %v", err)
	}

	// create asset
	asset := &Asset{
		ObjectType:        "Asset",
		AssetID:           assetInput.AssetID,
		AssetName:         assetInput.AssetName,
		AssetPrice:        assetInput.AssetPrice,
		ShippingAddress:   assetInput.ShippingAddress,
		OwnerOrg:          clientOrgID,
		PublicDescription: assetInput.PublicDescription,
	}
	assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	// === Save marble to state ===
	err = ctx.GetStub().PutPrivateData("AssetCollection", asset.AssetID, assetJSONasBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create Asset: %s", err.Error())
	}
	return asset, nil
}

// UpdateAsset 更新Asset，args传入assetID,assetName,desc,assetPrice，只能owner更新
func (t *CBPMChaincode) UpdateAsset(ctx contractapi.TransactionContextInterface, assetID string, assetName string, desc string, assetPrice float32) error {
	// TODO 重名检测
	asset, err := t.GetAsset(ctx, assetID)
	if err != nil {
		return err
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return err
	}
	if asset.OwnerOrg != clientOrgID {
		return fmt.Errorf("fail to create Asset: unauthorized updater %s", clientOrgID)
	}
	asset.AssetPrice = assetPrice
	asset.AssetName = assetName
	asset.PublicDescription = desc
	newAssetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutPrivateData("AssetCollection", assetID, newAssetBytes)
}

// DeleteAsset 删除Asset
func (t *CBPMChaincode) DeleteAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	// TODO owner才能删除
	exist, err := t.assetExists(ctx, assetID)
	if !exist {
		return fmt.Errorf("fail to delete asset: asset does not exist")
	}
	if err != nil {
		return fmt.Errorf("fail to delete asset: %v", err)
	}
	return ctx.GetStub().DelPrivateData("AssetCollection", assetID)
}

// GetAsset 获取Asset，args传入assetID
func (t *CBPMChaincode) GetAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetPrivateData("AssetCollection", assetID)
	if err != nil {
		return nil, fmt.Errorf("fail to get asset %s: %v", assetID, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("fail to get asset %s: asset does not exist", assetBytes)
	}
	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// GetAllAssets 获取所有Assets
func (t *CBPMChaincode) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {

	queryString := "{\"selector\":{\"objectType\":\"Asset\"}}"

	queryResults, err := t.getAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// QueryAssets 查询满足条件的assets，args传入查询语句
func (t *CBPMChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	queryResults, err := t.getAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

func (t *CBPMChaincode) assetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetPrivateData("AssetCollection", assetID)
	if err != nil {
		return false, fmt.Errorf("fail to read asset %s from world state. %v", assetID, err)
	}
	return assetBytes != nil, nil
}

// CreateOrder 创建Order，transient传入assetID，quantity，receivingAddress，note，链码生成UUID，返回创建好的Order
func (t *CBPMChaincode) CreateOrder(ctx contractapi.TransactionContextInterface) (*Order, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientOrderJSON, ok := transMap["order"]
	if !ok {
		return nil, fmt.Errorf("order not found in the transient map")
	}
	type orderTransientInput struct {
		AssetID          string `json:"assetID"`
		Quantity         int    `json:"quantity"`
		ReceivingAddress string `json:"receivingAddress"`
		Note             string `json:"note"`
	}
	var orderInput orderTransientInput
	err = json.Unmarshal(transientOrderJSON, &orderInput)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	// check input
	if len(orderInput.AssetID) == 0 {
		return nil, fmt.Errorf("asset ID must be a non-empty string")
	}
	if len(orderInput.ReceivingAddress) == 0 {
		return nil, fmt.Errorf("order address must be a non-empty string")
	}
	if orderInput.Quantity <= 0 {
		return nil, fmt.Errorf("asset quantity field must be a positive number")
	}
	asset, err := t.GetAsset(ctx, orderInput.AssetID)
	if err != nil {
		return nil, err
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("fail to get verified OrgID: %v", err)
	}
	newTradeID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("fail to generate Trade ID: %v", err)
	}

	// create order
	order := &Order{
		ObjectType:       "Order",
		TradeID:          newTradeID.String(),
		AssetID:          orderInput.AssetID,
		AssetName:        asset.AssetName,
		AssetPrice:       asset.AssetPrice,
		ShippingAddress:  asset.ShippingAddress,
		Quantity:         orderInput.Quantity,
		TotalPrice:       asset.AssetPrice * (float32(orderInput.Quantity)),
		ReceivingAddress: orderInput.ReceivingAddress,
		Status:           0,
		CreateTime:       time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime:       time.Now().Format("2006-01-02 15:04:05"),
		HandlerOrg:       "",
		OwnerOrg:         clientOrgID,
		Note:             orderInput.Note,
	}
	orderJSONasBytes, err := json.Marshal(order)
	err = ctx.GetStub().PutPrivateData("OrderCollection", order.TradeID, orderJSONasBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create Order: %s", err.Error())
	}
	return order, nil
}

// GetOrder 获取Order，args传入tradeID
func (t *CBPMChaincode) GetOrder(ctx contractapi.TransactionContextInterface, tradeID string) (*Order, error) {
	orderBytes, err := ctx.GetStub().GetPrivateData("OrderCollection", tradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to get order for trade %s: %v", tradeID, err)
	}
	if orderBytes == nil {
		return nil, fmt.Errorf("fail to get order for trade %s: asset does not exist", orderBytes)
	}
	var order Order
	err = json.Unmarshal(orderBytes, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetAllOrders 获取所有Orders
func (t *CBPMChaincode) GetAllOrders(ctx contractapi.TransactionContextInterface) ([]*Order, error) {
	queryString := "{\"selector\":{\"objectType\":\"Order\"}}"

	queryResults, err := t.getOrderQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// QueryOrders 查询满足条件的Orders，args传入查询语句
func (t *CBPMChaincode) QueryOrders(ctx contractapi.TransactionContextInterface, queryString string) ([]*Order, error) {
	queryResults, err := t.getOrderQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// DeleteOrder 删除Order，args传入tradeID
func (t *CBPMChaincode) DeleteOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	// TODO owner才能删除，非进行中的Order才能删除
	exists, err := t.orderExists(ctx, tradeID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("fail to delete order: order for trade #{tradeID} does not exist")
	}
	return ctx.GetStub().DelPrivateData("OrderCollection", tradeID)
}

// HandleOrder 非Owner开始处理Order，args传入tradeID
func (t *CBPMChaincode) HandleOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	order, err := t.GetOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to handle order: %v", err)
	}
	if order.Status != 0 {
		return fmt.Errorf("fail to handle order: order(status: %d) for trade #{tradeID} has been handled", order.Status)
	}
	if order.HandlerOrg != "" {
		return fmt.Errorf("fail to handle order: order for trade #{tradeID} has been handled by some org")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to handle order: %v", err)
	}
	if order.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to handle order: cannot handle as owner")
	}
	order.HandlerOrg = clientOrgID
	order.Status = 1
	order.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("fail to handle order: %v", err)
	}
	return ctx.GetStub().PutPrivateData("OrderCollection", tradeID, orderBytes)
}

// FinishOrder handler完成Order，args传入tradeID
func (t *CBPMChaincode) FinishOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	order, err := t.GetOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to finish order: %v", err)
	}
	if order.Status == 0 {
		return fmt.Errorf("fail to finish order: order for trade #{tradeID} has not been handled")
	}
	if order.Status == 2 {
		return fmt.Errorf("fail to finish order: order for trade #{tradeID} has been finished")
	}
	if order.HandlerOrg == "" {
		return fmt.Errorf("fail to finish order: no handler is specified")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to finish order: %v", err)
	}
	if order.OwnerOrg == clientOrgID {
		return fmt.Errorf("fail to finish order: cannot finish as owner")
	}
	if order.HandlerOrg != clientOrgID {
		return fmt.Errorf("fail to finish order: cannot finish by other org: %s", clientOrgID)
	}
	order.Status = 2
	order.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("fail to finish order: %v", err)
	}
	return ctx.GetStub().PutPrivateData("OrderCollection", tradeID, orderBytes)
}

// ConfirmFinishOrder owner确定完成，args传入tradeID
func (t *CBPMChaincode) ConfirmFinishOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	order, err := t.GetOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to confirm finish order: %v", err)
	}
	if order.Status == 0 {
		return fmt.Errorf("fail to confirm finish order: order has not been handled")
	}
	if order.Status == 1 {
		return fmt.Errorf("fail to confirm finish order: order has not been finished")
	}
	if order.Status == 3 {
		return fmt.Errorf("fail to confirm finish order: order has been confirmed finished")
	}
	if order.HandlerOrg == "" {
		return fmt.Errorf("fail to confirm finish order: no handler is specified")
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to confirm finish order: %v", err)
	}
	if order.OwnerOrg != clientOrgID {
		return fmt.Errorf("fail to confirm finish order: only owner can comfirm finish order")
	}
	order.Status = 3
	order.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("fail to confirm finish order: %v", err)
	}
	return ctx.GetStub().PutPrivateData("OrderCollection", tradeID, orderBytes)
}

func (t *CBPMChaincode) orderExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	orderBytes, err := ctx.GetStub().GetPrivateData("OrderCollection", tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read order for trade %s from world state. %v", tradeID, err)
	}
	return orderBytes != nil, nil
}
func (t *CBPMChaincode) getAssetQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("AssetCollection", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var results []*Asset
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		newAsset := new(Asset)
		err = json.Unmarshal(response.Value, newAsset)
		if err != nil {
			return nil, err
		}
		results = append(results, newAsset)
	}
	return results, nil
}

func (t *CBPMChaincode) getOrderQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Order, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("OrderCollection", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	var results []*Order
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		newOrder := new(Order)
		err = json.Unmarshal(response.Value, newOrder)
		if err != nil {
			return nil, err
		}
		results = append(results, newOrder)
	}
	return results, nil
}

// =========== Middleman-Supplier-Private-Data ==================

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

// CreateSupplyAsset 创建asset，需要transient传入assetName，assetPrice，shippingAddress，publicDescription 会生成一个UUID表示asset的名称，最终返回创建的asset对象
func (t *CBPMChaincode) CreateSupplyAsset(ctx contractapi.TransactionContextInterface) (*SupplyAsset, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientAssetJSON, ok := transMap["asset"]
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

	exists, err := t.supplyAssetExists(ctx, assetInput.AssetName)
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

	err = ctx.GetStub().PutPrivateData("SupplyAssetCollection", supplyAsset.AssetID, assetJSONasBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create supplyAsset: %v", err)
	}
	return supplyAsset, nil
}

// UpdateSupplyAsset 更新指定的asset，需要args传入assetID，assetName，assetPrice，shippingAddress，desc，只有owner能够更新
func (t *CBPMChaincode) UpdateSupplyAsset(ctx contractapi.TransactionContextInterface, assetID string, assetName string, assetPrice float32, shippingAddress string, desc string) error {
	asset, err := t.GetSupplyAsset(ctx, assetID)
	if err != nil {
		return fmt.Errorf("fail to update SupplyAsset: %v", err)
	}
	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("fail to update SupplyAsset: %v", err)
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
		return fmt.Errorf("fail to update SupplyAsset: %v", err)
	}
	return ctx.GetStub().PutPrivateData("SupplyAssetCollection", assetID, newAssetBytes)
}

// DeleteSupplyAsset 删除指定的asset，需要args传入assetID
func (t *CBPMChaincode) DeleteSupplyAsset(ctx contractapi.TransactionContextInterface, assetID string) error {
	// TODO owner删除
	exist, err := t.supplyAssetExists(ctx, assetID)
	if !exist {
		return fmt.Errorf("fail to delete asset: asset does not exist")
	}
	if err != nil {
		return fmt.Errorf("fail to delete asset: %v", err)
	}
	return ctx.GetStub().DelPrivateData("SupplyAssetCollection", assetID)
}

// GetSupplyAsset 获取指定的asset，需要args传入assetID
func (t *CBPMChaincode) GetSupplyAsset(ctx contractapi.TransactionContextInterface, assetID string) (*SupplyAsset, error) {
	supplyAssetBytes, err := ctx.GetStub().GetPrivateData("SupplyAssetCollection", assetID)
	if err != nil {
		return nil, fmt.Errorf("fail to get supply asset %s: %v", assetID, err)
	}
	if supplyAssetBytes == nil {
		return nil, fmt.Errorf("fail to get supply asset %s: asset does not exist", assetID)
	}
	var supplyAsset SupplyAsset
	err = json.Unmarshal(supplyAssetBytes, &supplyAsset)
	if err != nil {
		return nil, err
	}
	return &supplyAsset, nil
}

// GetAllSupplyAssets 获取所有的asset
func (t *CBPMChaincode) GetAllSupplyAssets(ctx contractapi.TransactionContextInterface) ([]*SupplyAsset, error) {

	queryString := "{\"selector\":{\"objectType\":\"SupplyAsset\"}}"

	queryResults, err := t.getSupplyAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// QuerySupplyAssets 查询assets，需要args传入查询语句
func (t *CBPMChaincode) QuerySupplyAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyAsset, error) {
	queryResults, err := t.getSupplyAssetQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

func (t *CBPMChaincode) supplyAssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetPrivateData("SupplyAssetCollection", assetID)
	if err != nil {
		return false, fmt.Errorf("fail to read supplyAsset %s from world state. %v", assetID, err)
	}
	return assetBytes != nil, nil
}

//func (t *CBPMChaincode) supplyAssetNameExists(ctx contractapi.TransactionContextInterface, assetName string) (bool, error) {
//	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"SupplyAsset\",\"assetName\":\"%s\"}}", assetName)
//	queryResults, err := t.getSupplyAssetQueryResultForQueryString(ctx, queryString)
//	if err != nil {
//		return false, fmt.Errorf("fail to check whether asset name exists: %v", err)
//	}
//	if len(queryResults) == 0 {
//		return false, nil
//	}
//	return true, nil
//}

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
	asset, err := t.GetAsset(ctx, orderInput.AssetID)
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
	err = ctx.GetStub().PutPrivateData("SupplyOrderCollection", order.TradeID, orderJSONasBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create supply order: %s", err.Error())
	}
	return order, nil
}

// GetSupplyOrder 获取supplyOrder，需要args传入tradeID
func (t *CBPMChaincode) GetSupplyOrder(ctx contractapi.TransactionContextInterface, tradeID string) (*SupplyOrder, error) {
	supplyOrderBytes, err := ctx.GetStub().GetPrivateData("SupplyOrderCollection", tradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to get supply order for trade %s: %v", tradeID, err)
	}
	if supplyOrderBytes == nil {
		return nil, fmt.Errorf("fail to get supply order for trade %s: order does not exist", supplyOrderBytes)
	}
	var supplyOrder SupplyOrder
	err = json.Unmarshal(supplyOrderBytes, &supplyOrder)
	if err != nil {
		return nil, err
	}
	return &supplyOrder, nil
}

// GetAllSupplyOrders 获取所有supplyOrders
func (t *CBPMChaincode) GetAllSupplyOrders(ctx contractapi.TransactionContextInterface) ([]*SupplyOrder, error) {
	queryString := "{\"selector\":{\"objectType\":\"SupplyOrder\"}}"

	queryResults, err := t.getSupplyOrderQueryResultForQueryString(ctx, queryString)
	if err != nil {
		return nil, err
	}
	return queryResults, nil
}

// QuerySupplyOrders 查询满足条件的supplyOrders，需要args传入查询条件
func (t *CBPMChaincode) QuerySupplyOrders(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyOrder, error) {
	queryResults, err := t.getSupplyOrderQueryResultForQueryString(ctx, queryString)
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
	return ctx.GetStub().DelPrivateData("SupplyOrderCollection", tradeID)
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
	return ctx.GetStub().PutPrivateData("SupplyOrderCollection", tradeID, orderBytes)
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
	return ctx.GetStub().PutPrivateData("SupplyOrderCollection", tradeID, orderBytes)
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
	return ctx.GetStub().PutPrivateData("SupplyOrderCollection", tradeID, orderBytes)
}
func (t *CBPMChaincode) supplyOrderExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	supplyOrderBytes, err := ctx.GetStub().GetPrivateData("SupplyOrderCollection", tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read supplyOrder for trade %s from world state. %v", tradeID, err)
	}
	return supplyOrderBytes != nil, nil
}

func (t *CBPMChaincode) getSupplyAssetQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyAsset, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("SupplyAssetCollection", queryString)
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

func (t *CBPMChaincode) getSupplyOrderQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*SupplyOrder, error) {

	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("SupplyOrderCollection", queryString)
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

// =========== Middleman-Carrier-Private-Data ==================

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
		Note:       arrangementInput.Note,
	}

	deliveryArrangementBytes, err := json.Marshal(deliveryArrangement)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery arrangement: %v", err)
	}
	err = ctx.GetStub().PutPrivateData("CarryArrangementCollection", arrangementInput.TradeID, deliveryArrangementBytes)
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
	return ctx.GetStub().DelPrivateData("CarryArrangementCollection", tradeID)
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
	return ctx.GetStub().PutPrivateData("CarryArrangementCollection", tradeID, deliveryArrangementBytes)
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
	return ctx.GetStub().PutPrivateData("CarryArrangementCollection", tradeID, deliveryArrangementBytes)
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
	return ctx.GetStub().PutPrivateData("CarryArrangementCollection", tradeID, deliveryArrangementBytes)
}

// GetDeliveryArrangement 获取DeliveryArrangement，args传入tradeID
func (t *CBPMChaincode) GetDeliveryArrangement(ctx contractapi.TransactionContextInterface, tradeID string) (*DeliveryArrangement, error) {
	deliveryArrangementBytes, err := ctx.GetStub().GetPrivateData("CarryArrangementCollection", tradeID)
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
	return getDeliveryArrangementQueryResultForQueryString(ctx, queryString)
}

// QueryDeliveryArrangements 查询满足条件的DeliveryArrangements，args传入查询语句
func (t *CBPMChaincode) QueryDeliveryArrangements(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryArrangement, error) {
	return getDeliveryArrangementQueryResultForQueryString(ctx, queryString)
}
func getDeliveryArrangementQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryArrangement, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("CarryArrangementCollection", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
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

func (t *CBPMChaincode) deliveryArrangementExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	deliveryArrangementBytes, err := ctx.GetStub().GetPrivateData("CarryArrangementCollection", tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read delivery arrangement for trade %s from world state. %v", tradeID, err)
	}
	return deliveryArrangementBytes != nil, nil
}

// =========== Supplier-Carrier-Private-Data ==================

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

// CreateDeliveryOrder 创建DeliveryOrder，transient传入tradeID，assetName，note
func (t *CBPMChaincode) CreateDeliveryOrder(ctx contractapi.TransactionContextInterface) (*DeliveryOrder, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientOrderJSON, ok := transMap["order"]
	if !ok {
		return nil, fmt.Errorf("order not found in the transient map")
	}
	type orderTransientInput struct {
		TradeID   string `json:"tradeID"`
		AssetName string `json:"assetName"`
		Note      string `json:"note"`
	}
	var orderInput orderTransientInput
	err = json.Unmarshal(transientOrderJSON, &orderInput)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	if len(orderInput.TradeID) == 0 {
		return nil, fmt.Errorf("trade ID must be a non-empty string")
	}
	if len(orderInput.AssetName) == 0 {
		return nil, fmt.Errorf("asset name must be a non-empty string")
	}
	exists, err := t.deliveryOrderExists(ctx, orderInput.TradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery order: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("fail to create delivery order: order for trade %s already exists", orderInput.TradeID)
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery order: %v", err)
	}

	deliveryOrder := &DeliveryOrder{
		ObjectType: "DeliveryOrder",
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
		return nil, fmt.Errorf("fail to create delivery order: %v", err)
	}
	err = ctx.GetStub().PutPrivateData("DeliveryOrderCollection", orderInput.TradeID, deliveryOrderBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery order: %v", err)
	}
	return deliveryOrder, nil
}

// DeleteDeliveryOrder 删除DeliveryOrder，args传入tradeID
func (t *CBPMChaincode) DeleteDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	// TODO 只有owner能删除，只有未进行的订单能删除
	deliveryOrder, err := t.GetDeliveryOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to delete delivery order: %v", err)
	}
	if deliveryOrder.Status != 0 {
		return fmt.Errorf("fail to delete delivery order: cannot delete order in progress")
	}
	return ctx.GetStub().DelPrivateData("DeliveryOrderCollection", tradeID)
}

// HandleDeliveryOrder 开始处理DeliveryOrder，args传入tradeID，非owner的组织才能护处理
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
	return ctx.GetStub().PutPrivateData("DeliveryOrderCollection", tradeID, deliveryOrderBytes)
}

// FinishDeliveryOrder handler处理完成订单，args传入tradeID
func (t *CBPMChaincode) FinishDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryOrder, err := t.GetDeliveryOrder(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to finish delivery order: %v", err)
	}
	if deliveryOrder.Status == 0 {
		return fmt.Errorf("fail to finish delivery order: order for trade #{tradeID} has not been handled")
	}
	if deliveryOrder.Status == 2 {
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
	return ctx.GetStub().PutPrivateData("DeliveryOrderCollection", tradeID, deliveryOrderBytes)
}

// GetDeliveryOrder 获取DeliveryOrder，args传入tradeID
func (t *CBPMChaincode) GetDeliveryOrder(ctx contractapi.TransactionContextInterface, tradeID string) (*DeliveryOrder, error) {
	deliveryOrderBytes, err := ctx.GetStub().GetPrivateData("DeliveryOrderCollection", tradeID)
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

// GetAllDeliveryOrders 获取所有DeliveryOrder
func (t *CBPMChaincode) GetAllDeliveryOrders(ctx contractapi.TransactionContextInterface) ([]*DeliveryOrder, error) {
	queryString := "{\"selector\":{\"objectType\":\"DeliveryOrder\"}}"
	return getDeliveryOrderQueryResultForQueryString(ctx, queryString)
}

// QueryDeliveryOrders 查询满足条件的deliveryOrders，args传入查询语句
func (t *CBPMChaincode) QueryDeliveryOrders(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryOrder, error) {
	return getDeliveryOrderQueryResultForQueryString(ctx, queryString)
}

// getDeliveryOrderQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getDeliveryOrderQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryOrder, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("DeliveryOrderCollection", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
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

func (t *CBPMChaincode) deliveryOrderExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	deliveryOrderBytes, err := ctx.GetStub().GetPrivateData("DeliveryOrderCollection", tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read delivery order for trade %s from world state. %v", tradeID, err)
	}
	return deliveryOrderBytes != nil, nil

}

// =========== Carrier-Manufacturer-Private-Data ==================

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

// CreateDeliveryDetail 创建DeliveryDetail，需要传入tradeID，assetName，startPlace，endPlace，contact，note
func (t *CBPMChaincode) CreateDeliveryDetail(ctx contractapi.TransactionContextInterface) (*DeliveryDetail, error) {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientDetailJSON, ok := transMap["detail"]
	if !ok {
		return nil, fmt.Errorf("detail not found in the transient map")
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
		return nil, fmt.Errorf("fail to unmarshal JSON: %s", err.Error())
	}
	if len(detailInput.TradeID) == 0 {
		return nil, fmt.Errorf("trade ID must be a non-empty string")
	}
	if len(detailInput.AssetName) == 0 {
		return nil, fmt.Errorf("asset name must be a non-empty string")
	}
	if len(detailInput.StartPlace) == 0 {
		return nil, fmt.Errorf("start place must be a non-empty string")
	}
	if len(detailInput.EndPlace) == 0 {
		return nil, fmt.Errorf("end place must be a non-empty string")
	}
	if len(detailInput.Contact) == 0 {
		return nil, fmt.Errorf("contact must be a non-empty string")
	}

	exists, err := t.deliveryDetailExists(ctx, detailInput.TradeID)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery detail: %v", err)
	}
	if exists {
		return nil, fmt.Errorf("fail to create delivery detail: detail for trade %s already exists", detailInput.TradeID)
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery detail: %v", err)
	}

	deliveryDetail := &DeliveryDetail{
		ObjectType: "DeliveryDetail",
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
		return nil, err
	}
	err = ctx.GetStub().PutPrivateData("DeliveryDetailCollection", detailInput.TradeID, deliveryDetailBytes)
	if err != nil {
		return nil, fmt.Errorf("fail to create delivery detail: %s", err.Error())
	}
	return deliveryDetail, nil

}

// DeleteDeliveryDetail 删除DeliveryDetail，args传入tradeID
func (t *CBPMChaincode) DeleteDeliveryDetail(ctx contractapi.TransactionContextInterface, tradeID string) error {
	// TODO 只有owner能删除，只能删除没有进行中
	deliveryDetail, err := t.GetDeliveryDetail(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to delete delivery detail: %v", err)
	}
	if deliveryDetail.Status != 0 {
		return fmt.Errorf("fail to delete delivery detail: cannot delete detail in progress")
	}
	return ctx.GetStub().DelPrivateData("DeliveryDetailCollection", tradeID)
}

// HandleDeliveryDetail 开始处理DeliveryDetail，args传入tradeID，只有owner能处理，在业务中，只有carrier能处理
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
	return ctx.GetStub().PutPrivateData("DeliveryDetailCollection", tradeID, deliveryDetailBytes)
}

// FinishDeliveryDetail 处理完成，args传入tradeID，只有owner能完成
func (t *CBPMChaincode) FinishDeliveryDetail(ctx contractapi.TransactionContextInterface, tradeID string) error {
	deliveryDetail, err := t.GetDeliveryDetail(ctx, tradeID)
	if err != nil {
		return fmt.Errorf("fail to finish delivery detail: %v", err)
	}
	if deliveryDetail.Status == 0 {
		return fmt.Errorf("fail to finish delivery detail: detail for trade #{tradeID} has not been handled")
	}
	if deliveryDetail.Status == 2 {
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
	return ctx.GetStub().PutPrivateData("DeliveryDetailCollection", tradeID, deliveryDetailBytes)
}

// GetDeliveryDetail 获取DeliveryDetail信息，args传入tradeID
func (t *CBPMChaincode) GetDeliveryDetail(ctx contractapi.TransactionContextInterface, tradeID string) (*DeliveryDetail, error) {
	deliveryDetailBytes, err := ctx.GetStub().GetPrivateData("DeliveryDetailCollection", tradeID)
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

// GetAllDeliveryDetails 获取所有DeliveryDetails
func (t *CBPMChaincode) GetAllDeliveryDetails(ctx contractapi.TransactionContextInterface) ([]*DeliveryDetail, error) {
	queryString := "{\"selector\":{\"objectType\":\"DeliveryDetail\"}}"
	return getDeliveryDetailQueryResultForQueryString(ctx, queryString)
}

// QueryDeliveryDetails 查询满足条件的DeliveryDetails,args传入查询语句
func (t *CBPMChaincode) QueryDeliveryDetails(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryDetail, error) {
	return getDeliveryDetailQueryResultForQueryString(ctx, queryString)
}

func getDeliveryDetailQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*DeliveryDetail, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataQueryResult("DeliveryDetailCollection", queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
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

func (t *CBPMChaincode) deliveryDetailExists(ctx contractapi.TransactionContextInterface, tradeID string) (bool, error) {
	deliveryDetailBytes, err := ctx.GetStub().GetPrivateData("DeliveryDetailCollection", tradeID)
	if err != nil {
		return false, fmt.Errorf("fail to read delivery detail for trade %s from world state. %v", tradeID, err)
	}
	return deliveryDetailBytes != nil, nil

}

// =========== Chaincode Common Functions ==================

func getClientOrgID(ctx contractapi.TransactionContextInterface, verifyOrg bool) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("fail to get client's orgID: %v", err)
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

func main() {
	chaincode, err := contractapi.NewChaincode(&CBPMChaincode{})
	if err != nil {
		log.Panicf("Error creating cbpmchaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting cbpmchaincode: %v", err)
	}
}
