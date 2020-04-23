package main

import (
	"github.com/zhj0811/fabric/chaincode/handler"
	"github.com/zhj0811/fabric/define"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("chaincode")

type handlerFunc func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error)

var funcHandler = map[string]handlerFunc{
	define.SaveData:       handler.SaveData,
	define.SaveACL:        handler.SaveACL,
	define.KeepaliveQuery: handler.KeepaliveQuery,
	define.QueryDataByKey: handler.QueryDataByKey,
	define.QueryListById:  handler.QueryListById,
	define.SaveUserInfo:   handler.SaveUserInfo,
	define.QueryUserdata:  handler.QueryUserInfo,
	define.SaveTable:      handler.SaveTable, // 报关单表头、表体信息
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
