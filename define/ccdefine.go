package define

type Page struct {
	CurrentPage  uint `json:"currentPage"`  //当前页码
	PageSize     uint `json:"pageSize"`     //每个页面显示个数
	TotalRecords uint `json:"totalRecords"` //总记录数
}
