package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	MySQLDBUri    string
	POSTGRESQLURI string
	REDISURI      string
	// POSTGRESQL struct {
	// 	Host     string
	// 	Port     int
	// 	User     string
	// 	Password interface{}
	// 	DBName   string
	// }
	Logger struct {
		logx.LogConf
	}
	Etcd discov.EtcdConf
}
