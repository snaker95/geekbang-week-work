package server

import (
	"context"
	"net/http"
)

/*
@Time : 2021/6/10 上午9:09
@Author : snaker95
@File : http
@Software: GoLand
*/

type HttpServer struct {
	http.Server
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{
		Server: http.Server{
			Addr: addr,
		},
	}
}

func (s *HttpServer) SetHandler(handler http.Handler) {
	s.Handler = handler
}

func (s *HttpServer) Start(ctx context.Context) error {
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
