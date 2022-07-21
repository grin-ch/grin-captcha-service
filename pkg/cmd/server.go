package cmd

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grin-ch/grin-api/api/captcha"
	"github.com/grin-ch/grin-captcha-service/cfg"
	"github.com/grin-ch/grin-captcha-service/pkg/model"
	"github.com/grin-ch/grin-captcha-service/pkg/service"
	center "github.com/grin-ch/grin-etcd-center"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// RunServer 运行服务
func RunServer() error {
	var cfgName, cfgPath string
	flag.StringVar(&cfgName, "cfgName", "service", "log file")
	flag.StringVar(&cfgPath, "cfgPath", "./cfg", "log path")
	flag.Parse()
	cfg.SetServerConfig(cfgName, cfgPath)

	return grpcServer(cfg.GetConfig())
}

// 运行grpc服务
func grpcServer(c *cfg.ServerConfig) error {
	initLogger(c.LogPath, c.LogLevel, c.LogColor, c.LogCaller)

	// grpc listener
	grpcListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", c.Port))
	if err != nil {
		log.Errorf("tcp listen err:%s", err.Error())
		return err
	}

	// grpc server
	svc := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryFunc)),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryFunc)),
		),
	)
	gracefulShutdown(svc)

	// 获取注册中心连接
	etcdCenter, err := center.NewEtcdCenter(c.RegEndpoint, c.RegTimeout)
	if err != nil {
		log.Errorf("new etcd client err:%s", err.Error())
		return err
	}
	// 服务注册
	registrar := etcdCenter.Registrar(c.Name, c.Host, c.Port, c.RegTimeout)
	err = registrar.Registry()
	if err != nil {
		log.Errorf("server registry err: %s", err.Error())
		return err
	}
	defer registrar.Deregistry()
	// 注册grpc服务
	if err := registryGrpcServices(svc, c); err != nil {
		log.Errorf("grpc server registry err:%s", err.Error())
		return err
	}
	log.Infof("grpc server is running: %s", fmt.Sprintf("%s:%d", c.Host, c.Port))
	return svc.Serve(grpcListener)
}

// 优雅退出
func gracefulShutdown(svc *grpc.Server) {
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sign := <-quit

		// 关闭服务链接
		svc.GracefulStop()
		log.Infof("grpc server shutdown: %s", sign)
	}()
}

func registryGrpcServices(svc *grpc.Server, c *cfg.ServerConfig) error {
	pv, err := model.RegistryProvider(c.RedisAddr, c.RedisPass, c.RedisDB)
	if err != nil {
		log.Errorf("registry redis error:%s", err.Error())
		return err
	}
	captcha.RegisterCaptchaServiceServer(svc, service.NewCaptchaServer(pv))

	return nil
}
