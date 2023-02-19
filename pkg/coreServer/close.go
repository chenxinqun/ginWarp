package coreServer

import (
	"context"
	"net/http"
	"time"

	"github.com/chenxinqun/ginWarpPkg/sysx/shutdown"
	"go.uber.org/zap"
)

var ser *Server
var httpSer *http.Server

func SetServer(s *Server) {
	ser = s
}

func SetHttpServer(hs *http.Server) {
	httpSer = hs
}

var (
	closeHandlers = make([]shutdown.CloseFunc, 0)
)

func RegisterCloseFunc(closeFunc shutdown.CloseFunc) {
	closeHandlers = append(closeHandlers, closeFunc)
}

func CloseCoreServer() {
	// 关闭 httpServer
	closeHandlers = append(closeHandlers, CloseHttp)
	// 关闭MySQL连接
	closeHandlers = append(closeHandlers, CloseMysql)
	// 关闭etcd连接
	closeHandlers = append(closeHandlers, CloseEtcd)
	// 关闭mongo连接
	closeHandlers = append(closeHandlers, CloseMongo)
	// 关闭Redis连接
	closeHandlers = append(closeHandlers, CloseRedis)
	// 关闭Kafka连接
	closeHandlers = append(closeHandlers, CloseKafka)

	// 优雅关闭
	shutdown.NewHook().Close(closeHandlers...)
}

// CloseHttp 关闭 httpServer
func CloseHttp() {
	if httpSer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := httpSer.Shutdown(ctx); err != nil {
			ser.Logger.Error("coreServer shutdown err", zap.Error(err))
		}
	}

}

// CloseEtcd 关闭etcd连接
func CloseEtcd() {
	if ser.EtcdRepo != nil {
		if err := ser.EtcdRepo.Close(); err != nil {
			ser.Logger.Error("EtcdRepo close err", zap.Error(err))
		}
	}
}

// CloseMysql 关闭 db连接
func CloseMysql() {
	if ser.MysqlRepo != nil {
		if err := ser.MysqlRepo.Close(); err != nil {
			ser.Logger.Error("dbw close err", zap.Error(err))
		}

		if err := ser.MysqlRepo.Close(); err != nil {
			ser.Logger.Error("dbr close err", zap.Error(err))
		}
	}
}

// CloseMongo 关闭MongoDB连接
func CloseMongo() {
	if ser.MongoRepo != nil {
		ser.MongoRepo.Close()
	}
}

// CloseRedis 关闭 Redis连接
func CloseRedis() {
	if ser.RedisRepo != nil {
		if err := ser.RedisRepo.Close(); err != nil {
			ser.Logger.Error("redis close err", zap.Error(err))
		}
	}
}

// CloseKafka 关闭Kafka连接
func CloseKafka() {
	if ser.ProducerRepo != nil {
		_ = ser.ProducerRepo.Close()
	}
	if ser.AsyncProducerRepo != nil {
		_ = ser.AsyncProducerRepo.Close()
	}
	if ser.ConsumerGroupRepo != nil {
		_ = ser.ConsumerGroupRepo.Close()
	}
}
