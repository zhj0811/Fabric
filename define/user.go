package define

type UserBaseInfo struct {
	Key           string `json:"key"`           // 业务key
	BusinessType  string `json:"businessType"`  // 业务类型 默认Import0001
	DataType      string `json:"dataType"`      // 业务数据类型
	WriteRoleType string `json:"writeRoleType"` // 写入数据角色类型
	Writer        string `json:"writer"`        // 写入人id
	Version       string `json:"version"`       // 数据版本号
}
type UserInfo struct {
	UserBaseInfo
	UserName string `json:"userName"` //用户名称
	UserID   string `json:"userID"`   //用户id
	UserType string `json:"userType"` //用户类型
	UserArea string `json:"userArea"` //用户区域
}

type SaveUserInfoRespone struct {
	UserBaseInfo
	ResponseCode    string `json:"responsecode"`    //返回码
	ResponseExplain string `json:"responseexplain"` //返回码信息
}

type QueryUserInfo struct {
	Key           string `json:"key"`           // 业务key
	BusinessType  string `json:"businessType"`  // 业务类型 默认Import0001
	DataType      string `json:"dataType"`      // 业务数据类型
	WriteRoleType string `json:"writeRoleType"` // 写入数据角色类型
	Reader        string `json:"reader"`        //读取人
}

type QueryUserInfoRespone struct {
	UserInfo
	ResponseCode    string `json:"responsecode"`    //返回码
	ResponseExplain string `json:"responseexplain"` //返回码信息
}
