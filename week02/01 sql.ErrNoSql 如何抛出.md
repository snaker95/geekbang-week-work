# 01 sql.ErrNoSql 如何抛出

## 问题

1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？



## 思路

###  sql.ErrNoRows 是什么

* 查看源码是如何抛出的

```go
func (r *Row) Scan(dest ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.rows.Close()
	for _, dp := range dest {
		if _, ok := dp.(*RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return ErrNoRows  //  重点 重点 重点
	}
	err := r.rows.Scan(dest...)
	if err != nil {
		return err
	}
	return r.rows.Close()
}
```

所以说, 从抛出的位置看到, 如果  `r.rows.Next() ` 为 **false**, 并且没有其他 error, 就是抛出 **ErrNoRows**

* 查 rows 是什么, Next是如何实现的

```go
type Row struct {
	err  error
	rows *Rows // 重点 重点 重点
}

type Rows struct {
	dc          *driverConn 
	releaseConn func(error)
	rowsi       driver.Rows
	cancel      func()      
	closeStmt   *driverStmt

	closemu sync.RWMutex
	closed  bool
	lasterr error 

	lastcols []driver.Value
}

// Next 的实现
func (rs *Rows) Next() bool {
	var doClose, ok bool
	withLock(rs.closemu.RLocker(), func() {
		doClose, ok = rs.nextLocked() // 重点 重点 重点
	})
	if doClose {
		rs.Close()
	}
	return ok
}

// nextLocked 实现
func (rs *Rows) nextLocked() (doClose, ok bool) {
	if rs.closed {
		return false, false
	}
	// ... 省略 ...

	rs.lasterr = rs.rowsi.Next(rs.lastcols)
	if rs.lasterr != nil {
		if rs.lasterr != io.EOF {
			return true, false
		}
		nextResultSet, ok := rs.rowsi.(driver.RowsNextResultSet)
		if !ok {
			return true, false
		}
		if !nextResultSet.HasNextResultSet() {
			doClose = true
		}
		return doClose, false
	}
	return false, true
}
```

其实深层次的并没有看懂, 只是了解到

1. go 官方制定了 database 的标准定义, 至于数据库驱动, 大多是有第三方完成;
2. 抛开数据库驱动实现, 该错误就是指**未查询到数据集**

###  发生sql.ErrNoRows，是否应该 Wrap 这个 error，抛给上层?

* 假设需要抛给上层，不应该 **Wrap**这个 **error**

  * 依据
    * 根据这个**error**产生的原因, 查询数据集为空时, 结合实际场景, 排查问题时不需要错误的堆栈信息, 没有必要通过 Wrap
  * 建议
    * 使用 **WithMessage** 添加上下文, 定位错误产生的源头

* 以上是建立在假设的基础上, 在实际使用场景中, 上层对于sql.ErrNoRows 应该是又爱又恨, 大体分为 **2** 种场景

  * 区分sql.ErrNoRows 和其他数据库错误
  * 不区分sql.ErrNoRows 和其他数据库错误

  所以, 还在纠结中..... 所以实现 2 种, 想咋就咋~~

### 代码

```go
package dao

// model 基本实现
// 第一部分: 统一实例化
// 第二部分: 表结构
// 第三部分: 表名
// 第四部分: 实例下的方法

import (
	"SnakerGin/config"
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// 统一实例化, 避免外部调用方法时多次实例化
var (
	CodeModel = Codes{}
)

// Codes 表结构
type Codes struct {
	ID         int64     `json:"-"` // 忽略主键, 业务场景下不需要
	Msg        string    `json:"msg"`
	RecoverMsg string    `json:"recover_msg"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName 表名
func (Codes) TableName() string {
	return "codes"
}

// GetListOrNil 查询列表, 空列表返回 nil
func (m *Codes) GetListOrNil(ctx context.Context, where *WhereCondition, limit, offset int, order ...string) ([]*Codes, error) {
	var list []*Codes
	db := getListHandler(ctx, where, limit, offset, order, &list)
	// 这里: 和 sql.ErrNotRows 是一样的道理
	// 不将ErrNotRows 返回上层
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return list, errors.WithMassage(db.Error, "Dao: ")
}

// GetListOrErr 查询列表, 空列表返回 err
func (m *Codes) GetListOrErr(ctx context.Context, where *WhereCondition, limit, offset int, order ...string) ([]*Codes, error) {
	var list []*Codes
	db := getListHandler(ctx, where, limit, offset, order, &list)
	return list, errors.WithMassage(db.Error, "Dao: ")
}

func (m *Codes) getListHandler(ctx context.Context, where *WhereCondition, limit, offset int, order []string, list interface{}) *gorm.DB {
	db := config.GetDBMaster().Model(m)
	if where != nil {
		db = db.Where(where.Query, where.Args...)
	}
	db = db.Limit(limit).Offset(offset)
	if len(order) > 0 && len(order[0]) > 0 {
		db = db.Order(order[0])
	}
	return db.Find(list)
}
```



