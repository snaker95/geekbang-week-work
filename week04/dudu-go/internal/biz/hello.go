package biz

import "context"

/*
@Time : 2021/6/20 下午5:44
@Author : snaker95
@File : hello.go
@Software: GoLand
*/

// po 数据层
type Hello struct {
	Name string `json:"name"`
}

type HelloRepo interface {
	SaveHello(context.Context, *Hello) error
	GetHellos(context.Context, map[string]interface{}) ([]*Hello, error)
}

type HelloUsecase struct {
	Repo HelloRepo
}

func (hu *HelloUsecase) Create(ctx context.Context, hello *Hello) error {
	return hu.Repo.SaveHello(ctx, hello)
}

func (hu *HelloUsecase) Get(ctx context.Context, where map[string]interface{}) ([]*Hello, error) {
	return hu.Repo.GetHellos(ctx, where)
}
