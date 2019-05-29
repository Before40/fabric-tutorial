# fabric-tutorial
Learning Fabric
该教程分为下面几个部分
- Fabric介绍
- 环境搭建
- 开发

## Fabric介绍
后续完善
## 环境搭建
搭建以下几种环境用来满足不同的需求

1. 2Peer + 1Orderer

主要用来完成基本的开发

2. 1CA + 2Peer + 1Orderer

CA测试

3. 生产环境

## 开发

### 开发环境搭建
Fabric 版本 1.4 LTS版本

开发环境采用1 ca + 2 Peer + 1 orderer结构

使用docker-compose编排环境
#### configgen.yaml文件
使用cryptogen showtemplate > configgen.yaml,生成模板文件，修改如下
```yaml
OrdererOrgs:
  - Name: Orderer
    Domain: example.com
    Specs:
      - Hostname: orderer
        SANS:
          - "localhost"
          - "127.0.0.1"
PeerOrgs:

  - Name: Org1
    Domain: org1.example.com
    EnableNodeOUs: false
    Template:
      Count: 2
      SANS:
        - "localhost"
        - "127.0.0.1"
    Users:
      Count: 1
```
执行命令：

cryptogen generate --config=cryptogen.yaml

然后在当前目录下生成crypto-config文件夹，包含相关的安全证书等。


### 交易配置
主要是编写configtx.yaml文件，然后利用configtx工具生成区块链相关配置。

#### configtx.yaml文件分析
这个文件主要用来配置创世块和联盟信息的

**Organizations** 定义了在后续配置中将要引用的不同组织的标识。
```yaml
Organizations:
    - &OrdererOrg
        Name: OrdererMSP
        ID: OrdererMSP
        MSPDir: ./crypto-config/ordererOrganizations/example.com/msp
        Policies: &OrdererOrgPolicies
            Readers:
              Type: Signature
              Rule: "OR('OrdererMSP.member')"                
            Writers:
              Type: Signature
              Rule: "OR('OrdererMSP.member')"
            Admins:
              Type: Signature
              Rule: "OR('OrdererMSP.admin')"
    - &Org1
        Name: Org1MSP
        ID: Org1MSP
        MSPDir: ./crypto-config/peerOrganizations/org1.example.com/msp
        AdminPrincipal: Role.ADMIN
        Policies: &Org1Policies
            Readers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
            Writers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
            Admins:
                Type: Signature
                Rule: "OR('Org1MSP.admin')"
        AnchorPeers:
            - Host: peer0.org1.example.com
              Port: 7051
```


**Orderer**定义了需要编码进配置交易或创世块中的与orderer相关的值
```yaml
Orderer: &OrdererDefaults
    OrdererType: solo
    Addresses:
        - orderer.example.com:7050
    BatchTimeout: 500ms

    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 98 MB
        PreferredMaxBytes: 512 KB
    MaxChannels: 0
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
    Organizations:

```


**Channel**与通道相关的配置
```yaml
Channel: &ChannelDefaults
    # Policies defines the set of policies at this level of the config tree
    # For Channel policies, their canonical path is
    #   /Channel/<PolicyName>
    Policies:
        # Who may invoke the 'Deliver' API
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        # Who may invoke the 'Broadcast' API
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        # By default, who may modify elements at this config level
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
    # Capabilities describes the channel level capabilities, see the
    # dedicated Capabilities section elsewhere in this file for a full
    # description
    Capabilities:
        <<: *ChannelCapabilities
```


**Profile**configtxgen工具的参数，不同的配置文件在这里配置
```yaml
Profiles:

    Profiles:

    SampleOrg:
      Orderer:
        <<: *OrdererDefaults
        Organizations:
          - *OrdererOrg
      Application:
        <<: *ApplicationDefaults
        Organizations:
          - *Org1
      Consortium:  SampleConsortium
      Consortiums:
        SampleConsortium:
          Organizations:
            - *Org1
            - *OrdererOrg
```

