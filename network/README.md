- manufacturer: orderer, peer1, peer2, admin
- supplier: peer1, peer2, admin
- carrier: peer1, peer2, admin
- middleman: peer1, peer2, admin





# CA配置

ca-tls

```shell
fabric-ca-client register -d --id.name peer1-supplier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-supplier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052

fabric-ca-client register -d --id.name peer1-carrier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-carrier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052


fabric-ca-client register -d --id.name peer1-middleman --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-middleman --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052


fabric-ca-client register -d --id.name peer1-manufacturer --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-manufacturer --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052


fabric-ca-client register -d --id.name orderer-manufacturer --id.secret ordererPW --id.type orderer -u https://0.0.0.0:7052
```





manufacturer

```shell
docker-compose -f docker-compose.yaml up manufacturer

export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/admin


fabric-ca-client enroll -d -u https://manufacturer-admin:manufacturer-adminpw@0.0.0.0:7053


fabric-ca-client register -d --id.name peer1-manufacturer --id.secret peer1PW --id.type peer -u https://0.0.0.0:7053
fabric-ca-client register -d --id.name peer2-manufacturer --id.secret peer2PW --id.type peer -u https://0.0.0.0:7053

fabric-ca-client register -d --id.name orderer-manufacturer --id.secret ordererPW --id.type orderer -u https://0.0.0.0:7053

fabric-ca-client register -d --id.name admin-manufacturer --id.secret adminpw --id.type admin --id.attrs "hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert" -u https://0.0.0.0:7053
```





supplier

```shell
docker-compose -f docker-compose.yaml up supplier



export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/admin


fabric-ca-client enroll -d -u https://supplier-admin:supplier-adminpw@0.0.0.0:7054


fabric-ca-client register -d --id.name peer1-supplier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name peer2-supplier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name admin-supplier --id.secret adminPW --id.type admin -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name user-supplier --id.secret userPW --id.type user -u https://0.0.0.0:7054
```



carrier

```shell
docker-compose -f docker-compose.yaml up carrier



export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/admin


fabric-ca-client enroll -d -u https://carrier-admin:carrier-adminpw@0.0.0.0:7055


fabric-ca-client register -d --id.name peer1-carrier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name peer2-carrier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name admin-carrier --id.secret adminPW --id.type admin -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name user-carrier --id.secret userPW --id.type user -u https://0.0.0.0:7055
```

middleman

```shell
docker-compose -f docker-compose.yaml up middleman



export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/admin


fabric-ca-client enroll -d -u https://middleman-admin:middleman-adminpw@0.0.0.0:7056


fabric-ca-client register -d --id.name peer1-middleman --id.secret peer1PW --id.type peer -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name peer2-middleman --id.secret peer2PW --id.type peer -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name admin-middleman --id.secret adminPW --id.type admin -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name user-middleman --id.secret userPW --id.type user -u https://0.0.0.0:7056
```







## 节点配置

MSP证书以及TLS证书



manufacturer msp

```shell
# 指定本组织的TLS根证书
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/ca-cert.pem
# orderer
# 指定order节点的HOME目录
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/orderer

fabric-ca-client enroll -d -u https://orderer-manufacturer:ordererPW@0.0.0.0:7053

# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer1


fabric-ca-client enroll -d -u https://peer1-manufacturer:peer1PW@0.0.0.0:7053


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer2


fabric-ca-client enroll -d -u https://peer2-manufacturer:peer2PW@0.0.0.0:7053


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/admin
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/admin/msp

fabric-ca-client enroll -d -u https://admin-manufacturer:adminpw@0.0.0.0:7053


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer1/msp/admincerts/supplier-admin-cert.pem

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer2/msp/admincerts/supplier-admin-cert.pem

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/orderer/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/orderer/msp/admincerts/supplier-admin-cert.pem

```



manufacturer tls 

**admin不用配置tls是因为我们生成admin的证书主要就是为了之后链码的安装和实例化，所以配不配置他的TLS证书也无关紧要了(关键是我们之前也没有将这个用户注册到tls服务器中)**

```shell
#指定TLS CA服务器生成的TLS根证书
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/tls/ca-cert.pem

#指定TLS根证书 orderer
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/orderer/tls

fabric-ca-client enroll -d -u https://orderer-manufacturer:ordererPW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts orderer-manufacturer

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/orderer/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA


#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer1/tls

fabric-ca-client enroll -d -u https://peer1-manufacturer:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-manufacturer

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer2/tls

fabric-ca-client enroll -d -u https://peer2-manufacturer:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-manufacturer

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/manufacturer/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA
```

启动peers

```shell
docker-compose -f docker-compose.yaml up peer1-manufacturer
docker-compose -f docker-compose.yaml up peer2-manufacturer
```



supplier msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer1


fabric-ca-client enroll -d -u https://peer1-supplier:peer1PW@0.0.0.0:7054


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer2


fabric-ca-client enroll -d -u https://peer2-supplier:peer2PW@0.0.0.0:7054


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/admin/msp

fabric-ca-client enroll -d -u https://admin-supplier:adminPW@0.0.0.0:7054

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer1/msp/admincerts/supplier-admin-cert.pem


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer2/msp/admincerts/supplier-admin-cert.pem



