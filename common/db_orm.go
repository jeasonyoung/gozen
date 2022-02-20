package common

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// DataQuery 数据查询
type DataQuery struct {
	model *gdb.Model
}

func NewDataQuery(model *gdb.Model) *DataQuery {
	query := &DataQuery{}
	query.model = model
	return query
}

func (dq *DataQuery) AddFn(fn ...func(model *gdb.Model)) *DataQuery {
	if fn != nil && len(fn) > 0 {
		m := dq.model
		for _, f := range fn {
			if f != nil {
				f(m)
			}
		}
	}
	return dq
}

func (dq *DataQuery) Add(condition bool, field string, val interface{}) *DataQuery {
	if condition && field != "" && val != nil {
		dq.model = dq.model.Where(field, val)
	}
	return dq
}

func (dq *DataQuery) AddLike(condition bool, field string, val interface{}) *DataQuery {
	if condition && field != "" && val != nil {
		dq.model = dq.model.WhereLike(field, val)
	}
	return dq
}

func (dq *DataQuery) AddWhere(field string, val interface{}) *DataQuery {
	if field != "" && val != nil {
		dq.model = dq.model.Where(field, val)
	}
	return dq
}

// QueryPaginate  分页查询
func (dq *DataQuery) QueryPaginate(query PagingQuery, items interface{}) (int, error) {
	return QueryPaginate(dq.model, query, items)
}

// QueryResut 查询结果
func (dq *DataQuery) QueryResut(result interface{}) error {
	if result != nil {
		return dq.model.Scan(result)
	}
	return nil
}

// DataUpdate 数据更新
type DataUpdate struct {
	model   *gdb.Model
	dataVal g.Map
}

func NewDataUpdate(model *gdb.Model) *DataUpdate {
	data := &DataUpdate{}
	data.model = model
	data.dataVal = g.Map{}
	return data
}

func (du *DataUpdate) Add(condition bool, field string, val interface{}) *DataUpdate {
	if condition && field != "" {
		du.dataVal[field] = val
	}
	return du
}

func (du *DataUpdate) AddSet(field string, val interface{}) *DataUpdate {
	if field != "" {
		du.dataVal[field] = val
	}
	return du
}

func (du *DataUpdate) UpdateWithPri(id uint64) error {
	if len(du.dataVal) > 0 && id > 0 {
		if _, err := du.model.Data(du.dataVal).WherePri(id).Update(); err != nil {
			g.Log().Error(err)
			return err
		}
	}
	return nil
}

// DataTxHandler 数据处理事务
func DataTxHandler(db gdb.DB, handlers ...func(tx *gdb.TX, err error)) error {
	if db != nil && len(handlers) > 0 {
		//开启事务
		tx, err := g.DB().Begin()
		if err != nil {
			g.Log().Error(err)
			return err
		}
		//事务提交处理
		defer func() {
			if err != nil {
				if err = tx.Rollback(); err != nil {
					g.Log().Error(err)
				}
			} else {
				if err = tx.Commit(); err != nil {
					g.Log().Error(err)
				}
			}
		}()
		//业务处理
		for _, handler := range handlers {
			//业务处理
			handler(tx, err)
			if err != nil {
				g.Log().Error(err)
				return err
			}
		}
	}
	return nil
}
