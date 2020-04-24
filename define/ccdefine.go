package define

type ResponseStatus struct {
	StatusCode int    `json:"statusCode"` //错误码0:成功1:失败
	StatusMsg  string `json:"statusMsg"`  //错误信息
}

type Page struct {
	CurrentPage  uint `json:"currentPage"`  //当前页码
	PageSize     uint `json:"pageSize"`     //每个页面显示个数
	TotalRecords uint `json:"totalRecords"` //总记录数
}
