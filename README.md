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

PeerOrgs:

  - Name: Org1
    Domain: org1.example.com
    EnableNodeOUs: false
    Template:
      Count: 2
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

    # SampleOrg defines an MSP using the sampleconfig.  It should never be used
    # in production but may be used as a template for other definitions
    - &OrdererOrg
        # DefaultOrg defines the organization which is used in the sampleconfig
        # of the fabric.git development environment
        Name: OrdererMSP

        # ID to load the MSP definition as
        ID: OrdererMSP

        # MSPDir is the filesystem path which contains the MSP configuration
        MSPDir: ../../v1/crypto-config/ordererOrganizations/example.com/msp

        # Policies defines the set of policies at this level of the config tree
        # For organization policies, their canonical path is usually
        #   /Channel/<Application|Orderer>/<OrgName>/<PolicyName>
        Policies: &OrdererOrgPolicies
            Readers:
                Type: Signature
                Rule: "OR('OrdererMSP.member')"
                # If your MSP is configured with the new NodeOUs, you might
                # want to use a more specific rule like the following:
                # Rule: "OR('OrdererMSP.admin', 'OrdererMSP.peer')"
            Writers:
                Type: Signature
                Rule: "OR('OrdererMSP.member')"
                # If your MSP is configured with the new NodeOUs, you might
                # want to use a more specific rule like the following:
                # Rule: "OR('OrdererMSP.admin', 'OrdererMSP.client'')"
            Admins:
                Type: Signature
                Rule: "OR('OrdererMSP.admin')"

    - &Org1
        # DefaultOrg defines the organization which is used in the sampleconfig
        # of the fabric.git development environment
        Name: Org1MSP

        # ID to load the MSP definition as
        ID: Org1MSP

        MSPDir: ../../v1/crypto-config/peerOrganizations/org1.example.com/msp

        # Policies defines the set of policies at this level of the config tree
        # For organization policies, their canonical path is usually
        #   /Channel/<Application|Orderer>/<OrgName>/<PolicyName>
        Policies: &Org1Policies
            Readers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
                # If your MSP is configured with the new NodeOUs, you might
                # want to use a more specific rule like the following:
                # Rule: "OR('Org1MSP.admin', 'Org1MSP.peer')"
            Writers:
                Type: Signature
                Rule: "OR('Org1MSP.member')"
                # If your MSP is configured with the new NodeOUs, you might
                # want to use a more specific rule like the following:
                # Rule: "OR('Org1MSP.admin', 'Org1MSP.client'')"
            Admins:
                Type: Signature
                Rule: "OR('Org1MSP.admin')"

        AnchorPeers:
            # AnchorPeers defines the location of peers which can be used
            # for cross org gossip communication.  Note, this value is only
            # encoded in the genesis block in the Application section context
            - Host: peer0.org1.example.com
              Port: 7051
```

**Capabilities**定义fabric网络的能力。v1.1的新概念，不能和v1.0.x的混合使用。
Capabilities定义了在fabric中必须实现的特征，fabric binary安全加入到fabric网络中。如，加入了一个新的MSP类型，新的binaries会识别和验证签名，而老的binariers没有这个支持的不能验证这个交易。这个会导致不同版本的fabric binaries有不同的世界状态。相反，定义一个通道，他的功能是通知其他没有这个能力的binaries必须停止处理交易直到他们更新为止。
```yaml
Capabilities:
    # Global capabilities apply to both the orderers and the peers and must be
    # supported by both.  Set the value of the capability to true to require it.
    Channel: &ChannelCapabilities
        # V1.1 for Global is a catchall flag for behavior which has been
        # determined to be desired for all orderers and peers running v1.0.x,
        # but the modification of which would cause imcompatibilities.  Users
        # should leave this flag set to true.
        V1_1: true

    # Orderer capabilities apply only to the orderers, and may be safely
    # manipulated without concern for upgrading peers.  Set the value of the
    # capability to true to require it.
    Orderer: &OrdererCapabilities
        # V1.1 for Order is a catchall flag for behavior which has been
        # determined to be desired for all orderers running v1.0.x, but the
        # modification of which  would cause imcompatibilities.  Users should
        # leave this flag set to true.
        V1_1: true

    # Application capabilities apply only to the peer network, and may be safely
    # manipulated without concern for upgrading orderers.  Set the value of the
    # capability to true to require it.
    Application: &ApplicationCapabilities
        # V1.1 for Application is a catchall flag for behavior which has been
        # determined to be desired for all peers running v1.0.x, but the
        # modification of which would cause incompatibilities.  Users should
        # leave this flag set to true.
        V1_2: true
