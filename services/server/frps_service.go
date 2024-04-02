package server

import (
	"context"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/metrics/mem"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/server"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

type ServerHandler interface {
	Run()
	Stop()
	GetCommonCfg() *v1.ServerConfig
	GetMem() *mem.ServerStats
}

type Server struct {
	srv    *server.Service
	Common *v1.ServerConfig
}

var (
	srv *Server
)

func InitGlobalServerService(svrCfg *v1.ServerConfig) {
	if srv != nil {
		logrus.Warn("server has been initialized")
		return
	}

	svrCfg.Complete()
	srv = NewServerHandler(svrCfg)
}

func GetGlobalServerSerivce() ServerHandler {
	if srv == nil {
		logrus.Panic("server has not been initialized")
	}
	return srv
}

func GetServerSerivce(svrCfg *v1.ServerConfig) ServerHandler {
	svrCfg.Complete()
	return NewServerHandler(svrCfg)
}

func NewServerHandler(svrCfg *v1.ServerConfig) *Server {
	warning, err := validation.ValidateServerConfig(svrCfg)
	if warning != nil {
		logrus.WithError(err).Warnf("validate server config warning: %+v", warning)
	}
	if err != nil {
		logrus.Panic(err)
	}

	log.InitLogger(svrCfg.Log.To, svrCfg.Log.Level, int(svrCfg.Log.MaxDays), svrCfg.Log.DisablePrintColor)

	var svr *server.Service

	if svr, err = server.NewService(svrCfg); err != nil {
		logrus.WithError(err).Panic("cannot create server, exit and restart")
	}

	return &Server{
		srv:    svr,
		Common: svrCfg,
	}
}

func (s *Server) Run() {
	wg := conc.NewWaitGroup()
	wg.Go(func() { s.srv.Run(context.Background()) })
	wg.Wait()
}

func (s *Server) Stop() {
	wg := conc.NewWaitGroup()
	wg.Go(func() {
		err := s.srv.Close()
		if err != nil {
			logrus.Errorf("close server error: %v", err)
		}
		logrus.Infof("server closed")
	})
	wg.Wait()
}

func (s *Server) GetCommonCfg() *v1.ServerConfig {
	return s.Common
}

func (s *Server) GetMem() *mem.ServerStats {
	return mem.StatsCollector.GetServer()
}
