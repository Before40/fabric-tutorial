package blockchain

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

// InvokeHello invoke hello method
func (setup *FabricSetup) InvokeHello(value string) (string, error) {
	var args []string
	args = append(args, "invoke")
	args = append(args, "invoke")
	args = append(args, "hello")
	args = append(args, value)
	eventID := "eventInvoke"

	transientDataMap := make(map[string][]byte)

	transientDataMap["result"] = []byte("Transient data in hello invoke")
	reg, notifier, err :=
		setup.event.RegisterChaincodeEvent(setup.ChainCodeID, eventID)
	if err != nil {
		return "", err
	}

	defer setup.event.Unregister(reg)

	response, err := setup.client.Execute(
		channel.Request{
			ChaincodeID:  setup.ChainCodeID,
			Fcn:          args[0],
			Args:         [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])},
			TransientMap: transientDataMap},
		channel.WithTargetEndpoints("peer0.org1.example.com;peer1.org1.example.com"))
	if err != nil {
		return "", fmt.Errorf("failed to move funds: %v", err)
	}

	select {
	case ccEvent := <-notifier:
		fmt.Printf("received cc event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	return string(response.TransactionID), nil
}
