package define

type ACLBaseInfo struct {
	Key          string `json:"key"`          // 业务key
	BusinessType string `json:"businessType"` // 业务类型 默认Import0001
	DataType     string `json:"dataType"`     // 业务数据类型
	Writer       string `json:"writer"`       // 写入人id
	Version      string `json:"version"`      // 数据版本号
}
type Access struct {
	ACLBaseInfo
	ReaderList []string `json:"readerList"` //权限访问列表
}
type AccessResponse struct {
	ACLBaseInfo
	ResponseCode    string `json:"responseCode"`    //返回码
	ResponseExplain string `json:"responseExplain"` //返回说明
}

type QueryACL struct {
	ACLBaseInfo
}

type QueryACLResponse struct {
	Access
	ResponseCode    string `json:"responseCode"`    //返回码
	ResponseExplain string `json:"responseExplain"` //返回说明
}
