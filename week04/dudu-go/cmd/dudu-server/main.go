package main

import (
	"context"
	"flag"
	"fmt"
	v1 "github.com/geekbang-week-work/week04/dudu-go/api/dudu-go/v1"
	"github.com/geekbang-week-work/week04/dudu-go/internal/biz"
	"github.com/geekbang-week-work/week04/dudu-go/internal/conf"
	"github.com/geekbang-week-work/week04/dudu-go/internal/data"
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

	// 加载配置
	var config = new(conf.Config)
	if err := conf.Scan(flagconf, config); err != nil {
		panic(fmt.Sprintf("conf.Scan err=%+v", err))
	}
	fmt.Printf("conf = %v", config)
	// 使用 wire 依赖注入, 获取启动入口
	app := initApp(ctx, config.Server, config.Data)
	if err := app.run(); err != nil {
		panic(fmt.Sprintf("app.run err=%+v", err))
	}
}

// wire完成依赖注入

type App struct {
	ctx   context.Context
	https []*server.HttpServer
	dbs   *conf.Data
}

// 初始化 app
func initApp(ctx context.Context, servers *conf.Server, dbs *conf.Data) *App {
	httpSrv := []*server.HttpServer{
		server.NewHttpServer(servers.Http.Addr),
	}
	app := &App{
		ctx:   ctx,
		https: httpSrv,
		dbs:   dbs,
	}
	return app
}

func (app *App) run() error {
	servers := app.https
	errG, ctx := errgroup.WithContext(app.ctx)

	// 循环启动服务
	for i := range servers {

		s := servers[i]
		mux := http.NewServeMux()
		mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			hello := new(service.HellService)
			resp, err := hello.SayHello(ctx, &v1.HelloReq{Name: "9900"})
			fmt.Fprintln(w, "Hello server --- ", resp, err)
		})
		mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
			d, err := data.NewData(app.dbs)
			if err != nil {
				panic(fmt.Sprintf("data.NewData err=%v", err))
			}
			repo := data.NewHelloRepo(d)
			uc := biz.NewHelloUsecase(repo)
			hello := service.NewHellService(uc)

			resp, err := hello.GetHello(ctx, &v1.GetHelloReq{Name: "get"})
			fmt.Fprintln(w, "Hello server get --- ", resp, err)
		})
		s.SetHandler(mux)
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
