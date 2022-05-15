package core

import (
	"context"
	"kubespace/server/service/alert"
	"fmt"
	"time"

	"go.uber.org/zap"
	"kubespace/server/global"
	"kubespace/server/initialize"
	"kubespace/server/service/system"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}

	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 Kubespace
	当前版本: V1.0.0
	默认swagger文档地址: http://127.0.0.1%s/swagger/index.html
	默认前端运行地址: http://127.0.0.1:8080
    默认后端运行地址: http://127.0.0.1:8888
`, address)
	ctx, _ := context.WithCancel(context.Background())
	go alert.Start(ctx)        //启动监控程序
	global.GVA_LOG.Error(s.ListenAndServe().Error())
	//ctx, cancel := context.WithCancel(context.Background())
	//go cron.TaskStart(ctx)     //启动定时任务
	//go utils.SyncLdapTask(ctx) //启动ldap同步任务
	//go alert.Start(ctx)        //启动监控程序
	//go cluster.NewClusterMgr() //启动contro-manager
	//ch := make(chan os.Signal)
	//signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)

/*	if err := s.ListenAndServe(); err != nil {
		global.GVA_LOG.Error(err.Error())
		cancel()
	}*/
}
