package data

import "github.com/geekbang-week-work/week04/dudu-go/internal/conf"

/*
@Time : 2021/6/20 下午7:00
@Author : snaker95
@File : data
@Software: GoLand
*/

type Data struct {
	// 数据库 handler
	dbHandler interface{}
}

func NewData(conf *conf.Data) (*Data, error) {
	return &Data{}, nil
}
