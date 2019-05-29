package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// HeroesServiceChaincode is xxx
type HeroesServiceChaincode struct {
}

// Init init chaincode
func (t *HeroesServiceChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("######### HeroesServiceChaincode Init ################")
	function, _ := stub.GetFunctionAndParameters()
	if function != "init" {
		return shim.Error("Unkown function call")
	}

	err := stub.PutState("hello", []byte("world"))
	if err != nil {
		return shim.Error(err.Error())

	}
	return shim.Success(nil)
}

// Invoke is jdkfjkdfj
func (t *HeroesServiceChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("################## Heroes Invoke###########")
	function, args := stub.GetFunctionAndParameters()
	if function != "invoke" {
		return shim.Error("unkown function call")
	}

	if len(args) < 1 {
		return shim.Error("the number of arguments is insuffient")

	}

	if args[0] == "query" {
		return t.query(stub, args)
	}

	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	return shim.Error("Unkown action, checkout the first argument")
}

func (t *HeroesServiceChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("############# HeroesServiceChainCode query ############")
	if len(args) < 2 {
		return shim.Error("The number of arguments is insuffficient.")
	}

	if args[1] == "hello" {
		state, err := stub.GetState("hello")
		if err != nil {
			return shim.Error("Failed to get state of hello")
		}

		return shim.Success(state)
	}

	return shim.Error("Unkown query action, check the second argument")
}

func (t *HeroesServiceChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("############# HeroesServiceChainCode invoke ############")
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient")
	}

	if args[1] == "hello" && len(args) == 3 {
		err := stub.PutState("hello", []byte(args[2]))
		if err != nil {
			return shim.Error("failed to update state of hello")
		}

		err = stub.SetEvent("eventInvoke", []byte{})
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success(nil)
	}

	return shim.Error("Unkown invoke action, check the second argument.")
}

func main() {
	err := shim.Start(new(HeroesServiceChaincode))
	if err != nil {
		fmt.Printf("Error starting Heroes Service chaincode：%s", err)
	}
}
