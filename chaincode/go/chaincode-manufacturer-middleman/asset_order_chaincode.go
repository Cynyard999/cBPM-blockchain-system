package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"time"
)

type CBPMChaincode struct {
	contractapi.Contract
}

type Asset struct {
	ObjectType        string  `json:"docType"`
	AssetID           string  `json:"assetID"`
	AssetName         string  `json:"assetName"`
	AssetPrice        float32 `json:"assetPrice"`
	OwnerOrg          string  `json:"ownerOrg"` // 链码中生成，就是中间商的org
	PublicDescription string  `json:"publicDescription"`
}

type Order struct {
	ObjectType string  `json:"docType"` // 固定
	TradeID    string  `json:"tradeID"`
	OrderID    string  `json:"orderID"` // 自动生成
	AssetID    string  `json:"assetID"`
	AssetName  string  `json:"assetName"`
	AssetPrice float32 `json:"assetPrice"`
	Quantity   int     `json:"quantity"`
	TotalPrice int     `json:"totalPrice"` // 自动生成
	Address    string  `json:"address"`
	Status     int     `json:"status"`     // 0: 未开始 1: 中间商开始处理 2: 中间商确认已完成 3: 生产商确认已完成
	CreateTime string  `json:"createTime"` // 自动生成
	UpdateTime string  `json:"updateTime"` // 自动生成
	OwnerOrg   string  `json:"ownerOrg"`   // 限制修改的权限
	HandlerOrg string  `json:"handlerOrg"` // 限制修改的权限
	Note       string  `json:"note"`
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Order    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// PaginatedQueryResult structure used for returning paginated query results and metadata
type PaginatedQueryResult struct {
	Records             []*Order `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
}

// TODO 查询另一通道的信息，与供货商提供的Asset信息校验
func (t *CBPMChaincode) CreateAsset(ctx contractapi.TransactionContextInterface) error {
	transMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("Error getting transient: " + err.Error())
	}
	transientAssetJSON, ok := transMap["asset"]
	if !ok {
		return fmt.Errorf("marble not found in the transient map")
	}
	type assetTransientInput struct {
		AssetID           string  `json:"assetID"`
		AssetName         string  `json:"assetName"`
		AssetPrice        float32 `json:"assetPrice"`
		PublicDescription string  `json:"publicDescription"`
	}
	var assetInput assetTransientInput
	err = json.Unmarshal(transientAssetJSON, &assetInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
	}
	// check input
	if len(assetInput.AssetID) == 0 {
		return fmt.Errorf("asset ID must be a non-empty string")
	}
	if len(assetInput.AssetName) == 0 {
		return fmt.Errorf("asset name must be a non-empty string")
	}
	if assetInput.AssetPrice <= 0 {
		return fmt.Errorf("asset price field must be a positive number")
	}
	exists, err := t.AssetExists(ctx, assetInput.AssetID)
	if err != nil {
		return fmt.Errorf("failed to create Asset: %v", err)
	}
	if exists {
		return fmt.Errorf("failed to create Asset: asset already exists")
	}

	clientOrgID, err := getClientOrgID(ctx, false)
	if err != nil {
		return fmt.Errorf("failed to get verified OrgID: %v", err)
	}

	// create asset
	asset := &Asset{
		ObjectType:        "Asset",
		AssetID:           assetInput.AssetID,
		AssetName:         assetInput.AssetName,
		AssetPrice:        assetInput.AssetPrice,
		OwnerOrg:          clientOrgID,
		PublicDescription: assetInput.PublicDescription,
	}
	assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// === Save marble to state ===
	err = ctx.GetStub().PutState(asset.AssetID, assetJSONasBytes)
	if err != nil {
		return fmt.Errorf("failed to put Asset: %s", err.Error())
	}
	return nil
}

func (t *CBPMChaincode) AssetExists(ctx contractapi.TransactionContextInterface, assetID string) (bool, error) {
	assetBytes, err := ctx.GetStub().GetState(assetID)
	if err != nil {
		return false, fmt.Errorf("failed to check whether asset %s exists: %v", assetID, err)
	}
	return assetBytes != nil, nil
}

func (t *CBPMChaincode) UpdateAsset(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) DeleteAsset(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) GetAsset(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) GetAllAssets(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) QueryAssets(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) PlaceOrder(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) GetOrder(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) GetAllOrders(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) DeleteOrder(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) HandleOrder(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) FinishOrder(ctx contractapi.TransactionContextInterface) {

}

func (t *CBPMChaincode) ConfirmFinishOrder(ctx contractapi.TransactionContextInterface) {

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
	chaincode, err := contractapi.NewChaincode(&CBPMChaincode{})
	if err != nil {
		log.Panicf("Error creating mmchaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting mmchaincode: %v", err)
	}
}
