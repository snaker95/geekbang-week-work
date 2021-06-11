package server

import "context"

/*
@Time : 2021/6/10 上午9:08
@Author : snaker95
@File : interface
@Software: GoLand
*/

type ServerInterface interface {
	Start(context.Context) error
	Stop(context.Context) error
}
