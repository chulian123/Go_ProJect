package main

import (
	"github.com/gin-gonic/gin"
	srv "test.com/project-common"
	"test.com/project-project/config"
	"test.com/project-project/router"
)

func main() {
	r := gin.Default()
	//	tp, tpErr := tracing.JaegerTraceProvider(config.C.JaegerConfig.Endpoints)
	//	if tpErr != nil {
	//		log.Fatal(tpErr)
	//	}
	//	otel.SetTracerProvider(tp)
	//	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	//路由
	router.InitRouter(r)
	//初始化rpc调用
	router.InitUserRpc()
	//grpc服务注册
	gc := router.RegisterGrpc()
	//grpc服务注册到etcd
	router.RegisterEtcdServer()
	//初始化kafka
	c := config.InitKafkaWriter()
	//初始化kafka消费者
	reader := config.NewCacheReader()
	go reader.DeleteCache()
	stop := func() {
		gc.Stop()
		c()
		reader.R.Close()
	}
	srv.Run(r, config.C.SC.Name, config.C.SC.Addr, stop)
}
