package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"github.com/zhj0811/fabric/define"
	mylogger "github.com/zhj0811/fabric/pkg/logger"
)

const (
	KEEPALIVE = "keepalive"
)

var logger = mylogger.NewSugaredLogger("DEBUG", "")

type handlerFunc func(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error)

var funcHandler = map[string]handlerFunc{
	define.SaveData:       SaveData,
	define.QueryData:      QueryData,
	define.KeepaliveQuery: KeepaliveQuery,
}

type FactoryChaincode struct {
}

// Init method will be called during deployment.
// The deploy transaction metadata is supposed to contain the administrator cert
func (t *FactoryChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debug("Init Chaincode...")
	err := stub.PutState(KEEPALIVE, []byte(KEEPALIVE))
	if err != nil {
		err = errors.WithMessage(err, "init chaincode error")
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("SUCCESS"))
}

func (t *FactoryChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	//logger.Debug("Invoke function=%v,args=%v\n", function, args)
	logger.Debugf("Invoke function=%s,args=%s", function, args)
	if len(args) < 1 || len(args[0]) == 0 {
		err := errors.New("the invoke args not exist or arg[0] is empty")
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}

	currentFunc := funcHandler[function]
	if currentFunc == nil {
		err := errors.New("the function name not exist!!")
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}

	payload, err := currentFunc(stub, function, args)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(payload)
}

func SaveData(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	//Save (request.Key,request.Value) as k,v
	logger.Debugf("Entering %s....", function)
	request := &define.Factory{}
	if err := request.UnmarshalJSON([]byte(args[0])); err != nil {
		return nil, errors.WithMessagef(err, "invoke func %s error for Unmarshal request", function)
	}

	//request.Key = stub.GetTxID() //txId
	logger.Debugf("request.Key:%s", request.Key)

	//Save (CreateCompositeKey,request.Value) as k,v
	//var attributes []string
	//attributes = append(attributes, message.Key, message.BusinessType, message.DataType, message.WriteRoleType)
	//key, err := stub.CreateCompositeKey(KEYTEST, attributes)

	value, err := request.MarshalJSON()
	if err != nil {
		return nil, errors.WithMessagef(err, "invoke func %s error for marshal value", function)
	}
	err = stub.PutState(request.Key, value)
	if err != nil {
		return nil, errors.WithMessagef(err, "inoke func %s err for PutState", function)
	}
	return nil, nil
}

func QueryData(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debugf("Entering %s....", function)
	//targetValue, _ := stub.GetState(compositekey)

	state, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.WithMessagef(err, "invoke func %s error with GetState", function)
	}
	return state, nil
}

func KeepaliveQuery(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	logger.Debugf("Entering %s....", function)
	targetValue, err := stub.GetState(KEEPALIVE)
	if err != nil {
		return nil, errors.WithMessagef(err, "invoke func %s failed with GetState", function)
	}
	if string(targetValue) != KEEPALIVE {
		return nil, errors.Errorf("%s get result is %s not %s", function, string(targetValue), KEEPALIVE)
	}
	return []byte(KEEPALIVE), nil
}

func main() {
	err := shim.Start(new(FactoryChaincode))
	if err != nil {
		logger.Errorf("Error starting chaincode: %s", err)
	}
}
