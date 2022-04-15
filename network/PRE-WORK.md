# 搭建网络+CA

## 环境变量

> 有些命令必须要用绝对路径

```shell
export PROJECT_PATH=$PROJECT_PATH
```



## CA配置

ca服务器配置

```shell
#设置环境变量指定根证书的路径(如果工作目录不同的话记得指定自己的工作目录,以下不再重复说明)
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/tls/ca-cert.pem
#设置环境变量指定CA客户端的HOME文件夹
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/tls/admin
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
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/cbpm/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/cbpm/admin

# 登陆
fabric-ca-client enroll -d -u https://cbpm-ca-admin:cbpm-adminpw@0.0.0.0:7053

#注册
fabric-ca-client register -d --id.name orderer-cbpm --id.secret ordererPW --id.type orderer -u https://0.0.0.0:7053

fabric-ca-client register -d --id.name admin-cbpm --id.secret adminPW --id.type admin --id.attrs "hf.Registrar.Roles=client,hf.Registrar.Attributes=*,hf.Revoker=true,hf.GenCRL=true,admin=true:ecert,abac.init=true:ecert" -u https://0.0.0.0:7053


```





manufacturer-ca

```shell
docker-compose up manufacturer-ca

export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/manufacturer/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/manufacturer/admin


fabric-ca-client enroll -d -u https://manufacturer-ca-admin:manufacturer-adminpw@0.0.0.0:7054




fabric-ca-client register -d --id.name peer1-manufacturer --id.secret peer1PW --id.type peer -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name peer2-manufacturer --id.secret peer2PW --id.type peer -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name admin-manufacturer --id.secret adminPW --id.type admin -u https://0.0.0.0:7054
fabric-ca-client register -d --id.name user-manufacturer --id.secret userPW --id.type user -u https://0.0.0.0:7054


```





supplier-ca

```shell
docker-compose up supplier-ca



export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/supplier/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/supplier/admin


fabric-ca-client enroll -d -u https://supplier-ca-admin:supplier-adminpw@0.0.0.0:7055


fabric-ca-client register -d --id.name peer1-supplier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name peer2-supplier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name admin-supplier --id.secret adminPW --id.type admin -u https://0.0.0.0:7055
fabric-ca-client register -d --id.name user-supplier --id.secret userPW --id.type user -u https://0.0.0.0:7055
```



carrier-ca

```shell
docker-compose up carrier-ca



export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/carrier/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/carrier/admin


fabric-ca-client enroll -d -u https://carrier-ca-admin:carrier-adminpw@0.0.0.0:7056


fabric-ca-client register -d --id.name peer1-carrier --id.secret peer1PW --id.type peer -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name peer2-carrier --id.secret peer2PW --id.type peer -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name admin-carrier --id.secret adminPW --id.type admin -u https://0.0.0.0:7056
fabric-ca-client register -d --id.name user-carrier --id.secret userPW --id.type user -u https://0.0.0.0:7056
```

middleman

```shell
docker-compose up middleman



export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/middleman/ca-cert.pem
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/middleman/admin


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
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/manufacturer/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/manufacturer/peer1


fabric-ca-client enroll -d -u https://peer1-manufacturer:peer1PW@0.0.0.0:7054


# peer2
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/manufacturer/peer2


fabric-ca-client enroll -d -u https://peer2-manufacturer:peer2PW@0.0.0.0:7054


# admin
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/manufacturer/admin
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/manufacturer/admin/msp

fabric-ca-client enroll -d -u https://admin-manufacturer:adminPW@0.0.0.0:7054


mkdir -p $PROJECT_PATH/manufacturer/peer1/msp/admincerts

cp $PROJECT_PATH/manufacturer/admin/msp/signcerts/cert.pem $PROJECT_PATH/manufacturer/peer1/msp/admincerts/manufacturer-admin-cert.pem

mkdir -p $PROJECT_PATH/manufacturer/peer2/msp/admincerts

cp $PROJECT_PATH/manufacturer/admin/msp/signcerts/cert.pem $PROJECT_PATH/manufacturer/peer2/msp/admincerts/manufacturer-admin-cert.pem

```



