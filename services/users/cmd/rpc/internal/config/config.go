package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DBUri   string
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	POSTGRESQLURI string
}
