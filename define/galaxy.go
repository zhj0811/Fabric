package define

const(
	SaveRegistration = "SaveRegistration"
	QueryDataByFabricTxID = "QueryDataByFabricTxID"
	QueryRegistrationByNo = "QueryRegistrationByNo"
)


type Registration struct {
	Timestamp        string  `json:"timestamp"`        //时间戳
	Channel          string  `json:"channel"`          //渠道
	AdvanceOrderNo   string  `json:"advanceOrderNo"`   //预入库单编号
	WarehouseNo      string  `json:"warehouseNo"`      //仓库编号
	MaerchantNo      string  `json:"maerchantNo"`      //商户编号
	SubMaerchantNo   string  `json:"subMaerchantNo"`   //子商户编号
	LoanAmount       string  `json:"loanAmount"`       //放款额度
	ApprovalResults  string  `json:"approvalResults"`  //审批结果
	Approver         string  `json:"approver"`         //审批人
	ApprovalDate     string  `json:"approvalDate"`     //审批日期
	ApprovalComments string  `json:"approvalComments"` //审批意见
	OrderArrayList   []Order `json:"orderArrayList"`
}

type Order struct {
	Brand          	      string `json:"brand"`               //品牌
	Sku                   string `json:"sku"`                 //sku
	SkuName               string `json:"skuName"`             //sku名称
	PledgeRate            string `json:"pledgeRate"`          //质押率
	Merchandise           string `json:"merchandise"`         //商品规格
	Quantity              string `json:"quantity"`            //预入库数量
	PassQuantity          string `json:"passQuantity"`        //审核通过数量
	Price                 string `json:"price"`               //价格
	ProductionDate        string `json:"productionDate"`      //生产日期
	ShelfLife      	      string `json:"shelfLife,omitempty"` //保质期
	DueDate               string `json:"dueDate,omitempty"`   //到期日
	Attribute             string `json:"attribute,omitempty"` 	          //商品属性
	ProductCode    		  string `json:"productCode"`	                  //外包装大条码
	ProductIdentifierCode string `json:"productIdentifierCode,omitempty"` //外包装小条码
}