package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

// FabricSetup impl
type FabricSetup struct {
	ConfigFile      string
	OrgID           string
	OrdererID       string
	ChannelID       string
	ChainCodeID     string
	initialized     bool
	ChannelConfig   string
	ChaincodeGoPath string
	ChaincodePath   string
	OrgAdmin        string
	OrgName         string
	UserName        string
	client          *channel.Client
	admin           *resmgmt.Client
	sdk             *fabsdk.FabricSDK
	event           *event.Client
}

// Initialize 读取
func (setup *FabricSetup) Initialize() error {
	if setup.initialized {
		return errors.New("sd already initialized")
	}

	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create sdk")
	}

	setup.sdk = sdk
	fmt.Println("SDK created")

	resourceManagerClientContext := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))

	if err != nil {
		return errors.WithMessage(err, "failed to load Admin identify")
	}

	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}

	setup.admin = resMgmtClient
	fmt.Println("Resouce management client created")

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		return errors.WithMessage(err, "failed to create msp client")
	}

	adminIdentify, err := mspClient.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		return errors.WithMessage(err, "failed to get admin signing identity")
	}

	req := resmgmt.SaveChannelRequest{ChannelID: setup.ChannelID,
		ChannelConfigPath: setup.ChannelConfig,
		SigningIdentities: []msp.SigningIdentity{adminIdentify}}
	txID, err := setup.admin.SaveChannel(req, resmgmt.WithOrdererEndpoint(setup.OrdererID))
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}

	fmt.Println("Channel created")

	if err = setup.admin.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
		return errors.WithMessage(err, "failed to make admin join channel")
	}

	fmt.Println("Channel joined")
	fmt.Println("Initialization successful")
	setup.initialized = true
	return nil
}

// InstallAndInstantiateCC install and initantiate chaincode
func (setup *FabricSetup) InstallAndInstantiateCC() error {
	ccPkg, err := gopackager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("ccPkg created")

	installCCReq := resmgmt.InstallCCRequest{
		Name:    setup.ChainCodeID,
		Path:    setup.ChaincodePath,
		Version: "0",
		Package: ccPkg}
	_, err = setup.admin.InstallCC(
		installCCReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")

	}
	fmt.Println("Chaincode installed")

	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	resp, err := setup.admin.InstantiateCC(
		setup.ChannelID,
		resmgmt.InstantiateCCRequest{
			Name:    setup.ChainCodeID,
			Path:    setup.ChaincodeGoPath,
			Version: "0",
			Args:    [][]byte{[]byte("init")},
			Policy:  ccPolicy})

	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincde")
	}

	fmt.Println("Chaincode instantiated")

	clientContext := setup.sdk.ChannelContext(setup.ChannelID,
		fabsdk.WithUser(setup.UserName))
	setup.client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	setup.event, err = event.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")
	fmt.Println("Chaincode Intallation & Instantiation Successful")
	return nil
}

// CloseSDK close the sdk
func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close()
}
