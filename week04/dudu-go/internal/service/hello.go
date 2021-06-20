package service

import (
	"context"
	v1 "github.com/geekbang-week-work/week04/dudu-go/api/dudu-go/v1"
	"github.com/geekbang-week-work/week04/dudu-go/internal/biz"
)

/*
@Time : 2021/6/20 下午5:28
@Author : snaker95
@File : hello.go
@Software: GoLand
*/

type HellService struct {
	Uc *biz.HelloUsecase
}

func (s *HellService) SayHello(ctx context.Context, in *v1.HelloReq) (*v1.HelloResp, error) {
	return &v1.HelloResp{Greeting: "Hello " + in.Name}, nil
}

func (s *HellService) GetHello(ctx context.Context, in *v1.GetHelloReq) (*v1.GetHelloResp, error) {
	// pto -> po
	req := map[string]interface{}{
		"name": in.Name,
	}
	hellos, err := s.Uc.Get(ctx, req)
	resp := &v1.GetHelloResp{
		List: make([]*v1.HelloOne, 0, len(hellos)),
	}
	// hellos -> resp 拷贝
	// 示例: 简单处理
	for i := range hellos {
		resp.List = append(resp.List, &v1.HelloOne{
			Name: hellos[i].Name,
		})
	}
	return resp, err
}
