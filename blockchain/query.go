package blockchain

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// QueryHello query the chaincode to get the state of hello
func (setup *FabricSetup) QueryHello() (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "query")
	args = append(args, "hello")

	response, err := setup.client.Query(
		channel.Request{
			ChaincodeID: setup.ChainCodeID,
			Fcn:         args[0],
			Args:        [][]byte{[]byte(args[1]), []byte(args[2])}},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints("peer0.org1.example.com"))
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err)
	}

	return string(response.Payload), nil
}
