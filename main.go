package main

import (
	"fmt"
	"os"

	"github.com/tramsyck/fabric-tutorial/blockchain"
	"github.com/tramsyck/fabric-tutorial/web"
	"github.com/tramsyck/fabric-tutorial/web/controllers"
)

func main() {
	fSetup := blockchain.FabricSetup{
		OrdererID:     "orderer.example.com",
		ChannelID:     "samplechannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/tramsyck/fabric-tutorial/fixtures/channel-artifacts/samplechannel.tx",

		ChainCodeID:     "samplechannel",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/tramsyck/fabric-tutorial/chaincode",
		OrgAdmin:        "Admin",
		OrgName:         "Org1",
		ConfigFile:      "config.yaml",
		UserName:        "User1",
	}

	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK:%v\n", err)
		return
	}

	defer fSetup.CloseSDK()

	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode:%v\n", err)
		return
	}

	app := &controllers.Application{
		Fabric: &fSetup,
	}

	web.Serve(app)
}
