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