package data

import (
	"context"
	"github.com/geekbang-week-work/week04/dudu-go/internal/biz"
	"github.com/pkg/errors"
)

/*
@Time : 2021/6/20 下午5:44
@Author : snaker95
@File : hello.go
@Software: GoLand
*/

var helloTable = "hello"

// Do 数据层
type hello struct {
	Id   int64  `json:"id" db:"name" ak:"true"`
	Name string `json:"name" db:"name"`
}

type helloRepo struct {
	data *Data
}

func NewHelloRepo(data *Data) biz.HelloRepo {
	return &helloRepo{
		data: data,
	}
}

func (hr *helloRepo) SaveHello(ctx context.Context, biz *biz.Hello) error {
	// biz(DO) -> hello(PO) 数据深度拷贝(转化)

	// 数据数操作 todo
	// hr.dbHandler.Insert()

	return nil
}

func (hr *helloRepo) GetHellos(ctx context.Context, where map[string]interface{}) ([]*biz.Hello, error) {
	// 数据库操作 todo

	// hello(PO) -> biz(DO) 数据深度拷贝(转化)
	var err error
	return nil, errors.Wrap(err, "data:")
}
