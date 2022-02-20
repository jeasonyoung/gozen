package common

// ReqPagingQuery 分页请求-报文体
type ReqPagingQuery struct {
	Index int `json:"index" p:"index" default:"1"` //当前页码
	Rows  int `json:"rows" p:"rows" default:"20"`  //当前页数据量
}

func (p *ReqPagingQuery) GetIndex() int {
	return p.Index
}

func (p *ReqPagingQuery) GetRows() int {
	return p.Rows
}
