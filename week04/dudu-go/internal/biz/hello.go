package biz

import (
	"context"
	"github.com/pkg/errors"
)

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
	repo HelloRepo
}

func NewHelloUsecase(repo HelloRepo) *HelloUsecase {
	return &HelloUsecase{
		repo: repo,
	}
}

func (hu *HelloUsecase) Create(ctx context.Context, hello *Hello) error {
	return errors.WithMessage(hu.repo.SaveHello(ctx, hello), "biz: ")
}

func (hu *HelloUsecase) Get(ctx context.Context, where map[string]interface{}) ([]*Hello, error) {
	resp, err := hu.repo.GetHellos(ctx, where)
	return resp, errors.WithMessage(err, "biz")
}
