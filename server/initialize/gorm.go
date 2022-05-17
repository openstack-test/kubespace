package initialize

import (
	"fmt"
	"kubespace/server/model/kubernetes"
	"os"

	"kubespace/server/global"
	"kubespace/server/model/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"kubespace/server/utils/monitoring/prometheus"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	default:
		return GormMysql()
	}
}

//初始化prometheus链接
func InitPrometheusClient() *prometheus.Prometheus {
	cli, err := prometheus.NewPrometheus(
		&prometheus.Options{Endpoint: fmt.Sprintf("http://%s:%s/", global.GVA_CONFIG.Monitor.Host, global.GVA_CONFIG.Monitor.Port)})
	if err != nil {
		fmt.Println("init monitor client error", err)
		panic(err)
	}
	return cli
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		// 系统模块表
		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCode{},

		// Kubernetes模块表
        kubernetes.K8SCluster{},

	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
