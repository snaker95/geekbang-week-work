package service

import (
	"context"
	v1 "github.com/geekbang-week-work/week04/dudu-go/api/dudu-go/v1"
)

/*
@Time : 2021/6/20 下午5:28
@Author : snaker95
@File : hello.go
@Software: GoLand
*/

type HellService struct {
}

func (s *HellService) SayHello(ctx context.Context, in v1.HelloReq) (*v1.HelloResp, error) {
	return &v1.HelloResp{Greeting: "Hello " + in.Name}, nil
}
