- cbpm: orderer, admin
- manufacturer: peer1, peer2, admin
- supplier: peer1, peer2, admin
- carrier: peer1, peer2, admin
- middleman: peer1, peer2, admin



# 下面的都不用管了

## CA配置

ca服务器配置

```shell
#设置环境变量指定根证书的路径(如果工作目录不同的话记得指定自己的工作目录,以下不再重复说明)
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/tls/ca-cert.pem
#设置环境变量指定CA客户端的HOME文件夹
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/tls/admin
#登录管理员用户用于之后的节点身份注册
fabric-ca-client enroll -d -u https://tls-ca-admin:tls-ca-adminpw@0.0.0.0:7052
```







ca-tls

注册节点身份

```shell
fabric-ca-client register -d --id.name peer1-supplier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-supplier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052

fabric-ca-client register -d --id.name peer1-carrier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-carrier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052


fabric-ca-client register -d --id.name peer1-middleman --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-middleman --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052


fabric-ca-client register -d --id.name peer1-manufacturer --id.secret peer1PW --id.type peer -u https://0.0.0.0:7052
fabric-ca-client register -d --id.name peer2-manufacturer --id.secret peer2PW --id.type peer -u https://0.0.0.0:7052


fabric-ca-client register -d --id.name orderer-cbpm --id.secret ordererPW --id.type orderer -u https://0.0.0.0:7052
```

cbpm-ca

配置了cbpm-ca过后

```shell
docker-compose up cbpm-ca
```

打开另一个

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/admin

# 登陆
fabric-ca-client enroll -d -u https://cbpm-ca-admin:cbpm-adminpw@0.0.0.0:7053

#注册
fabric-ca-client register -d --id.name orderer-cbpm --id.secret ordererPW --id.type orderer -u https://0.0.0.0:7053

fabric-ca-client register -d --id.name admin-cbpm --id.secret adminPW --id.type admin --id.attrs "hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert" -u https://0.0.0.0:7053


```





manufacturer-ca

```shell
docker-compose up manufacturer-ca

export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/admin


fabric-ca-client enroll -d -u https://manufacturer-ca-admin:manufacturer-adminpw@0.0.0.0:7054




fabric-ca-client register -d --id.name peer1-manufacturer --id.secret peer1PW --id.type peer -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name peer2-manufacturer --id.secret peer2PW --id.type peer -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name admin-manufacturer --id.secret adminPW --id.type admin -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name user-manufacturer --id.secret userPW --id.type user -u https://0.0.0.0:7054


```





supplier-ca

```shell
docker-compose up supplier-ca



export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/admin


fabric-ca-client enroll -d -u https://supplier-ca-admin:supplier-adminpw@0.0.0.0:7055


fabric-ca-client register -d --id.name peer1-supplier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name peer2-supplier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name admin-supplier --id.secret adminPW --id.type admin -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name user-supplier --id.secret userPW --id.type user -u https://0.0.0.0:7055
```



carrier-ca

```shell
docker-compose up carrier-ca



export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/networkrrier/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/networkrrier/admin


fabric-ca-client enroll -d -u https://carrier-ca-admin:carrier-adminpw@0.0.0.0:7056


fabric-ca-client register -d --id.name peer1-carrier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name peer2-carrier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name admin-carrier --id.secret adminPW --id.type admin -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name user-carrier --id.secret userPW --id.type user -u https://0.0.0.0:7056
```

middleman

```shell
docker-compose up middleman



export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/admin


fabric-ca-client enroll -d -u https://middleman-ca-admin:middleman-adminpw@0.0.0.0:7057


fabric-ca-client register -d --id.name peer1-middleman --id.secret peer1PW --id.type peer -u https://0.0.0.0:7057
fabric-ca-client register -d --id.name peer2-middleman --id.secret peer2PW --id.type peer -u https://0.0.0.0:7057
fabric-ca-client register -d --id.name admin-middleman --id.secret adminPW --id.type admin -u https://0.0.0.0:7057
fabric-ca-client register -d --id.name user-middleman --id.secret userPW --id.type user -u https://0.0.0.0:7057
```







## 节点配置

MSP证书以及TLS证书



manufacturer msp

```shell
# 指定本组织的TLS根证书
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer1


fabric-ca-client enroll -d -u https://peer1-manufacturer:peer1PW@0.0.0.0:7054


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer2