manufacturer tls 

**admin不用配置tls是因为我们生成admin的证书主要就是为了之后链码的安装和实例化，所以配不配置他的TLS证书也无关紧要了(关键是我们之前也没有将这个用户注册到tls服务器中)**

```shell
#指定TLS CA服务器生成的TLS根证书
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/tls/ca-cert.pem


#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/manufacturer/peer1/tls

fabric-ca-client enroll -d -u https://peer1-manufacturer:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-manufacturer

cd $PROJECT_PATH/manufacturer/peer1/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/manufacturer/peer2/tls

fabric-ca-client enroll -d -u https://peer2-manufacturer:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-manufacturer

cd $PROJECT_PATH/manufacturer/peer2/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH
```

启动peers

```shell
docker-compose up -d peer1-manufacturer
docker-compose up -d peer2-manufacturer
```



supplier msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/supplier/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/supplier/peer1


fabric-ca-client enroll -d -u https://peer1-supplier:peer1PW@0.0.0.0:7055


# peer2
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/supplier/peer2


fabric-ca-client enroll -d -u https://peer2-supplier:peer2PW@0.0.0.0:7055


# admin
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/supplier/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/supplier/admin/msp

fabric-ca-client enroll -d -u https://admin-supplier:adminPW@0.0.0.0:7055

mkdir -p $PROJECT_PATH/supplier/peer1/msp/admincerts

cp $PROJECT_PATH/supplier/admin/msp/signcerts/cert.pem $PROJECT_PATH/supplier/peer1/msp/admincerts/supplier-admin-cert.pem


mkdir -p $PROJECT_PATH/supplier/peer2/msp/admincerts

cp $PROJECT_PATH/supplier/admin/msp/signcerts/cert.pem $PROJECT_PATH/supplier/peer2/msp/admincerts/supplier-admin-cert.pem



```

supplier tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/supplier/peer1/tls

fabric-ca-client enroll -d -u https://peer1-supplier:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-supplier

cd $PROJECT_PATH/supplier/peer1/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/supplier/peer2/tls

fabric-ca-client enroll -d -u https://peer2-supplier:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-supplier

cd $PROJECT_PATH/supplier/peer2/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH

```

启动peers

```shell
docker-compose up -d peer1-supplier
docker-compose up -d peer2-supplier
```



carrier msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/carrier/ca-cert.pem

# peer1
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/carrier/peer1


fabric-ca-client enroll -d -u https://peer1-carrier:peer1PW@0.0.0.0:7056


# peer2
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/carrier/peer2


fabric-ca-client enroll -d -u https://peer2-carrier:peer2PW@0.0.0.0:7056


# admin
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/carrier/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/carrier/admin/msp

fabric-ca-client enroll -d -u https://admin-carrier:adminPW@0.0.0.0:7056

mkdir -p $PROJECT_PATH/carrier/peer1/msp/admincerts

cp $PROJECT_PATH/carrier/admin/msp/signcerts/cert.pem $PROJECT_PATH/carrier/peer1/msp/admincerts/carrier-admin-cert.pem


mkdir -p $PROJECT_PATH/carrier/peer2/msp/admincerts

cp $PROJECT_PATH/carrier/admin/msp/signcerts/cert.pem $PROJECT_PATH/carrier/peer2/msp/admincerts/carrier-admin-cert.pem


```

carrier tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/carrier/peer1/tls

fabric-ca-client enroll -d -u https://peer1-carrier:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-carrier

cd $PROJECT_PATH/carrier/peer1/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/carrier/peer2/tls

fabric-ca-client enroll -d -u https://peer2-carrier:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-carrier

cd $PROJECT_PATH/carrier/peer2/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH


```

启动peers

```shell
docker-compose up -d peer1-carrier
docker-compose up -d peer2-carrier
```



middleman msp

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/middleman/ca-cert.pem
# peer1
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/middleman/peer1


fabric-ca-client enroll -d -u https://peer1-middleman:peer1PW@0.0.0.0:7057


# peer2
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/middleman/peer2