```

supplier tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer1/tls

fabric-ca-client enroll -d -u https://peer1-supplier:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-supplier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer2/tls

fabric-ca-client enroll -d -u https://peer2-supplier:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-supplier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/supplier/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA

```

启动peers

```shell
docker-compose -f docker-compose.yaml up peer1-supplier
docker-compose -f docker-compose.yaml up peer2-supplier
```



carrier msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer1


fabric-ca-client enroll -d -u https://peer1-carrier:peer1PW@0.0.0.0:7055


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer2


fabric-ca-client enroll -d -u https://peer2-carrier:peer2PW@0.0.0.0:7055


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/admin/msp

fabric-ca-client enroll -d -u https://admin-carrier:adminPW@0.0.0.0:7055

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer1/msp/admincerts/carrier-admin-cert.pem


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer2/msp/admincerts/carrier-admin-cert.pem


```

carrier tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer1/tls

fabric-ca-client enroll -d -u https://peer1-carrier:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-carrier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer2/tls

fabric-ca-client enroll -d -u https://peer2-carrier:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-carrier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/carrier/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA


```

启动peers

```shell
docker-compose -f docker-compose.yaml up peer1-carrier
docker-compose -f docker-compose.yaml up peer2-carrier
```



middleman msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/ca-cert.pem
# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer1


fabric-ca-client enroll -d -u https://peer1-middleman:peer1PW@0.0.0.0:7056


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer2


fabric-ca-client enroll -d -u https://peer2-middleman:peer2PW@0.0.0.0:7056


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/admin/msp

fabric-ca-client enroll -d -u https://admin-middleman:adminPW@0.0.0.0:7056

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer1/msp/admincerts/middleman-admin-cert.pem


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer2/msp/admincerts/middleman-admin-cert.pem



```

middleman tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer1/tls

fabric-ca-client enroll -d -u https://peer1-middleman:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-middleman

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer2/tls

fabric-ca-client enroll -d -u https://peer2-middleman:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-middleman

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA/middleman/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/CA


```

启动peers

```shell
docker-compose -f docker-compose.yaml up peer1-middleman
docker-compose -f docker-compose.yaml up peer2-middleman
```





MSP创世区块配置

```shell

# manufacturer
cd ./manufacturer/msp
mkdir admincerts && mkdir tlscacerts 
cd ..
cp ./admin/msp/signcerts/cert.pem ./msp/admincerts/ca-cert.pem
cp ./ca-cert.pem ./msp/cacerts/ca-cert.pem
cp ../tls/ca-cert.pem ./msp/tlscacerts/ca-cert.pem
cd ..

# supplier
cd ./supplier/msp
mkdir admincerts && mkdir tlscacerts 
cd ..
cp ./admin/msp/signcerts/cert.pem ./msp/admincerts/ca-cert.pem
cp ./ca-cert.pem ./msp/cacerts/ca-cert.pem
cp ../tls/ca-cert.pem ./msp/tlscacerts/ca-cert.pem
cd ..

# carrier
cd ./carrier/msp
mkdir admincerts && mkdir tlscacerts 
cd ..
cp ./admin/msp/signcerts/cert.pem ./msp/admincerts/ca-cert.pem
cp ./ca-cert.pem ./msp/cacerts/ca-cert.pem
cp ../tls/ca-cert.pem ./msp/tlscacerts/ca-cert.pem
cd ..

# middleman
cd ./middleman/msp
mkdir admincerts && mkdir tlscacerts 
cd ..
cp ./admin/msp/signcerts/cert.pem ./msp/admincerts/ca-cert.pem
cp ./ca-cert.pem ./msp/cacerts/ca-cert.pem
cp ../tls/ca-cert.pem ./msp/tlscacerts/ca-cert.pem
cd ..
```

生成创世区块

```shell
mkdir system-genesis-block
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./system-genesis-block/genesis.block -channelID syschannel
```

生成通道信息

```shell
mkdir channel
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel/mychannel.tx -channelID mychannel
```

运行orderer

```
docker-compose -f docker-compose.yaml up orderer-manufacturer
```



配置cli

启动cli并且进入

```shell
docker-compose up -d 
```

supplier-cli

```shell
docker exec -it cli-supplier sh

export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/supplier/admin/msp

peer channel create -o orderer-manufacturer:7050 -c mychannel --ordererTLSHostnameOverride orderer-manufacturer -f /tmp/hyperledger/channel/mychannel.tx --outputBlock /tmp/hyperledger/channel/mychannel.block --tls --cafile /tmp/hyperledger/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem

export CORE_PEER_ADDRESS=peer1-supplier:7051
peer channel join -b /tmp/hyperledger/channel/mychannel.block


export CORE_PEER_ADDRESS=peer2-supplier:7051
peer channel join -b /tmp/hyperledger/channel/mychannel.block



```

carrier-cli

```shell
docker exec -it cli-carrier sh


export CORE_PEER_ADDRESS=peer1-supplier:7051
peer channel join -b /tmp/hyperledger/channel/mychannel.block

export CORE_PEER_ADDRESS=peer2-supplier:7051
peer channel join -b /tmp/hyperledger/channel/mychannel.block
```

