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





### go-marbles-couchdb

```shell
# supplier-peer1环境 MSPCONFIGPATH设置为admin的

export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_ADDRESS=peer1-supplier:7051
export CORE_PEER_LOCALMSPID=SupplierMSP

peer chaincode install -n scchaincode -v 1.0  -p test-go/go

export CORE_PEER_ADDRESS=peer2-supplier:7051
peer chaincode install -n scchaincode -v 1.0  -p test-go/go


export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_ADDRESS=peer1-carrier:7051
export CORE_PEER_LOCALMSPID=CarrierMSP
peer chaincode install -n scchaincode -v 1.0  -p test-go/go


export CORE_PEER_ADDRESS=peer2-carrier:7051
peer chaincode install -n scchaincode -v 1.0  -p test-go/go



export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051


peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -l golang -v 1.0 -c '{"Args":["init"]}' -P "OR('Supplier.peer','Carrier.peer')"




```

调用查询

```shell
peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"Args":["initMarble","marble1","blue","35","tom"]}'

peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C scchannel -n scchaincode -c '{"function": "initMarble","Args":["marble1","blue","35","tom"]}'

# 不知道为啥，上面的没有写入数据库，导致下面的命令读不到

peer chaincode query -C scchannel -n scchaincode -c '{"function": "queryMarbles","Args":["{\"selector\":{\"owner\":\"tom\"}}"]}'


peer chaincode query -C scchannel -n scchaincode -c '{"Args":["readMarble","marble1"]}'
```



### private-data-java-couchdb

!!! 没有尝试

```shell
# supplier-peer1环境 MSPCONFIGPATH设置为admin的

export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_ADDRESS=peer1-supplier:7051
export CORE_PEER_LOCALMSPID=SupplierMSP

peer chaincode install -n scchaincode -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/java/test-java-private-data

export CORE_PEER_ADDRESS=peer2-supplier:7051
peer chaincode install -n scchaincode -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/java/test-java-private-data


export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_ADDRESS=peer1-carrier:7051
export CORE_PEER_LOCALMSPID=CarrierMSP
peer chaincode install -n scchaincode -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/java/test-java-private-data


export CORE_PEER_ADDRESS=peer2-carrier:7051

peer chaincode install -n scchaincode -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/java/test-java-private-data




```

初始化链码

```shell
# cli
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" --collections-config /tmp/hyperledger/fabric/chaincode/java/test-java-private-data/collections_config.json -C scchannel -n scchaincode -l java -v 1.0 -c '{"Args":[]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

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







如果报错`error starting container: Failed to generate platform-specific docker build: Error executing build: API error (404): network xxxx not found ""` 

新开terminal 输入`docker network ls` 然后根据网络名修改`docker-compose.yml` 中的`CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE` 环境变量



如果报错`Error: could not assemble transaction, err proposal response was not successful, error code 500, msg failed to execute transaction 8e95ab096404723808c9b0d2350db55dae5d81964f3b02176d10ac5639316e12: error sending: timeout expired while executing transaction`

多试几次？

