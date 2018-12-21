package utils

import (
	"encoding/json"
	"errors"

	"github.com/peersafe/tradetrain/define"
)

// FormatRequestMessage format requset to message json
func FormatRequestMessage(request define.Factor) ([]byte, error) {
	var invokeRequest define.InvokeRequest
	var err error
	message := &define.Message{}
	if request.Key == "" {
		err = errors.New("Key is NULL")
		return nil, err
	}
	message.Key = request.Key
	if request.BusinessType == "" {
		err = errors.New("BusinessType is NULL")
		return nil, err
	}
	message.BusinessType = request.BusinessType
	if request.DataType == "" {
		err = errors.New("DataType is NULL")
		return nil, err
	}
	message.DataType = request.DataType
	if request.WriteRoleType == "" {
		err = errors.New("WriteRoleType is NULL")
		return nil, err
	}
	message.WriteRoleType = request.WriteRoleType
	if request.Writer == "" {
		err = errors.New("Writer is NULL")
		return nil, err
	}
	message.Writer = request.Writer
	message.Version = request.Version
	message.BusinessData = request.BusinessData

	message.Expand1 = request.Expand1
	message.Expand2 = request.Expand2

	b, _ := json.Marshal(message)
	invokeRequest.Value = string(b)
	invokeRequest.Key = request.Key

	return json.Marshal(invokeRequest)
}

// FormatRequestFormMessage format requset to message json
func FormatRequestFormMessage(request define.CustomsDeclarationInfo) ([]byte, error) {
	var invokeRequest define.InvokeRequest
	message := &define.CustomsDeclarationMessage{}
	message.Key = request.Key
	message.BusinessType = request.BusinessType
	message.DataType = request.DataType
	message.WriteRoleType = request.WriteRoleType
	message.Writer = request.Writer
	message.EntryID = request.EntryID
	message.Version = request.Version
	message.BusinessData = request.BusinessData

	message.Expand1 = request.Expand1
	message.Expand2 = request.Expand2
	message.Version = request.Version

	b, _ := json.Marshal(message)
	invokeRequest.Value = string(b)
	invokeRequest.Key = request.Key

	return json.Marshal(invokeRequest)
}

func FormatRequestAccessMessage(request define.Access) ([]byte, error) {
	//CryptoAlgorithm is used by hoperun, we don't need to care about it. By wuxu, 20170901.
	//request.CryptoAlgorithm = "aes"
	var invokeRequest define.InvokeRequest
	var err error
	accessMessage := &define.AccessMessage{}
	if request.Key == "" {
		err = errors.New("Key is NULL")
		return nil, err
	}
	accessMessage.Key = request.Key
	if request.BusinessType == "" {
		err = errors.New("BusinessType is NULL")
		return nil, err
	}
	accessMessage.BusinessType = request.BusinessType
	if request.DataType == "" {
		err = errors.New("DataType is NULL")
		return nil, err
	}
	accessMessage.DataType = request.DataType
	if request.Writer == "" {
		err = errors.New("Writer is NULL")
		return nil, err
	}
	accessMessage.Writer = request.Writer
	accessMessage.Version = request.Version
	accessMessage.ReaderList = request.ReaderList
	arrLen := len(request.ReaderList)
	if arrLen == 0 {
		err = errors.New("ReaderList is NULL")
		return nil, err
	}
	b, _ := json.Marshal(accessMessage)
	invokeRequest.Value = string(b)
	invokeRequest.Key = request.Key

	return json.Marshal(invokeRequest)
}

func VerifyQueryDataRequestFormat(request *define.QueryData) error {
	var err error
	if request.Key == "" {
		err = errors.New("Key is NULL")
		return err
	}
	if request.BusinessType == "" {
		err = errors.New("BusinessType is NULL")
		return err
	}
	if request.DataType == "" {
		err = errors.New("DataType is NULL")
		return err
	}

	if request.Reader == "" {
		err = errors.New("Reader is NULL")
		return err
	}
	if request.WriteRoleType == "" {
		err = errors.New("WriteRoleType is NULL")
		return err
	}

	return nil
}

func VerifyQueryAclRequestFormat(request *define.QueryACL) error {
	var err error
	if request.Key == "" {
		err = errors.New("Key is NULL")
		return err
	}
	if request.BusinessType == "" {
		err = errors.New("BusinessType is NULL")
		return err
	}
	if request.DataType == "" {
		err = errors.New("DataType is NULL")
		return err
	}

	if request.Writer == "" {
		err = errors.New("Writer is NULL")
		return err
	}

	return nil
}

// FormatResponseFormMessage format response to message json
func FormatResponseFormMessage(userId string, request *define.CustomsDeclarationInfo, messages *define.CustomsDeclarationMessage) error {
	*request = messages.CustomsDeclarationInfo
	return nil
}

// FormatResponseMessage format response to message json
func FormatResponseMessage(request *define.Factor, messages *define.Message) error {
	*request = messages.Factor
	return nil
}

func FormatResponseAccessMessage(request *define.Access, messages *define.AccessMessage) error {
	*request = messages.Access
	return nil
}

func FormatUserInfoRequestMessage(request define.UserInfo) ([]byte, error) {
	var invokeRequest define.InvokeRequest
	var err error
	message := &define.UserInfoMessage{}
	if request.Key == "" {
		err = errors.New("Key is NULL")
		return nil, err
	}
	message.Key = request.Key
	if request.BusinessType == "" {
		err = errors.New("BusinessType is NULL")
		return nil, err
	}
	message.BusinessType = request.BusinessType
	if request.DataType == "" {
		err = errors.New("DataType is NULL")
		return nil, err
	}
	message.DataType = request.DataType
	if request.WriteRoleType == "" {
		err = errors.New("WriteRoleType is NULL")
		return nil, err
	}
	message.WriteRoleType = request.WriteRoleType
	if request.Writer == "" {
		err = errors.New("Writer is NULL")
		return nil, err
	}
	message.Writer = request.Writer
	message.Version = request.Version
	if request.UserName == "" {
		err = errors.New("UserName is NULL")
		return nil, err
	}
	message.UserName = request.UserName
	if request.UserID == "" {
		err = errors.New("UserID is NULL")
		return nil, err
	}
	message.UserID = request.UserID
	if request.UserType == "" {
		err = errors.New("UserType is NULL")
		return nil, err
	}
	message.UserType = request.UserType
	if request.UserArea == "" {
		err = errors.New("UserArea is NULL")
		return nil, err
	}
	message.UserArea = request.UserArea

	b, _ := json.Marshal(message)
	invokeRequest.Value = string(b)
	invokeRequest.Key = request.Key

	return json.Marshal(invokeRequest)
}

func VerifyQueryUserInfoRequestFormat(request *define.QueryUserInfo) error {
	var err error
	if request.Key == "" {
		err = errors.New("Key is NULL")
		return err
	}
	if request.BusinessType == "" {
		err = errors.New("BusinessType is NULL")
		return err
	}
	if request.DataType == "" {
		err = errors.New("DataType is NULL")
		return err
	}

	if request.WriteRoleType == "" {
		err = errors.New("WriteRoleType is NULL")
		return err
	}

	return nil
}
