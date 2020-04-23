package router

import (
	"sync"

	"github.com/zhj0811/fabric/apiserver/handler"

	"github.com/gin-gonic/gin"
)

// Router 全局路由
var router *gin.Engine
var onceCreateRouter sync.Once

func GetRouter() *gin.Engine {
	onceCreateRouter.Do(func() {
		router = createRouter()
	})

	return router
}

func createRouter() *gin.Engine {
	router := gin.Default()

	traders := router.Group("/traders") // 贸易商——天物
	{
		traders.POST("/saveImportContract", handler.SaveData)        // 1 外贸合同
		traders.POST("/queryImportContract", handler.QueryDataByKey) // 外贸合同查询

		traders.POST("/saveLCInfo", handler.SaveData)        // 2 信用证信息
		traders.POST("/queryLCInfo", handler.QueryDataByKey) // 信用证信息查询

		traders.POST("/saveArriveNotice", handler.SaveData)        // 6 到港信息
		traders.POST("/queryArriveNotice", handler.QueryDataByKey) // 到港信息查询

		traders.POST("/saveInCustomsDelegations", handler.SaveData)        // 7 入区报关委托协议
		traders.POST("/queryInCustomsDelegations", handler.QueryDataByKey) // 入区报关委托协议查询

		// traders.POST("/saveOutCustomsDelegations", handler.SaveData)        // 8 出区报关委托协议
		// traders.POST("/queryOutCustomsDelegations", handler.QueryDataByKey) // 出区报关委托协议查询

		// traders.POST("/saveVehicleSale", handler.SaveData)        // 10 车辆销售
		// traders.POST("/queryVehicleSale", handler.QueryDataByKey) // 车辆销售查询

		// traders.POST("/saveInOutInvoice", handler.SaveData)        // 11 入出区发票
		// traders.POST("/queryInOutInvoice", handler.QueryDataByKey) // 入出区发票查询
	}

	foreign := router.Group("/foreign") // 海外端——车企
	{
		// foreign.POST("/saveGoodsCategory", handler.SaveData)        // 1 商品分类要素
		foreign.POST("/queryGoodsCategory", handler.QueryDataByKey) // 商品分类要素查询

		// foreign.POST("/saveAutoCategory", handler.SaveData)        // 2 车辆分类要素
		// foreign.POST("/queryAutoCategory", handler.QueryDataByKey) // 车辆分类要素查询

		// foreign.POST("/saveVehiclePrices", handler.SaveData)        // 3 车辆海外报价
		// foreign.POST("/queryVehiclePrices", handler.QueryDataByKey) // 车辆海外报价查询

		// foreign.POST("/saveProviderInvoice", handler.SaveData)        // 4 采购发票要素
		// foreign.POST("/queryProviderInvoice", handler.QueryDataByKey) // 采购发票要素查询

		// foreign.POST("/saveExportContract", handler.SaveData)        // 5 海外外贸合同要素
		// foreign.POST("/queryExportContract", handler.QueryDataByKey) // 海外外贸合同要素查询

		// foreign.POST("/saveExportInvoice", handler.SaveData)        // 6 海外外贸发票要素
		// foreign.POST("/queryExportInvoice", handler.QueryDataByKey) // 海外外贸发票要素查询

		// foreign.POST("/saveOverseasClearance", handler.SaveData)        // 7 海外通关要素信息
		// foreign.POST("/queryOverseasClearance", handler.QueryDataByKey) // 海外通关要素信息查询

		// foreign.POST("/saveDispatchInfo", handler.SaveData)  // 8 海外发货信息
		// foreign.POST("/queryDispatchInfo", handler.SaveData) // 海外发货信息查询

		// foreign.POST("/savePremiumsFreight ", handler.SaveData)       // 9 运保费
		// foreign.POST("/queryPremiumsFreight", handler.QueryDataByKey) // 运保费查询
	}

	customs := router.Group("/customs") // 海关
	{
		customs.POST("/saveStatusOfCustoms", handler.SaveData)        // 1 报关单工作流
		customs.POST("/queryStatusOfCustoms", handler.QueryDataByKey) // 报关单工作流查询

		customs.POST("/saveCheckResult", handler.SaveData)        // 2 审核反馈
		customs.POST("/queryCheckResult", handler.QueryDataByKey) // 审核反馈查询

		customs.POST("/saveTaxBill", handler.SaveData)        // 3 税单
		customs.POST("/queryTaxBill", handler.QueryDataByKey) // 税单查询
	}

	warehouse := router.Group("/warehouse") // 仓储
	{
		// warehouse.POST("/saveWarehousesInfo", handler.SaveData)        // 1 入区
		warehouse.POST("/queryWarehousesInfo", handler.QueryDataByKey) // 入区查询

		// warehouse.POST("/saveDeliveryInfo", handler.SaveData)        // 2 出区
		// warehouse.POST("/queryDeliveryInfo", handler.QueryDataByKey) // 出区查询
	}

	bank := router.Group("/bank") // 银行——工商银行总
	{
		bank.POST("/saveIssueLC", handler.SaveData)        // 1 信用证开立信息
		bank.POST("/queryIssueLC", handler.QueryDataByKey) // 信用证开立信息查询

		bank.POST("/saveExaminationLC", handler.SaveData)        // 2 信用证回单信息审核
		bank.POST("/queryExaminationLC", handler.QueryDataByKey) // 信用证回单信息审核查询

		bank.POST("/saveRemittanceReceipt", handler.SaveData)        // 3 付汇水单
		bank.POST("/queryRemittanceReceipt", handler.QueryDataByKey) // 付汇水单查询

		bank.POST("/saveBankDeclarations", handler.SaveData)        // 4 银行补充
		bank.POST("/queryBankDeclarations", handler.QueryDataByKey) // 银行补充查询
	}

	declare := router.Group("/declare") // 申报端
	{
		// declare.POST("/saveBillLadingInfo", handler.SaveData)        // 1 提单及状态信息
		declare.POST("/queryBillLadingInfo", handler.QueryDataByKey) // 提单及状态信息查询

		// declare.POST("/saveInCustomsFiduciary", handler.SaveData)        // 2 入区申报受托
		// declare.POST("/queryInCustomsFiduciary", handler.QueryDataByKey) // 入区申报受托查询

		// declare.POST("/saveOutCustomsFiduciary", handler.SaveData)        // 3 出区申报受托
		// declare.POST("/queryOutCustomsFiduciary", handler.QueryDataByKey) // 出区申报受托查询

		declare.POST("/saveCustomsFormHeader", handler.SaveCustomsForm)    // 4 报关单表头信息
		declare.POST("/queryCustomsFormHeder", handler.QueryFormDataByKey) // 报关单表头信息查询

		// declare.POST("/saveCustomsFromBody", handler.SaveCustomsForm)     // 5 报关单表体信息
		// declare.POST("/queryCustomsFormBody", handler.QueryFormDataByKey) // 报关单表体信息查询

		// declare.POST("/saveDeclarations", handler.SaveData)        // 6 补充申报
		// declare.POST("/queryDeclarations", handler.QueryDataByKey) // 补充申报查询
	}

	all := router.Group("/all") // 所有角色
	{
		all.POST("/saveACL", handler.SaveACL)        // 设置业务访问权限控制
		all.POST("/queryACL", handler.QueryListById) // 查询业务访问权限控制
	}

	operator := router.Group("/operator") // 运营商
	{
		operator.POST("/saveUserInfo", handler.SaveUserInfo)   // 用户信息
		operator.POST("/queryUserInfo", handler.QueryUserInfo) // 用户信息查询
	}

	block := router.Group("/block") // 区块链网络操作
	{
		block.HEAD("/keepalive", handler.KeepaliveQuery) // 探活查询
		block.GET("/blockheight", handler.BlockHeight)   // 区块高度
		block.GET("/kafkaNumber", handler.KafkaNumber)   // kafka数量

		//log interface
		operator.POST("/setloglevel", handler.SetLogLevel)
		operator.GET("/getloglevel", handler.GetLogLevel)
	}

	return router
}