```
**Orderer**定义了需要编码进配置交易或创世块中的与orderer相关的值
```yaml
Orderer: &OrdererDefaults

    # Orderer Type: The orderer implementation to start
    # Available types are "solo" and "kafka"
    OrdererType: solo

    Addresses:
        - orderer.example.com:7050

    # Batch Timeout: The amount of time to wait before creating a batch
    BatchTimeout: 500ms

    # Batch Size: Controls the number of messages batched into a block
    BatchSize:

        # Max Message Count: The maximum number of messages to permit in a batch
        MaxMessageCount: 10

        # Absolute Max Bytes: The absolute maximum number of bytes allowed for
        # the serialized messages in a batch.
        AbsoluteMaxBytes: 98 MB

        # Preferred Max Bytes: The preferred maximum number of bytes allowed for
        # the serialized messages in a batch. A message larger than the preferred
        # max bytes will result in a batch larger than preferred max bytes.
        PreferredMaxBytes: 512 KB

    # Max Channels is the maximum number of channels to allow on the ordering
    # network. When set to 0, this implies no maximum number of channels.
    MaxChannels: 0

    Kafka:
        # Brokers: A list of Kafka brokers to which the orderer connects
        # NOTE: Use IP:port notation
        Brokers:
            - 127.0.0.1:9092

    # Organizations is the list of orgs which are defined as participants on
    # the orderer side of the network
    Organizations:

    # Policies defines the set of policies at this level of the config tree
    # For Orderer policies, their canonical path is
    #   /Channel/Orderer/<PolicyName>
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
        # BlockValidation specifies what signatures must be included in the block
        # from the orderer for the peer to validate it.
        BlockValidation:
            Type: ImplicitMeta
            Rule: "ANY Writers"

    # Capabilities describes the orderer level capabilities, see the
    # dedicated Capabilities section elsewhere in this file for a full
    # description
    Capabilities:
        <<: *OrdererCapabilities
```
**Applications**和应用相关的参数需要加入到配置交易或者创世块的值
```yaml
Application: &ApplicationDefaults
    ACLs: &ACLsDefault
        #This section provides defaults for policies for various resources
        #in the system. These "resources" could be functions on system chaincodes
        #(e.g., "GetBlockByNumber" on the "qscc" system chaincode) or other resources
        #(e.g.,who can receive Block events). This section does NOT specify the resource's
        #definition or API, but just the ACL policy for it.
        #
        #User's can override these defaults with their own policy mapping by defining the
        #mapping under ACLs in their channel definition

        #---Lifecycle System Chaincode (lscc) function to policy mapping for access control---#

        #ACL policy for lscc's "getid" function
        lscc/ChaincodeExists: /Channel/Application/Readers

        #ACL policy for lscc's "getdepspec" function
        lscc/GetDeploymentSpec: /Channel/Application/Readers

        #ACL policy for lscc's "getccdata" function
        lscc/GetChaincodeData: /Channel/Application/Readers

        #---Query System Chaincode (qscc) function to policy mapping for access control---#

        #ACL policy for qscc's "GetChainInfo" function
        qscc/GetChainInfo: /Channel/Application/Readers

        #ACL policy for qscc's "GetBlockByNumber" function
        qscc/GetBlockByNumber: /Channel/Application/Readers

        #ACL policy for qscc's  "GetBlockByHash" function
        qscc/GetBlockByHash: /Channel/Application/Readers

        #ACL policy for qscc's "GetTransactionByID" function
        qscc/GetTransactionByID: /Channel/Application/Readers

        #ACL policy for qscc's "GetBlockByTxID" function
        qscc/GetBlockByTxID: /Channel/Application/Readers

        #---Configuration System Chaincode (cscc) function to policy mapping for access control---#

        #ACL policy for cscc's "GetConfigBlock" function
        cscc/GetConfigBlock: /Channel/Application/Readers

        #ACL policy for cscc's "GetConfigTree" function
        cscc/GetConfigTree: /Channel/Application/Readers

        #ACL policy for cscc's "SimulateConfigTreeUpdate" function
        cscc/SimulateConfigTreeUpdate: /Channel/Application/Writers

        #---Miscellanesous peer function to policy mapping for access control---#

        #ACL policy for invoking chaincodes on peer
        peer/Proposal: /Channel/Application/Writers

        #ACL policy for chaincode to chaincode invocation
        peer/ChaincodeToChaincode: /Channel/Application/Readers

        #---Events resource to policy mapping for access control###---#

        #ACL policy for sending block events
        event/Block: /Channel/Application/Readers

        #ACL policy for sending filtered block events
        event/FilteredBlock: /Channel/Application/Readers

    # Organizations is the list of orgs which are defined as participants on
    # the application side of the network.
    Organizations:

    # Policies defines the set of policies at this level of the config tree
    # For Application policies, their canonical path is
    #   /Channel/Application/<PolicyName>
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
        Org1MemberPolicy:
            Type: Signature
            Rule: "OR('Org1MSP.member')"
        Org2MemberPolicy:
            Type: Signature
            Rule: "OR('Org2MSP.member')"
        Org1Org2MemberPolicy:
            Type: Signature
            Rule: "OR('Org1MSP.member','Org2MSP.member')"

    # Capabilities describes the application level capabilities, see the
    # dedicated Capabilities section elsewhere in this file for a full
    # description
    Capabilities:
        <<: *ApplicationCapabilities
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

    TwoOrgsOrdererGenesis:
        <<: *ChannelDefaults
        Orderer:
            <<: *OrdererDefaults
            Organizations:
            - <<: *OrdererOrg
              Policies:
                  <<: *OrdererOrgPolicies
                  Admins:
                      Type: Signature
                      Rule: "OR('OrdererMSP.admin')"
        Consortiums:
            TestConsortium:
                Organizations:
                - <<: *Org1
                  Policies:
                      <<: *Org1Policies
                      Admins:
                          Type: Signature
                          Rule: "OR('Org1MSP.admin')"

            SampleConsortium:
                Organizations:
                - <<: *Org1
                  Policies:
                      <<: *Org1Policies
                      Admins:
                          Type: Signature
                          Rule: "OR('Org1MSP.admin')"
                - <<: *Org2
                  Policies:
                      <<: *Org2Policies
                      Admins:
                          Type: Signature
                          Rule: "OR('Org2MSP.admin')"

    OneOrgChannel:
        Consortium: TestConsortium
        Application:
            <<: *ApplicationDefaults
            Organizations:
            - *Org1

    TwoOrgsChannel:
        Consortium: SampleConsortium
        Application:
            <<: *ApplicationDefaults
            Organizations:
            - *Org1
            - *Org2

    DsChannel:
        Consortium: SampleConsortium
        Application:
            <<: *ApplicationDefaults
            Organizations:
            - *Org1
            - *Org2
