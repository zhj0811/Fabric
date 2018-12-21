package define

type FactorBaseInfo struct {
	Key           string `json:"key"`           // 业务key
	BusinessType  string `json:"businessType"`  // 业务类型 默认Import0001
	DataType      string `json:"dataType"`      // 业务数据类型
	WriteRoleType string `json:"writeRoleType"` // 写入数据角色类型
	Writer        string `json:"writer"`        // 写入人id
	Version       string `json:"version"`       // 数据版本号
}
type Factor struct {
	FactorBaseInfo
	BusinessData string `json:"businessData"` // 业务数据
	Expand1      string `json:"expand1"`      // 扩展字段1
	Expand2      string `json:"expand2"`      // 扩展地段2
}

type CustomsDeclarationInfo struct {
	FactorBaseInfo
	EntryID      string `json:"EntryID"`      // 报关单号
	BusinessData string `json:"businessData"` // 业务数据
	Expand1      string `json:"expand1"`      // 扩展字段1
	Expand2      string `json:"expand2"`      // 扩展地段2
}

type FactorResponse struct {
	FactorBaseInfo
	ResponseCode    string `json:"responseCode"`    //返回码
	ResponseExplain string `json:"responseExplain"` //返回说明
}

type QueryData struct {
	Key           string `json:"key"`           // 业务key
	BusinessType  string `json:"businessType"`  // 业务类型 默认Import0001
	DataType      string `json:"dataType"`      // 业务数据类型
	WriteRoleType string `json:"writeRoleType"` // 写入数据角色类型
	Reader        string `json:"reader"`        //读取人
}

type QueryFormData struct {
	Key           string `json:"key"`           // 业务key
	BusinessType  string `json:"businessType"`  // 业务类型 默认Import0001
	DataType      string `json:"dataType"`      // 业务数据类型
	WriteRoleType string `json:"writeRoleType"` // 写入数据角色类型
	Reader        string `json:"reader"`        //读取人
	EntryID       string `json:"entryID"`       // 报关单号
}

type QueryDataResponse struct {
	FactorBaseInfo
	BusinessData    string `json:"businessData"`    // 业务数据
	ResponseCode    string `json:"responseCode"`    //返回码
	ResponseExplain string `json:"responseExplain"` //返回说明
}

type QueryFormDataResponse struct {
	FactorBaseInfo
	EntryID         string `json:"entryID"`         // 报关单号
	BusinessData    string `json:"businessData"`    // 业务数据
	ResponseCode    string `json:"responseCode"`    //返回码
	ResponseExplain string `json:"responseExplain"` //返回说明
}

type FileInfo struct {
	Name string `json:"name"` //文件名称
	Hash string `json:"hash"` //文件hash
	Path string `json:"path"` //文件下载地址
}

// BlockchainData 区块信息模型
type BlockchainData struct {
	TxId           string      `json:"txId"`        // 交易ID
	TxHash         string      `json:"txHash"`      // 交易请求hash
	BlockHash      string      `json:"blockHash"`   // 当前区块hash
	BlockHeight    uint64      `json:"blockHeight"` // 当前区块高度
	Bidbond        string      `json:"bidbond"`     // 投标保函唯一编号
	Bid            string      `json:"bid"`         // 开立投标保函对应招标需求唯一编号
	Progress       string      `json:"progress"`    // 进度描述信息
	CreateBy       string      `json:"createBy"`
	CreateTime     uint64      `json:"createTime"`
	Sender         string      `json:"sender"`
	Receiver       []string    `json:"receiver"`
	LastUpdateTime uint64      `json:"lastUpdateTime"`
	LastUpdateBy   string      `json:"lastUpdateBy"`
	BlockData      string      `json:"blockData"`
	Remark         string      `json:"remark"`
	Status         StateEntity `json:"status"` // 状态
}

// StateEntity  状态存储结构(json)
type StateEntity struct {
	ChangeEvent string `json:"changeEvent"` // 变更事件
	PreState    string `json:"preState"`    // 上一状态如"cancel"
	CurrState   string `json:"currState"`   // 当前状态如"applied"
}

type Events struct {
	ChaincodeId string      `json:"chaincodeId"` //链码ID
	TxId        string      `json:"txId"`        //交易ID
	EventName   string      `json:"eventName"`   //事件名
	Payload     interface{} `json:"payload"`
}

// BlockchainData 区块信息模型
type BlockDataObj struct {
	BlockHash    string      `json:"blockHash"`    //当前区块hash
	BlockHeight  uint64      `json:"blockHeight"`  //当前区块高度
	PreviousHash string      `json:"previousHash"` //前一个区块哈希
	Events       interface{} `json:"events"`       //生成的事件结构
}
