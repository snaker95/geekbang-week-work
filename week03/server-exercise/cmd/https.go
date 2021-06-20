package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"server-exercise/internal/server"
	"syscall"
)

/*
@Time : 2021/6/10 上午10:03
@Author : snaker95
@File : https
@Software: GoLand
*/
type ab struct{
	A int
}
func main() {
	// 注册路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// var a *ab
		// fmt.Printf("R%v", a.A)
		fmt.Fprintln(w, "Hello server --- ", r.Host)
	})

	if err := run(); err != nil {
		fmt.Printf("==== main.run err ==== \n%v\n", err)
	}
}

func run() error {
	errG, ctx := errgroup.WithContext(context.Background())

	servers := []*server.HttpServer{
		server.NewHttpServer("127.0.0.1:8020"),
		server.NewHttpServer("127.0.0.1:8010"),
	}
	for i := range servers {
		s := servers[i]
		errG.Go(func() error {
			return s.Start(ctx)
		})
	}
	// 监听退出信号
	signChan := make(chan os.Signal, 1)
	defer close(signChan)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	errG.Go(func() error {
		select {
		case <-signChan:
			fmt.Printf("httpServer stop")
			// 应该有优化空间
			for _, s := range servers {
				fmt.Println("停止", s.Stop(ctx))
			}
			return nil
		}
	})
	return errG.Wait()
}
