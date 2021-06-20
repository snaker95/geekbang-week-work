package data

import (
	"context"
	"github.com/geekbang-week-work/week04/dudu-go/internal/biz"
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

type HelloRepo struct {
}

func (hr *HelloRepo) SaveHello(ctx context.Context, biz *biz.Hello) error {
	// biz -> hello 数据深度拷贝(转化)

	// 数据数操作 todo

	return nil
}

func (hr *HelloRepo) GetHellos(ctx context.Context, where map[string]interface{}) ([]*biz.Hello, error) {
	// 数据库操作 todo

	// hello -> biz 数据深度拷贝(转化)
	return nil, nil
}