```

第一个configtx.yaml配置文件如下：
```yaml
Organizations:
    - &OrdererOrg
        Name: OrdererMSP
        ID: OrdererMSP
        MSPDir: ./crypto-config/ordererOrganizations/example.com/msp
        AdminPrincipal: Role.ADMIN

    - &Org1
        Name: Org1MSP
        ID: Org1MSP
        MSPDir: ./crypto-config/peerOrganizations/org1.example.com/msp
        AdminPrincipal: Role.ADMIN



Orderer: &OrdererDefaults

    # Orderer Type: The orderer implementation to start
    # Available types are "solo" and "kafka"
    OrdererType: solo
    Addresses:
        - orderer.example.com:7050
    BatchTimeout: 500ms

    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 98 MB
        PreferredMaxBytes: 512 KB
    MaxChannels: 0
    Organizations:

        Application: &ApplicationDefaults
   # Organizations is the list of orgs which are defined as participants on
    # the application side of the network.
    Organizations:
   

################################################################################
#
#   Profile
#
#   - Different configuration profiles may be encoded here to be specified
#   as parameters to the configtxgen tool
#
################################################################################
Profiles:

    SampleOrg:
      Orderer:
        <<: *OrdererDefaults
        Organizations:
          - *OrdererOrg
      Consortiums:
        SampleConsortium:
          Organizations:
            - *OrdererOrg
            - *Org1
    
    SampleChannel:
      Consortium: SampleConsortium
      Application:
        Organizations:
          - *Org1
```

1.创建创世块

执行 configtxgen -profile SampleOrg -outputBlock ./channel-artifacts/gensis.block

2. 通道配置

configtxgen -profile SampleChannel -outputCreateChannelTx ./channel-artifacts/samplechannel.tx -channelID samplechannel

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
<<<<<<< HEAD
      - ORDERER_GENERAL_LOCALMSPID=example.com
=======
      - ORDERER_GENERAL_LOCALMSPID=org1.example.com
>>>>>>> basic
      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
      - ORDERER_GENERAL_TLS_ENABLED=true
      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyper/ledger/orderer/tls/ca.crt]
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: orderer
    volumes:
      - ./channel-artifacts/genesis.block:/var/hyperledger/orderer/genesis.block
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
      - CORE_PEER_LOCALMSPID=org1.example.com
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
      - CORE_PEER_LOCALMSPID=org1.example.com
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