fabric-ca-client enroll -d -u https://peer2-middleman:peer2PW@0.0.0.0:7057


# admin
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/middleman/admin

#这里多了一个环境变量，是指定admin用户的msp证书文件夹的
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/middleman/admin/msp

fabric-ca-client enroll -d -u https://admin-middleman:adminPW@0.0.0.0:7057

mkdir -p $PROJECT_PATH/middleman/peer1/msp/admincerts

cp $PROJECT_PATH/middleman/admin/msp/signcerts/cert.pem $PROJECT_PATH/middleman/peer1/msp/admincerts/middleman-admin-cert.pem


mkdir -p $PROJECT_PATH/middleman/peer2/msp/admincerts

cp $PROJECT_PATH/middleman/admin/msp/signcerts/cert.pem $PROJECT_PATH/middleman/peer2/msp/admincerts/middleman-admin-cert.pem



```

middleman tls

```shell
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/tls/ca-cert.pem
#指定TLS证书的HOME目录 peer1
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/middleman/peer1/tls

fabric-ca-client enroll -d -u https://peer1-middleman:peer1PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer1-middleman

cd $PROJECT_PATH/middleman/peer1/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH

#指定TLS证书的HOME目录 peer2
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/middleman/peer2/tls

fabric-ca-client enroll -d -u https://peer2-middleman:peer2PW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts peer2-middleman

cd $PROJECT_PATH/middleman/peer2/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH


```

启动peers

```shell
docker-compose up -d peer1-middleman
docker-compose up -d peer2-middleman
```



orderer配置

```shell
# orderer

export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/cbpm/orderer

export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/cbpm/ca-cert.pem

fabric-ca-client enroll -d -u https://orderer-cbpm:ordererPW@0.0.0.0:7053

# 指定TLS CA服务器生成的TLS根证书
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/cbpm/orderer/tls
#指定TLS根证书
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/tls/ca-cert.pem

fabric-ca-client enroll -d -u https://orderer-cbpm:ordererPW@0.0.0.0:7052 --enrollment.profile tls --csr.hosts orderer-cbpm


cd $PROJECT_PATH/cbpm/orderer/tls/keystore
mv *_sk key.pem
cd $PROJECT_PATH


# admin
export FABRIC_CA_CLIENT_HOME=$PROJECT_PATH/cbpm/admin
export FABRIC_CA_CLIENT_TLS_CERTFILES=$PROJECT_PATH/cbpm/ca-cert.pem
export FABRIC_CA_CLIENT_MSPDIR=$PROJECT_PATH/cbpm/admin/msp

fabric-ca-client enroll -d -u https://admin-cbpm:adminPW@0.0.0.0:7053

