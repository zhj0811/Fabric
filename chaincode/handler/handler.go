package handler

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/zhj0811/fabric/chaincode/utils"
	"github.com/zhj0811/fabric/define"
)

var KEEPALIVETEST = "test"
var KEYTEST = "keytest"
var ACLTEST = "acltest"
var USERINFO = "userinfo"

var myLogger = shim.NewLogger("handler")

func init() {
	myLogger.SetLevel(shim.LogDebug)
}

func SaveData(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	//Save (request.Key,request.Value) as k,v
	request := &define.InvokeRequest{}
	if err = json.Unmarshal([]byte(args[0]), request); err != nil {
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	//request.Key = stub.GetTxID() //txId

	message := &define.Message{}
	err = json.Unmarshal([]byte(request.Value), message)
	if err != nil {
		myLogger.Errorf("Failed to unmarshal request.Value:%s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	message.FabricTxId = stub.GetTxID() //txid

	myLogger.Debugf("message.BusinessType:%s", message.BusinessType)
	myLogger.Debugf("message.Key:%s", message.Key)

	//Save (CreateCompositeKey,request.Value) as k,v
	var attributes []string
	attributes = append(attributes, message.Key, message.BusinessType, message.DataType, message.WriteRoleType)
	key, err := stub.CreateCompositeKey(KEYTEST, attributes)
	if err != nil {
		myLogger.Errorf("Failed to CreateCompositeKey: %s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	myLogger.Debugf("CreateCompositeKey is: ", key)
	value, _ := json.Marshal(message)
	err = stub.PutState(key, value)
	if err != nil {
		myLogger.Errorf("saveData err: %s", err.Error())
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	return utils.InvokeResponse(stub, nil, function, []string{request.Value}, false)
}

func SaveACL(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	//Save (request.Key,request.Value) as k,v
	request := &define.InvokeRequest{}
	if err = json.Unmarshal([]byte(args[0]), request); err != nil {
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	message := &define.AccessMessage{}
	err = json.Unmarshal([]byte(request.Value), message)
	if err != nil {
		myLogger.Errorf("Failed to unmarshal request.Value:%s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	message.FabricTxId = stub.GetTxID() //txid

	myLogger.Debugf("message.BusinessType:%s", message.BusinessType)
	myLogger.Debugf("message.Key:%s", message.Key)

	//Save (CreateCompositeKey,request.Value) as k,v
	var attributes []string
	attributes = append(attributes, message.Key, message.BusinessType, message.DataType, message.Writer)
	key, err := stub.CreateCompositeKey(ACLTEST, attributes)
	if err != nil {
		myLogger.Errorf("Failed to CreateCompositeKey: %s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	myLogger.Debugf("CreateCompositeKey is: %s", key)
	value, _ := json.Marshal(message)
	err = stub.PutState(key, value)
	if err != nil {
		myLogger.Errorf("saveData err: %s", err.Error())
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	return utils.InvokeResponse(stub, nil, function, []string{request.Value}, false)
}

func QueryDataByKey(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	myLogger.Debugf("QueryDataByKey....")
	var attributes []string
	attributes = append(attributes, args[0], args[1], args[2], args[3])
	reader := args[4]
	compositekey, _ := stub.CreateCompositeKey(KEYTEST, attributes)
	targetValue, _ := stub.GetState(compositekey)
	message := &define.Message{}
	if len(targetValue) == 0 {
		err = errors.New(define.ValueOfKeyNil)
		return utils.QueryResponse(err, nil, define.Page{})
	}
	_ = json.Unmarshal(targetValue, message)
	writer := message.Writer
	response := string(targetValue)
	if writer == reader {
		return utils.QueryResponse(nil, response, define.Page{})
	}

	var Keyattributes []string
	Keyattributes = append(Keyattributes, message.Key, message.BusinessType, message.DataType, message.Writer)
	AclKey, err := stub.CreateCompositeKey(ACLTEST, Keyattributes)
	if err != nil {
		myLogger.Errorf("Failed to CreateCompositeKey: %s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	aclValue, err := stub.GetState(AclKey)
	if len(aclValue) == 0 {
		err = errors.New(define.PermissionNotFound)
		return utils.QueryResponse(err, nil, define.Page{})
	}

	accessmessage := &define.AccessMessage{}
	_ = json.Unmarshal(aclValue, accessmessage)
	for _, user := range accessmessage.ReaderList {
		if reader == user {
			return utils.QueryResponse(nil, response, define.Page{})
		}
	}
	err = errors.New(define.NoPermission)
	return utils.QueryResponse(err, nil, define.Page{})
}

func QueryListById(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	var response string
	var attributes []string
	Key := args[0]
	BusinessType := args[1]
	DataType := args[2]
	Writer := args[3]
	attributes = append(attributes, Key, BusinessType, DataType, Writer)
	compositekey, err := stub.CreateCompositeKey(ACLTEST, attributes)
	if err != nil {
		myLogger.Errorf("Failed to CreateCompositeKey: %s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	myLogger.Errorf("CreateCompositeKey is:%s ", compositekey)

	targetValue, err := stub.GetState(compositekey)
	if err != nil {
		myLogger.Errorf("Failed to GetState() as key %s: %s", args[0], err)
		utils.QueryResponse(nil, nil, define.Page{})
	}
	if len(targetValue) == 0 {
		myLogger.Errorf("targetValue is:%d ", len(targetValue))
		err = errors.New(define.ValueOfKeyNil)
		return utils.QueryResponse(err, nil, define.Page{})
	}
	response = string(targetValue)
	return utils.QueryResponse(nil, response, define.Page{})
}

func KeepaliveQuery(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	targetValue, err := stub.GetState(KEEPALIVETEST)
	if err != nil {
		err = fmt.Errorf("ERROR! KeepaliveQuery get failed, err = %s", err.Error())
		return []byte("UnReached"), err
	}

	if string(targetValue) != KEEPALIVETEST {
		err = fmt.Errorf("ERROR! KeepaliveQuery get result is %s", string(targetValue))
		return []byte("UnReached"), err
	}

	return []byte("Reached"), nil
}

func SaveUserInfo(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error
	//Save (request.Key,request.Value) as k,v
	request := &define.InvokeRequest{}
	if err = json.Unmarshal([]byte(args[0]), request); err != nil {
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	message := &define.UserInfoMessage{}
	err = json.Unmarshal([]byte(request.Value), message)
	if err != nil {
		myLogger.Errorf("Failed to unmarshal request.Value:%s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	message.FabricTxId = stub.GetTxID() //txid

	myLogger.Debugf("message.BusinessType:%s", message.BusinessType)
	myLogger.Debugf("message.Key:%s", message.Key)
	//Save (CreateCompositeKey,request.Value) as k,v
	var attributes []string
	attributes = append(attributes, message.Key, message.BusinessType, message.DataType, message.WriteRoleType)
	key, err := stub.CreateCompositeKey(KEYTEST, attributes)
	if err != nil {
		myLogger.Errorf("Failed to CreateCompositeKey: %s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	myLogger.Errorf("SaveUserInfo CreateCompositeKey is: %s", key)
	value, _ := json.Marshal(message)
	err = stub.PutState(key, value)
	if err != nil {
		myLogger.Errorf("saveData err: %s", err.Error())
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	return utils.InvokeResponse(stub, nil, function, []string{request.Value}, false)
}

func QueryUserInfo(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	myLogger.Debugf("QueryACL is: %s\n", args[0])
	request := &define.QueryUserInfo{}
	if err = json.Unmarshal([]byte(args[0]), request); err != nil {
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	var attributes []string
	attributes = append(attributes, request.Key, request.BusinessType, request.DataType, request.WriteRoleType)
	compositekey, err := stub.CreateCompositeKey(KEYTEST, attributes)
	if err != nil {
		myLogger.Errorf("Failed to CreateCompositeKey: %s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	myLogger.Errorf("CreateCompositeKey is:%s ", compositekey)

	targetValue, err := stub.GetState(compositekey)
	if err != nil {
		myLogger.Errorf("Failed to GetState() as key %s: %s", args[0], err)
		utils.QueryResponse(nil, nil, define.Page{})
	}

	if len(targetValue) == 0 {
		myLogger.Errorf("targetValue is:%d ", len(targetValue))
		err = errors.New(define.ValueOfKeyNil)
		return utils.QueryResponse(err, nil, define.Page{})
	}

	myLogger.Errorf("targetValue is:%d ", len(targetValue))
	response := string(targetValue)
	return utils.QueryResponse(nil, response, define.Page{})
}

// SaveTable 报关单表头、表体信息
func SaveTable(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	//Save (request.Key,request.Value) as k,v
	request := &define.InvokeRequest{}
	if err = json.Unmarshal([]byte(args[0]), request); err != nil {
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	//request.Key = stub.GetTxID() //txId

	message := &define.CustomsDeclarationMessage{}
	err = json.Unmarshal([]byte(request.Value), message)
	if err != nil {
		myLogger.Errorf("Failed to unmarshal request.Value:%s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	message.FabricTxId = stub.GetTxID() //txid

	myLogger.Debugf("message.BusinessType:%s", message.BusinessType)
	myLogger.Debugf("message.Key:%s", message.Key)

	//Save (CreateCompositeKey,request.Value) as k,v
	var attributes []string
	attributes = append(attributes, message.Key, message.BusinessType, message.DataType, message.WriteRoleType, message.EntryID)
	key, err := stub.CreateCompositeKey(KEYTEST, attributes)
	if err != nil {
		myLogger.Errorf("Failed to CreateCompositeKey: %s", err)
		return utils.InvokeResponse(stub, err, function, nil, false)
	}
	myLogger.Debugf("CreateCompositeKey is: ", key)

	err = stub.PutState(key, []byte(request.Value))
	if err != nil {
		myLogger.Errorf("saveData err: %s", err.Error())
		return utils.InvokeResponse(stub, err, function, nil, false)
	}

	return utils.InvokeResponse(stub, nil, function, []string{request.Value}, false)
}
