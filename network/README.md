# Network Nodes

- ca-tls
- cbpm: orderer-cbpm, cbpm-ca
- manufacturer: peer1, peer2, manufacturer-ca,  couchdb1
- supplier: peer1, peer2, supplier-ca, couchdb2
- carrier: peer1, peer2, carrier-ca, couchdb3
- middleman: peer1, peer2, middleman-ca, couchdb4
- cli



# 启动cBPM网络

```shell
# 清除原先的网络
./clean.sh
# 部署网络 1. 创建fabric网络 2. 创建channel 4. 节点加入channel并广播身份认证leader 4. 部署链码
./start.sh
```



# 更新网络

```shell
### 更新链码
docker exec -it cli /bin/bash

# 在容器中
export CHANNEL=cbpmchannel
export CHAINCODE=chaincode-cbpm
export CHAINCODE_LANG=golang
export CHAINCODE_VERSION=$SPECIFIED_VERSION
export CHAINCODE_NAME=cbpmchaincode

./scripts/upgrade-chaincode.sh

```