mkdir $PROJECT_PATH/cbpm/orderer/msp/admincerts
#将签名证书拷贝过去
cp $PROJECT_PATH/cbpm/admin/msp/signcerts/cert.pem $PROJECT_PATH/cbpm/orderer/msp/admincerts/orderer-admin-cert.pem

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
Organizations:
    - &OrdererOrg
        Name: OrdererOrg
        ID: OrdererMSP
        MSPDir: ./cbpm/msp
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('OrdererMSP.member')"
            Writers:
                Type: Signature
                Rule: "OR('OrdererMSP.member')"
            Admins:
                Type: Signature
                Rule: "OR('OrdererMSP.admin')"

        OrdererEndpoints:
            - orderer-cbpm:7050

    - &ManufacturerOrg
        Name: ManufacturerOrg
        ID: ManufacturerMSP
        MSPDir: ./manufacturer/msp
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('ManufacturerMSP.member')"
            Writers:
                Type: Signature
                Rule: "OR('ManufacturerMSP.member')"
            Admins:
                Type: Signature
                Rule: "OR('ManufacturerMSP.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('ManufacturerMSP.peer')"
        AnchorPeers:
            - Host: peer1-manufacturer
              Port: 7051


    - &SupplierOrg
        Name: SupplierOrg
        ID: SupplierMSP
        MSPDir: ./supplier/msp
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('SupplierMSP.admin', 'SupplierMSP.peer', 'SupplierMSP.client')"
            Writers:
                Type: Signature
                Rule: "OR('SupplierMSP.admin', 'SupplierMSP.client')"
            Admins:
                Type: Signature
                Rule: "OR('SupplierMSP.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('SupplierMSP.peer')"
        AnchorPeers:
            - Host: peer1-supplier
              Port: 7051

    - &CarrierOrg
        Name: CarrierOrg
        ID: CarrierMSP
        MSPDir: ./carrier/msp
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('CarrierMSP.admin', 'CarrierMSP.peer', 'CarrierMSP.client')"
            Writers:
                Type: Signature
                Rule: "OR('CarrierMSP.admin', 'CarrierMSP.client')"
            Admins:
                Type: Signature
                Rule: "OR('CarrierMSP.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('CarrierMSP.peer')"

        AnchorPeers:
            - Host: peer1-carrier
              Port: 7051

    - &MiddlemanOrg
        Name: MiddlemanOrg
        ID: MiddlemanMSP

        MSPDir: ./middleman/msp
        Policies:
            Readers:
                Type: Signature
                Rule: "OR('MiddlemanMSP.admin', 'MiddlemanMSP.peer', 'MiddlemanMSP.client')"
            Writers:
                Type: Signature
                Rule: "OR('MiddlemanMSP.admin', 'MiddlemanMSP.client')"
            Admins:
                Type: Signature
                Rule: "OR('MiddlemanMSP.admin')"
            Endorsement:
                Type: Signature
                Rule: "OR('MiddlemanMSP.peer')"

        AnchorPeers:
            - Host: peer1-middleman
              Port: 7051
Capabilities:
    Channel: &ChannelCapabilities
        V1_4_3: true
        V1_3: false
        V1_1: false
    Orderer: &OrdererCapabilities
        V1_4_2: true
        V1_1: false
    Application: &ApplicationCapabilities
        V1_4_2: true
        V1_3: false
        V1_2: false
        V1_1: false  
Application: &ApplicationDefaults
    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        LifecycleEndorsement:
            Type: ImplicitMeta
            Rule: "MAJORITY Endorsement"
        Endorsement:
            Type: ImplicitMeta
            Rule: "MAJORITY Endorsement"

    Capabilities:
        <<: *ApplicationCapabilities
Orderer: &OrdererDefaults
    OrdererType: solo

    Addresses:
        - orderer-cbpm:7050
    BatchTimeout: 2s
    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 99 MB
        PreferredMaxBytes: 512 KB
    Organizations:
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
        BlockValidation:
            Type: ImplicitMeta
            Rule: "ANY Writers"
Channel: &ChannelDefaults
    Policies:
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
    Capabilities:
        <<: *ChannelCapabilities