第一个configtx.yaml配置文件如下：
```yaml
Organizations:
    - &OrdererOrg
        Name: OrdererMSP
        ID: OrdererMSP
        MSPDir: ./crypto-config/ordererOrganizations/example.com/msp
        Policies: &OrdererOrgPolicies
            Readers:
              Type: Signature
              Rule: "OR('OrdererMSP.member')"                
            Writers:
              Type: Signature
              Rule: "OR('OrdererMSP.member')"
            Admins:
              Type: Signature
              Rule: "OR('OrdererMSP.admin')"
    - &Org1
        Name: Org1MSP
        ID: Org1MSP
        MSPDir: ./crypto-config/peerOrganizations/org1.example.com/msp
        AdminPrincipal: Role.ADMIN
        Policies: &Org1Policies
            Readers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
            Writers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
            Admins:
                Type: Signature
                Rule: "OR('Org1MSP.admin')"
        AnchorPeers:
            - Host: peer0.org1.example.com
              Port: 7051

Orderer: &OrdererDefaults
    OrdererType: solo
    Addresses:
        - orderer.example.com:7050
    BatchTimeout: 500ms

    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 98 MB
        PreferredMaxBytes: 512 KB
    MaxChannels: 0
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
    Organizations:
   
Channel: &ChannelDefaults
    Policies:
        # Who may invoke the 'Deliver' API
        Readers:
            Type: ImplicitMeta
            Rule: "ANY Readers"
        # Who may invoke the 'Broadcast' API
        Writers:
            Type: ImplicitMeta
            Rule: "ANY Writers"
        # By default, who may modify elements at this config level
        Admins:
            Type: ImplicitMeta
            Rule: "MAJORITY Admins"
    # Capabilities:
    #     <<: *ChannelCapabilities
Application: &ApplicationDefaults
  Organizations:

Profiles:

    SampleOrg:
      Orderer:
        <<: *OrdererDefaults
        Organizations:
          - *OrdererOrg
      Application:
        <<: *ApplicationDefaults
        Organizations:
          - *Org1
      Consortium:  SampleConsortium
      Consortiums:
        SampleConsortium:
          Organizations:
            - *Org1
            - *OrdererOrg
```

1.创建创世块

执行 configtxgen -profile SampleOrg -outputBlock ./channel-artifacts/gensis.block

2. 通道配置

configtxgen -profile SampleOrg -outputCreateChannelTx ./channel-artifacts/samplechannel.tx -channelID samplechannel

至此，配置文件生成结束

### Docker虚拟机配置
我们采用了1Orderer + 2Peers配置，docker-compose.yaml内容如下：
```yaml
version: '2'

networks:
  default:

services:

  orderer.example.com:
    image: hyperledger/fabric-orderer:amd64-1.4.1
    container_name: orderer.example.com
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_LISTENPORT=7050
      - ORDERER_GENERAL_GENESISPROFILE=SampleChain
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./channel-artifacts/gensis.block:/var/hyperledger/orderer/genesis.block
      - ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp:/var/hyperledger/orderer/msp
      - ./crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls:/var/hyperledger/orderer/tls
    ports:
      - 7050:7050
    networks:
      default:
        aliases:
          - orderer.example.com

  peer0.org1.example.com:
    container_name: peer0.org1.example.com
    image: hyperledger/fabric-peer:amd64-1.4.1
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_ATTACHSTDOUT=true
      - CORE_LOGGING_LEVEL=DEBUG
      - CORE_PEER_NETWORKID=SampleChain
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/var/hyperledger/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/var/hyperledger/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/tls/ca.crt
      - CORE_PEER_ID=peer0.org1.example.com
      - CORE_PEER_ADDRESSAUTODETECT=true
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org1.example.com:7051      
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_GOSSIP_SKIPHANDSHAKE=true
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_PEER_MSPCONFIGPATH=/var/hyperledger/msp
      - CORE_PEER_TLS_SERVERHOSTOVERRIDE=peer0.org1.example.com
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    volumes:
     - /var/run:/host/var/run
     - ./crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/var/hyperledger/msp
     - ./crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls:/var/hyperledger/tls
    ports:
      - 7051:7051
      - 7053:7053
    depends_on:
      - orderer.example.com
    links:
      - orderer.example.com
    networks:
      default:
        aliases:
          - peer0.org1.example.com

  peer1.org1.example.com:
    container_name: peer1.org1.example.com
    image: hyperledger/fabric-peer:amd64-1.4.1
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_VM_DOCKER_ATTACHSTDOUT=true
      - CORE_LOGGING_LEVEL=DEBUG
      - CORE_PEER_NETWORKID=SampleChain
      - CORE_PEER_PROFILE_ENABLED=true
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/var/hyperledger/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/var/hyperledger/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/var/hyperledger/tls/ca.crt
      - CORE_PEER_ID=peer1.org1.example.com
      - CORE_PEER_ADDRESSAUTODETECT=true
      - CORE_PEER_ADDRESS=peer1.org1.example.com:7051
      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.org1.example.com:7051      
      - CORE_PEER_GOSSIP_USELEADERELECTION=true
      - CORE_PEER_GOSSIP_ORGLEADER=false
      - CORE_PEER_GOSSIP_SKIPHANDSHAKE=true
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_PEER_MSPCONFIGPATH=/var/hyperledger/msp
      - CORE_PEER_TLS_SERVERHOSTOVERRIDE=peer1.org1.example.com
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: peer node start
    volumes:
     - /var/run:/host/var/run
     - ./crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/msp:/var/hyperledger/msp
     - ./crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls:/var/hyperledger/tls
    ports:
      - 8051:7051
      - 8053:7053
    depends_on:
      - orderer.example.com
    links:
      - orderer.example.com
    networks:
      default:
        aliases:
          - peer1.org1.example.com
             
```

执行docker-compose up,
然后查看日志，可以看出启动成功


### Fabsdk测试
修改chainhero代码

go build编译

执行，然后在浏览器中访问就能查询和修改。
