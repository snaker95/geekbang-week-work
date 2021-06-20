package main

import (
	"context"
	"flag"
	"fmt"
	v1 "github.com/geekbang-week-work/week04/dudu-go/api/dudu-go/v1"
	"github.com/geekbang-week-work/week04/dudu-go/internal/conf"
	"github.com/geekbang-week-work/week04/dudu-go/internal/server"
	"github.com/geekbang-week-work/week04/dudu-go/internal/service"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

/*
@Time : 2021/6/17 下午10:20
@Author : snaker95
@File : main
@Software: GoLand
*/

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/release.yaml", "config path, eg: -conf config.yaml")
	flag.Parse()
}
func main() {
	// 初始化上下文
	ctx := context.Background()
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		hello := new(service.HellService)
		resp, err := hello.SayHello(ctx, v1.HelloReq{Name:"oooxxx"})
		fmt.Fprintln(w, "Hello server --- ", resp, err)
	})


	// 加载配置
	var config = new(conf.Config)
	if err := conf.Scan(flagconf, config); err != nil {
		panic(fmt.Sprintf("conf.Scan err=%+v", err))
	}
	fmt.Printf("conf = %v", config)
	// 使用 wire 依赖注入, 获取启动入口
	app := initApp(ctx, config.Server)
	if err := app.run(); err != nil {
		panic(fmt.Sprintf("app.run err=%+v", err))
	}
}

// wire完成依赖注入

type App struct {
	ctx   context.Context
	https []*server.HttpServer
}

// 初始化 app
func initApp(ctx context.Context, servers *conf.Server) *App {
	httpSrv := []*server.HttpServer{
		server.NewHttpServer(servers.Http.Addr),
	}
	app := &App{
		ctx:   ctx,
		https: httpSrv,
	}
	return app
}

func (a *App) run() error {
	servers := a.https
	errG, ctx := errgroup.WithContext(a.ctx)
	// 循环启动服务
	for i := range servers {
		s := servers[i]
		errG.Go(func() error {
			fmt.Println("httpServer start", s.Addr)
			return s.Start(ctx)
		})
	}
	// 监听退出信号
	signChan := make(chan os.Signal, 1)
	defer close(signChan)
	signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	// 异步处理退出信号
	errG.Go(func() error {
		select {
		case <-signChan:
			fmt.Println("httpServer stop")
			// 应该有优化空间
			for _, s := range servers {
				fmt.Println("停止", s.Addr, s.Stop(ctx))
			}
			return nil
		}
	})
	// errGroup 等待
	return errG.Wait()
}