Profiles:
    CBPMOrdererGenesis:
        <<: *ChannelDefaults
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
            Capabilities: 
                <<: *OrdererCapabilities
        Consortiums:
            CBPMConsortium:
                Organizations:
                    - *SupplierOrg
                    - *CarrierOrg
                    - *ManufacturerOrg
                    - *MiddlemanOrg
            # SCConsortium:
            #     Organizations:
            #         - *SupplierOrg
            #         - *CarrierOrg
            # MaMiConsortium:
            #     Organizations:
            #         - *ManufacturerOrg
            #         - *MiddlemanOrg
            # MiSConsortium:
            #     Organizations:
            #         - *MiddlemanOrg
            #         - *SupplierOrg   
            # MiCConsortium:
            #     Organizations:
            #         - *MiddlemanOrg
            #         - *CarrierOrg 
            # CMaConsortium:
            #     Organizations:
            #         - *CarrierOrg
            #         - *ManufacturerOrg
    CBPMChannel:
        <<: *ChannelDefaults
        Consortium: CBPMConsortium
        Application:
            <<: *ApplicationDefaults
            Organizations:
                - *ManufacturerOrg
                - *MiddlemanOrg
                - *SupplierOrg
                - *CarrierOrg
            Capabilities: *ApplicationCapabilities
    # SCChannel:
    #     <<: *ChannelDefaults
    #     Consortium: SCConsortium
    #     Application:
    #         <<: *ApplicationDefaults
    #         Organizations:
    #             - *SupplierOrg
    #             - *CarrierOrg
    #         Capabilities: *ApplicationCapabilities

    # MaMiChannel:
    #     <<: *ChannelDefaults
    #     Consortium: MaMiConsortium
    #     Application:
    #         <<: *ApplicationDefaults
    #         Organizations:
    #             - *ManufacturerOrg
    #             - *MiddlemanOrg
    #         Capabilities: *ApplicationCapabilities
    # MiSChannel:
    #     <<: *ChannelDefaults
    #     Consortium: MiSConsortium
    #     Application:
    #         <<: *ApplicationDefaults
    #         Organizations:
    #             - *MiddlemanOrg
    #             - *SupplierOrg
    #         Capabilities: *ApplicationCapabilities
    # MiCChannel:
    #     <<: *ChannelDefaults
    #     Consortium: MiCConsortium
    #     Application:
    #         <<: *ApplicationDefaults
    #         Organizations:
    #             - *MiddlemanOrg
    #             - *CarrierOrg
    #         Capabilities: *ApplicationCapabilities
    # CMaChannel:
    #     <<: *ChannelDefaults
    #     Consortium: CMaConsortium
    #     Application:
    #         <<: *ApplicationDefaults
    #         Organizations:
    #             - *CarrierOrg
    #             - *ManufacturerOrg
    #         Capabilities: *ApplicationCapabilities

```

生成创世区块

```shell
configtxgen -profile CBPMOrdererGenesis -outputBlock ./channel-artifacts/genesis.block -channelID syschannel
```

生成通道信息

```shell
configtxgen -profile CBPMChannel -outputCreateChannelTx ./channel-artifacts/cbpmchannel.tx -channelID cbpmchannel
```

配置anchor-peers

```shell
configtxgen -profile CBPMChannel -outputAnchorPeersUpdate ./channel-artifacts/ManufacturerMSPanchors.tx -channelID cbpmchannel -asOrg ManufacturerOrg

configtxgen -profile CBPMChannel -outputAnchorPeersUpdate ./channel-artifacts/MiddlemanMSPanchors.tx -channelID cbpmchannel -asOrg MiddlemanOrg

configtxgen -profile CBPMChannel -outputAnchorPeersUpdate ./channel-artifacts/CarrierMSPanchors.tx -channelID cbpmchannel -asOrg CarrierOrg

configtxgen -profile CBPMChannel -outputAnchorPeersUpdate ./channel-artifacts/SupplierMSPanchors.tx -channelID cbpmchannel -asOrg SupplierOrg
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



# 测试网络

## Create a supplier-carrier channel

### 创建并进入通道

生成scchannel.tx启动cli并且进入

```shell
configtxgen -profile SCChannel -outputCreateChannelTx ./channel-artifacts/scchannel.tx -channelID scchannel
docker exec -it cli /bin/bash
```



用supplier创建scchannel通道

```shell
# 配置通道名
export CHANNEL=scchannel


# 配置supplier-peer1环境 MSPCONFIGPATH设置为admin的
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP


peer channel create -c $CHANNEL -f /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.tx -o orderer-cbpm:7050 --outputBlock /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.block --tls --cafile /tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem

#### 或者
peer channel create -c $CHANNEL -f /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.tx -o orderer-cbpm:7050 --outputBlock /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.block --tls --cafile /tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem
```

Supplier节点加入通道

```shell
# cli
export CORE_PEER_ADDRESS=peer1-supplier:7051
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.block
export CORE_PEER_ADDRESS=peer2-supplier:7051
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.block
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
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.block
export CORE_PEER_ADDRESS=peer2-carrier:7051
peer channel join -b /tmp/hyperledger/fabric/channel-artifacts/$CHANNEL.block
```

检测（如果没有的话会报错

```shell
peer channel list
```

### 选择链码

链码路径-链码语言

- chaincode-supplier-carrier: golang
- test-go : golang
- test-go-asset-transfer : golang
- test-go-private-data : golang

### 安装链码

