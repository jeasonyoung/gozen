package common

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// defPagIdx 默认页索引
const defPagIdx = 1

// defPagRows 默认每页数据量
const defPagRows = 10

// PagingQuery 分页查询条件
type PagingQuery interface {
	// GetIndex 获取页码
	GetIndex() int
	// GetRows 获取每页数量
	GetRows() int
}

func ParsePagingIndex(paging PagingQuery) (index, rows int) {
	index = defPagIdx
	rows = defPagRows

	if paging.GetIndex() > 0 {
		index = paging.GetIndex()
	}

	if paging.GetRows() > 0 {
		rows = paging.GetRows()
	}
	return
}

func Paginate(paging PagingQuery) func(m *gdb.Model) *gdb.Model {
	return func(m *gdb.Model) *gdb.Model {
		return m.Page(ParsePagingIndex(paging))
	}
}

// QueryPaginate 数据分页查询
func QueryPaginate(m *gdb.Model, query PagingQuery, items interface{}) (int, error) {
	//总行数
	total, err := m.Count()
	if err != nil {
		g.Log().Error(err)
		return 0, err
	}
	//分页处理查询
	if err = m.Handler(Paginate(query)).Scan(items); err != nil {
		g.Log().Error(err)
		return total, err
	}
	return total, err
}
