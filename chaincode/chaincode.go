package main

import (
	"github.com/peersafe/factoring/chaincode/handler"
	"github.com/peersafe/factoring/define"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
)

var logger = shim.NewLogger("factorChaincode")

type handlerFunc func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error)

var funcHandler = map[string]handlerFunc{
	define.SaveData:              handler.SaveData,
	define.KeepaliveQuery:        handler.KeepaliveQuery,
	define.QueryDataByFabricTxId: handler.QueryDataByFabricTxId,
	define.QueryDataByBusinessNo: handler.QueryDataByBusinessNo,
}

type FactorChaincode struct {
}

func init() {
	logger.SetLevel(shim.LogDebug)
}

// Init method will be called during deployment.
// The deploy transaction metadata is supposed to contain the administrator cert
func (t *FactorChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init Chaincode...")

	err := stub.PutState(handler.KEEPALIVETEST, []byte(handler.KEEPALIVETEST))
	if err != nil {
		logger.Error("Init Chaincode error:%s", err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}

func (t *FactorChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Debugf("Invoke function=%v,args=%v\n", function, args)
    logger.Infof("Invoke to chaincode txId is : %s.",stub.GetTxID())
 	logger.Debugf("Invoke to chaincode time is %d .",time.Now().UnixNano()/1000000)

	if len(args) < 1 || len(args[0]) == 0 {
		logger.Error("the invoke args not exist or arg[0] is empty")
		return shim.Error("the invoke args not exist  or arg[0] is empty")
	}

	currentFunc := funcHandler[function]
	if currentFunc == nil {
		logger.Error("the function name not exist!!")
		return shim.Error("the function name not exist!!")
	}

	payload, err := currentFunc(stub, function, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(payload)
}

func main() {
	err := shim.Start(new(FactorChaincode))
	if err != nil {
		logger.Error("Error starting Factorchaincode: %s", err)
	}
}
