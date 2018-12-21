package define

// Message 消息转发平台 结构
type Message struct {
	Factor            // 业务逻辑数据
	FabricTxId string `json:"fabricTxId"`
}

type CustomsDeclarationMessage struct {
	CustomsDeclarationInfo
	FabricTxId string `json:"fabricTxId"`
}

//AccessMessage
type AccessMessage struct {
	Access            // 业务逻辑数据
	FabricTxId string `json:"fabricTxId"`
}

type QueryContents struct {
	Schema  string      `json:"$schema"`
	Payload interface{} `json:"payload"`
}
type UserInfoMessage struct {
	UserInfo          // 业务逻辑数据
	FabricTxId string `json:"fabricTxId"`
}

type LogLevel struct {
	Level string `json:"level"`
}
