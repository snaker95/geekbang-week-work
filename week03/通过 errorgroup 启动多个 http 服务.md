# 通过 errorgroup 启动多个 http 服务

## 思路:

* 定义 server 服务的接口:
  * Start()
  * Stop()
* 定义 自定义的 server Struct; 
  * http.server
  * stop chan struct{}
* 实现 server 接口
* 定义 NewServer 方法



## 信号监听

```go
signChan := make(chan os.Signal, 1)
defer close(signChan)
signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
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
```