fabric-ca-client enroll -d -u https://peer2-manufacturer:peer2PW@0.0.0.0:7054


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/admin
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/admin/msp

fabric-ca-client enroll -d -u https://admin-manufacturer:adminPW@0.0.0.0:7054


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer1/msp/admincerts/manufacturer-admin-cert.pem

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer2/msp/admincerts/manufacturer-admin-cert.pem

```



manufacturer tls 

**admin不用配置tls是因为我们生成admin的证书主要就是为了之后链码的安装和实例化，所以配不配置他的TLS证书也无关紧要了(关键是我们之前也没有将这个用户注册到tls服务器中)**

```shell
#指定TLS CA服务器生成的TLS根证书
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/tls/ca-cert.pem


#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer1/tls

fabric-ca-client enroll -d -u https://peer1-manufacturer:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-manufacturer

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer2/tls

fabric-ca-client enroll -d -u https://peer2-manufacturer:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-manufacturer

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/manufacturer/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network
```

启动peers

```shell
docker-compose up -d peer1-manufacturer
docker-compose up -d peer2-manufacturer
```



supplier msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer1


fabric-ca-client enroll -d -u https://peer1-supplier:peer1PW@0.0.0.0:7055


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer2


fabric-ca-client enroll -d -u https://peer2-supplier:peer2PW@0.0.0.0:7055


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/admin/msp

fabric-ca-client enroll -d -u https://admin-supplier:adminPW@0.0.0.0:7055

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer1/msp/admincerts/supplier-admin-cert.pem


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer2/msp/admincerts/supplier-admin-cert.pem



```

supplier tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer1/tls

fabric-ca-client enroll -d -u https://peer1-supplier:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-supplier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer2/tls

fabric-ca-client enroll -d -u https://peer2-supplier:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-supplier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/supplier/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network

```

启动peers

```shell
docker-compose up -d peer1-supplier
docker-compose up -d peer2-supplier
```



carrier msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer1


fabric-ca-client enroll -d -u https://peer1-carrier:peer1PW@0.0.0.0:7056


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer2


fabric-ca-client enroll -d -u https://peer2-carrier:peer2PW@0.0.0.0:7056


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/admin/msp

fabric-ca-client enroll -d -u https://admin-carrier:adminPW@0.0.0.0:7056

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer1/msp/admincerts/carrier-admin-cert.pem


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer2/msp/admincerts/carrier-admin-cert.pem


```

carrier tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer1/tls

fabric-ca-client enroll -d -u https://peer1-carrier:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-carrier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer2/tls

fabric-ca-client enroll -d -u https://peer2-carrier:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-carrier

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/carrier/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network


```

启动peers

```shell
docker-compose up -d peer1-carrier
docker-compose up -d peer2-carrier
```



middleman msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/ca-cert.pem
# peer1
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer1


fabric-ca-client enroll -d -u https://peer1-middleman:peer1PW@0.0.0.0:7057


# peer2
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer2


fabric-ca-client enroll -d -u https://peer2-middleman:peer2PW@0.0.0.0:7057


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/admin/msp

fabric-ca-client enroll -d -u https://admin-middleman:adminPW@0.0.0.0:7057

mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer1/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer1/msp/admincerts/middleman-admin-cert.pem


mkdir -p /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer2/msp/admincerts

cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer2/msp/admincerts/middleman-admin-cert.pem



```

middleman tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer1/tls

fabric-ca-client enroll -d -u https://peer1-middleman:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-middleman

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer1/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer2/tls

fabric-ca-client enroll -d -u https://peer2-middleman:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-middleman

cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/middleman/peer2/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network


```

启动peers

```shell
docker-compose up -d peer1-middleman
docker-compose up -d peer2-middleman
```



orderer配置

```shell
# orderer

export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/orderer

export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/ca-cert.pem

fabric-ca-client enroll -d -u https://orderer-cbpm:ordererPW@0.0.0.0:7053

# 指定TLS CA服务器生成的TLS根证书
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/orderer/tls
#指定TLS根证书
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/tls/ca-cert.pem

fabric-ca-client enroll -d -u https://orderer-cbpm:ordererPW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts orderer-cbpm


cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/orderer/tls/keystore
mv *_sk key.pem
cd /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network


# admin
export FABRIC_CA_CLIENT_HOME=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/admin
export FABRIC_CA_CLIENT_TLS_CERTFILES=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=/Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/admin/msp

