- cbpm: orderer, admin
- manufacturer: peer1, peer2, admin
- supplier: peer1, peer2, admin
- carrier: peer1, peer2, admin
- middleman: peer1, peer2, admin



# 启动网络

```shell
docker-compose up -d 
```



# Create a supplier- carrier channel

## 创建并进入通道

启动cli并且进入

```shell
docker exec -it cli /bin/bash
```



用supplier创建scchannel通道

```shell
# 配置supplier-peer1环境 MSPCONFIGPATH设置为admin的
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP


peer channel create -c scchannel -f /tmp/hyperledger/fabric/channel-artifacts/scchannel.tx -o orderer-cbpm:7050 --outputBlock /tmp/hyperledger/fabric/channel-artifacts/scchannel.block --tls --cafile /tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem

#### 或者
peer channel create -c scchannel -f /tmp/hyperledger/fabric/channel-artifacts/scchannel.tx -o orderer-cbpm:7050 --outputBlock /tmp/hyperledger/fabric/channel-artifacts/scchannel.block --tls --cafile /tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem
```

Supplier节点加入通道

```shell
# cli
export CORE_PEER_ADDRESS=peer1-supplier:7051
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/scchannel.block
export CORE_PEER_ADDRESS=peer2-supplier:7051
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/scchannel.block
```

Carrier节点加入通道

```shell
# cli
# 配置carrier-peer1环境 MSPCONFIGPATH设置为admin的
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=CarrierMSP

export CORE_PEER_ADDRESS=peer1-carrier:7051
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/scchannel.block
export CORE_PEER_ADDRESS=peer2-carrier:7051
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/scchannel.block
```

检测（如果没有的话会报错

```shell
peer channel list
```

## 安装测试链码

TODO：

1. 将最简单的java链码部署成功，并且成功调用
2. 通道间流程定义好后，先写一个通道的智能合约

### 选择链码

链码路径-链码语言

- chaincode-supplier-carrier: golang

- test-go : golang
- test-go-asset-transfer : golang
- test-go-private-data : golang
- test-go-secured-transfer: golang ！！！不能使用，版本过高，该链码使用了` implicit private data namespace reserved for organization-specific private data.` 

### 安装链码

```shell
# 输入想要安装的链码路径以及语言版本
export CHAINCODE=chaincode-supplier-carrier
export CHAINCODE_LANG=golang
export CHAINCODE_VERSION=1.0
export CHAINCODE_NAME=scchaincode

export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_ADDRESS=peer1-supplier:7051
export CORE_PEER_LOCALMSPID=SupplierMSP

peer chaincode install -n $CHAINCODE_NAME -v $CHAINCODE_VERSION  -p $CHAINCODE -l $CHAINCODE_LANG

export CORE_PEER_ADDRESS=peer2-supplier:7051
peer chaincode install -n $CHAINCODE_NAME -v $CHAINCODE_VERSION  -p $CHAINCODE -l $CHAINCODE_LANG


export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_ADDRESS=peer1-carrier:7051
export CORE_PEER_LOCALMSPID=CarrierMSP
peer chaincode install -n $CHAINCODE_NAME -v $CHAINCODE_VERSION  -p $CHAINCODE -l $CHAINCODE_LANG


export CORE_PEER_ADDRESS=peer2-carrier:7051
peer chaincode install -n $CHAINCODE_NAME -v $CHAINCODE_VERSION  -p $CHAINCODE -l $CHAINCODE_LANG

```

### 初始化并调用链码

如果报错`Error: could not assemble transaction, err proposal response was not successful, error code 500, msg failed to execute transaction 8e95ab096404723808c9b0d2350db55dae5d81964f3b02176d10ac5639316e12: error sending: timeout expired while executing transaction`

对于gradle的项目，可以在宿主机执行`gradle build` 

对于go的项目，可以先在宿主机执行`export GO111MODULE=auto && go mod vendor` 



如果invoke成功执行了，但是查询发现没有数据增加到数据库中，检查-P参数的设置





如果一直报错`Error: endorsement failure during query. response: status:500 message:"Received unknown function invocation" ` 

请删除本地docker中的名为`dev-peer*` 的所有镜像，因为可能如果scchannel和scchaincode名字相同的话，可能会用到之前的智能合约的镜像...导致函数调用不成功



#### chaincode-supplier-carrier

```shell
peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"


peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -c '{"Args":["createDeliveryOrder", "trade1","asset1","100","placeA","placeB"]}'

```



#### test-go-private-data

```shell
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')" --collections-config $GOPATH/src/$CHAINCODE/collections_config.json

export MARBLE=$(echo -n "{\"name\":\"marble1\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)


peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -c '{"Args":["InitMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"

export MARBLE=$(echo -n "{\"name\":\"marble2\",\"color\":\"red\",\"size\":50,\"owner\":\"tom\",\"price\":102}" | base64 | tr -d \\n)

peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -c '{"Args":["InitMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"


peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["readMarble","marble1"]}'


peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["readMarblePrivateDetails","marble1"]}'



# switch to carrier
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_LOCALMSPID=CarrierMSP
export CORE_PEER_ADDRESS=peer1-carrier:7051

peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["readMarble","marble1"]}'


peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["readMarblePrivateDetails","marble1"]}'

peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["GetMarbleHash","collectionMarbles","marble1"]}'

peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["QueryMarblesByOwner","tom"]}'

peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["QueryMarbles","{\"selector\":{\"owner\":\"tom\"}}"]}'

```

#### go-marbles-couchdb

```shell
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051


peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

```

调用查询

```shell
peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -c '{"Args":["initMarble","marble1","blue","35","tom"]}'

peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -c '{"function": "initMarble","Args":["marble1","blue","35","tom"]}'


peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"function": "queryMarbles","Args":["{\"selector\":{\"owner\":\"tom\"}}"]}'


peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["readMarble","marble1"]}'
```







#### test-go

```shell
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":["initLedger"]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

```



```shell

peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["ReadAsset","asset1"]}'


peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -c '{"Args":["CreateAsset","asset10","black","5","cynyard","10"]}'


peer chaincode query -C scchannel -n $CHAINCODE_NAME -c '{"Args":["ReadAsset","asset10"]}'


```







#### go-secured-transfer-couchdb-不能使用-版本不支持

```shell
# 用peer1-supplier初始化（其实都无所谓）
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

peer chaincode upgrade -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"


```

调用

```shell
# supplier中 创建一个asset，以下是私密数据（但是没有使用private data collection，在智能合约加密）公开数据只有assetID和publicDescription
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

export ASSET_PROPERTIES=$(echo -n "{\"object_type\":\"asset_properties\",\"asset_id\":\"asset1\",\"color\":\"blue\",\"size\":35,\"salt\":\"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3\"}" | base64 | tr -d \\n)


peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"Args":["CreateAsset", "asset1", "A new asset for Org1MSP"]}' --transient "{\"asset_properties\":\"$ASSET_PROPERTIES\"}"


peer chaincode query -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"function":"GetAssetPrivateProperties","Args":["asset1"]}'

```







## 更新链码

`peer chaincode install` 的命令，修改-v的参数



`peer chaincode instantiate` 的命令，将该命令改为`peer chaincode upgrade` ，并且修改-v的参数