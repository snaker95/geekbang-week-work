package biz

import "github.com/google/wire"

/*
@Time : 2021/6/20 下午7:14
@Author : snaker95
@File : wire_set
@Software: GoLand
*/

var ProviderSet = wire.NewSet(NewHelloUsecase)
