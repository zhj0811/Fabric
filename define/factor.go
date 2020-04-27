package define

//go:generate easyjson -all factor.go
// Factory struct
type Factory struct {
	Key     string `json:"key" binding:"required"`   //业务key
	Value   string `json:"value" binding:"required"` // 业务数据
	Expand1 string `json:"expand1,omitempty"`        // 扩展字段1
	Expand2 string `json:"expand2,omitempty"`        // 扩展地段2
}

type ResData struct {
	ResCode int    `json:"resCode"` //错误码0:成功1:失败
	ResMsg  string `json:"resMsg"`  //错误信息
}
