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
docker exec -it cli sh
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

链码描述-链码路径-链码语言

- go-marbles-couchdb : test-go-asset-transfer : golang
- go-asset-transfer-couchdb: test-go-asset-transfer : golang
- java-private-asset-coudhdb: /tmp/hyperledger/fabric/chaincode/java/test-java-private-data : java

### 安装链码

```shell
# 输入想要安装的链码路径以及语言版本
export CHAINCODE=
export CHAINCODE_LANG=


export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_ADDRESS=peer1-supplier:7051
export CORE_PEER_LOCALMSPID=SupplierMSP

peer chaincode install -n scchaincode -v 1.0  -p $CHAINCODE -l $CHAINCODE_LANG

export CORE_PEER_ADDRESS=peer2-supplier:7051
peer chaincode install -n scchaincode -v 1.0  -p $CHAINCODE -l $CHAINCODE_LANG


export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_ADDRESS=peer1-carrier:7051
export CORE_PEER_LOCALMSPID=CarrierMSP
peer chaincode install -n scchaincode -v 1.0  -p $CHAINCODE -l $CHAINCODE_LANG


export CORE_PEER_ADDRESS=peer2-carrier:7051
peer chaincode install -n scchaincode -v 1.0  -p $CHAINCODE -l $CHAINCODE_LANG

```

### 初始化链码

#### go-marbles-couchdb

```shell
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051


peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -l $CHAINCODE_LANG -v 1.0 -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

```

调用查询

```shell
peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"Args":["initMarble","marble1","blue","35","tom"]}'

peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"function": "initMarble","Args":["marble1","blue","35","tom"]}'


peer chaincode query -C scchannel -n scchaincode -c '{"function": "queryMarbles","Args":["{\"selector\":{\"owner\":\"tom\"}}"]}'


peer chaincode query -C scchannel -n scchaincode -c '{"Args":["readMarble","marble1"]}'
```

#### go-asset-transfer-couchdb

```shell
peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -l $CHAINCODE_LANG -v 1.1 -c '{"Args":["initLedger"]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"


peer chaincode query -C scchannel -n scchaincode -c '{"Args":["ReadAsset","asset1"]}'


peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"Args":["CreateAsset","asset10","black","5","cynyard","10"]}'

peer chaincode query -C scchannel -n scchaincode -c '{"Args":["ReadAsset","asset10"]}'
```

#### java-private-asset-coudhdb

！！没有成功

初始化链码

```shell
# cli
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" --collections-config $CHAINCODE/collections_config.json -C scchannel -n scchaincode -l $CHAINCODE_LANG -v 1.0 -c '{"Args":[]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

```

- -P：指定背书策略，上例中的背书策略是，两个组织需要参与链码invoke或query，chaincode执行才能生效。这个参数可为空，则任意安装了链码的节点无约束地调用链码。
- -n：指定链码名称
- -v：指定链码版本号
- -l：指定链码使用的语言，可以是golang, java, nodejs
- -o：指定排序节点
- -C：指定通道名
- -c：指定初始化参数



调用

```shell
export ASSET_PROPERTIES=$(echo -n "{\"objectType\":\"asset\",\"assetID\":\"asset1\",\"color\":\"green\",\"size\":20,\"appraisedValue\":100}" | base64 | tr -d \\n)

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"function":"CreateAsset","Args":[]}' --transient "{\"asset_properties\":\"$ASSET_PROPERTIES\"}"
```



如果报错`Error: could not assemble transaction, err proposal response was not successful, error code 500, msg failed to execute transaction 8e95ab096404723808c9b0d2350db55dae5d81964f3b02176d10ac5639316e12: error sending: timeout expired while executing transaction`

对于gradle的项目，可以在宿主机执行`gradle build` 

对于go的项目，可以先在宿主机执行`export GO111MODULE=auto && go mod vendor` 



如果invoke成功执行了，但是查询发现没有数据增加到数据库中，检查-P参数的设置



## 更新链码

`peer chaincode install` 的命令，修改-v的参数



`peer chaincode instantiate` 的命令，将该命令改为`peer chaincode upgrade` ，并且修改-v的参数