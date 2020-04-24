package define

//go:generate easyjson -all factor.go
// Factory struct
type Factory struct {
	Key     string `json:"key"`     //业务key
	Value   string `json:"value"`   // 业务数据
	Expand1 string `json:"expand1"` // 扩展字段1
	Expand2 string `json:"expand2"` // 扩展地段2
}