fabric-ca-client enroll -d -u https://admin-cbpm:adminPW@0.0.0.0:7053

mkdir /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/orderer/msp/admincerts
#将签名证书拷贝过去
cp /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/admin/msp/signcerts/cert.pem /Users/cynyard/Documents/College/2022-Spring/graduation-project/cBPM-blockchain-system/network/cbpm/orderer/msp/admincerts/orderer-admin-cert.pem

```





MSP创世区块配置

```shell
#### ！！！一定要配置orderer的msp，不然一直报错panic: proto: Marshal called with nil
# orderer
cd ./cbpm/msp
mkdir admincerts && mkdir tlscacerts 
cd ..
cp ./admin/msp/signcerts/cert.pem ./msp/admincerts/ca-cert.pem
cp ./ca-cert.pem ./msp/cacerts/ca-cert.pem
cp ../tls/ca-cert.pem ./msp/tlscacerts/ca-cert.pem
cd ..


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

创建`configtx.yaml`  capabilities不用v2_0，不然会出问题

```shell

```



生成创世区块

```shell
export FABRIC_CFG_PATH=$PWD
configtxgen -profile CBPMOrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID syschannel
```

生成通道信息

```shell
configtxgen -profile SCChannel -outputCreateChannelTx ./channel-artifacts/scchannel.tx -channelID scchannel
```

配置config.yml

在\$PROJECT_PATH/\$ORG/admin/msp下面统一创建并配置（不用配置orderer的组织）

Certificate路径按情况修改

```shell
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/0-0-0-0-7057.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/0-0-0-0-7057.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/0-0-0-0-7057.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/0-0-0-0-7057.pem
    OrganizationalUnitIdentifier: orderer
```



运行orderer

```
docker-compose up -d orderer-cbpm
```



# 每次docker-compose up过后需要执行



配置cli

启动cli并且进入

```shell
docker-compose up -d 

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







安装测试链码

```shell
# cli
# peer lifecycle chaincode package test.tar.gz --path /tmp/hyperledger/fabric/chaincode/test-java --lang java --label test

# supplier-peer1环境 MSPCONFIGPATH设置为admin的

export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_ADDRESS=peer1-supplier:7051
export CORE_PEER_LOCALMSPID=SupplierMSP

peer chaincode install -n sccc -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/test-java

export CORE_PEER_ADDRESS=peer2-supplier:7051
peer chaincode install -n sccc -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/test-java


export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_ADDRESS=peer1-carrier:7051
export CORE_PEER_LOCALMSPID=CarrierMSP
peer chaincode install -n sccc -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/test-java


export CORE_PEER_ADDRESS=peer2-carrier:7051

peer chaincode install -n sccc -v 1.0 -l java -p /tmp/hyperledger/fabric/chaincode/test-java

```

初始化链码

```shell
# cli
peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile /tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem -C scchannel -n sccc -l java -v 1.0 -c '{"Args":[]}' -P 'OR ('\''SupplierMSP.peer'\'','\''CarrierMSP.peer'\'')'
```

- P：指定背书策略，上例中的背书策略是，两个组织需要参与链码invoke或query，chaincode执行才能生效。这个参数可为空，则任意安装了链码的节点无约束地调用链码。
- -n：指定链码名称
- -v：指定链码版本号
- -l：指定链码使用的语言，可以是golang, java, nodejs
- -o：指定排序节点
- -C：指定通道名
- -c：指定初始化参数

如果报错error starting container: Failed to generate platform-specific docker build: Error executing build: API error (404): network xxxx not found ""

新开terminal 输入`docker network ls` 然后根据网络名修改`docker-compose.yml` 中的`CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE` 环境变量





```shell
export ASSET_PROPERTIES=$(echo -n "{\"objectType\":\"asset\",\"assetID\":\"asset1\",\"color\":\"green\",\"size\":20,\"appraisedValue\":100}" | base64 | tr -d \\n)
peer chaincode invoke -o orderer-cbpm:7050 --ordererTLSHostnameOverride orderer-cbpm --tls --cafile /tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem -C channel -n private -c '{"function":"CreateAsset","Args":[]}' --transient "{\"asset_properties\":\"$ASSET_PROPERTIES\"}"

peer chaincode query -C scchannel -n sccc -c '{"f
unction":"CreateAsset","Args":[]}' --transient "{\"asset_properties\":\"$ASSET_PROPERTIES\"}"
```

