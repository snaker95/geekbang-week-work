package data

import "github.com/google/wire"

/*
@Time : 2021/6/20 下午7:06
@Author : snaker95
@File : wire_set
@Software: GoLand
*/

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewHelloRepo)