```shell
# 输入想要安装的链码路径以及语言版本
export CHAINCODE=chaincode-supplier-carrier
export CHAINCODE_LANG=golang
export CHAINCODE_VERSION=1.0
export CHAINCODE_NAME=scchaincode


export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_LOCALMSPID=SupplierMSP

export CORE_PEER_ADDRESS=peer1-supplier:7051
peer chaincode install -n $CHAINCODE_NAME -v $CHAINCODE_VERSION  -p $CHAINCODE -l $CHAINCODE_LANG

export CORE_PEER_ADDRESS=peer2-supplier:7051
peer chaincode install -n $CHAINCODE_NAME -v $CHAINCODE_VERSION  -p $CHAINCODE -l $CHAINCODE_LANG


export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/keystore/key.pem
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_LOCALMSPID=CarrierMSP

export CORE_PEER_ADDRESS=peer1-carrier:7051
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
peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

peer chaincode upgrade -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"
```

#### test-go-private-data

```shell
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')" --collections-config $GOPATH/src/$CHAINCODE/collections_config.json

export MARBLE=$(echo -n "{\"name\":\"marble1\",\"color\":\"blue\",\"size\":35,\"owner\":\"tom\",\"price\":99}" | base64 | tr -d \\n)


peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["InitMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"

export MARBLE=$(echo -n "{\"name\":\"marble2\",\"color\":\"red\",\"size\":50,\"owner\":\"tom\",\"price\":102}" | base64 | tr -d \\n)

peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["InitMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"


peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["readMarble","marble1"]}'


peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["readMarblePrivateDetails","marble1"]}'



# switch to carrier
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/carrier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/carrier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_LOCALMSPID=CarrierMSP
export CORE_PEER_ADDRESS=peer1-carrier:7051

export MARBLE=$(echo -n "{\"name\":\"marble3\",\"color\":\"red\",\"size\":50,\"owner\":\"tom\",\"price\":222}" | base64 | tr -d \\n)

peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["InitMarble"]}' --transient "{\"marble\":\"$MARBLE\"}"

peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["readMarble","marble1"]}'


peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["readMarblePrivateDetails","marble1"]}'

peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["GetMarbleHash","collectionMarbles","marble1"]}'

peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["QueryMarblesByOwner","tom"]}'

peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["QueryMarbles","{\"selector\":{\"owner\":\"tom\"}}"]}'

```

#### go-marbles-couchdb

```shell
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051


peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":[""]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

```

调用查询

```shell
peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["initMarble","marble1","blue","35","tom"]}'

peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -c '{"function": "initMarble","Args":["marble1","blue","35","tom"]}'


peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"function": "queryMarbles","Args":["{\"selector\":{\"owner\":\"tom\"}}"]}'


peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["readMarble","marble1"]}'
```

#### test-go

```shell
export CORE_PEER_MSPCONFIGPATH=/tmp/hyperledger/fabric/peer/supplier/admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/tlscacerts/tls-0-0-0-0-7052.pem
export CORE_PEER_TLS_CERT_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/signcerts/cert.pem
export CORE_PEER_TLS_KEY_FILE=/tmp/hyperledger/fabric/peer/supplier/peer1/tls/keystore/key.pem
export CORE_PEER_LOCALMSPID=SupplierMSP
export CORE_PEER_ADDRESS=peer1-supplier:7051

peer chaincode instantiate -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -l $CHAINCODE_LANG -v $CHAINCODE_VERSION -c '{"Args":["initLedger"]}' -P "OR('SupplierMSP.peer','CarrierMSP.peer')"

```



```shell

peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["ReadAsset","asset1"]}'


peer chaincode invoke -o orderer-cbpm:7050 --tls --cafile "/tmp/hyperledger/fabric/peer/cbpm/orderer/tls/tlscacerts/tls-0-0-0-0-7052.pem" -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["CreateAsset","asset10","black","5","cynyard","10"]}'


peer chaincode query -C $CHANNEL -n $CHAINCODE_NAME -c '{"Args":["ReadAsset","asset10"]}'


```
