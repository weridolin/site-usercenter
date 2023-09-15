package svc

import (
	"github.com/weridolin/site-gateway/services/users/cmd/rpc/internal/config"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	UserModel models.UserModel
	RoleModel models.Role
}

func NewServiceContext(c config.Config) *ServiceContext {
	// db, err := gorm.Open(mysql.Open(c.DBUri), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{
	// 		// TablePrefix:   "auth_", // 表名前缀，`User` 的表名应该是 `t_users`
	// 		SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
	// 	},
	// })
	db, err := gorm.Open(postgres.Open(c.POSTGRESQLURI), &gorm.Config{})
	if err != nil {
		logx.Error(err)
	}
	return &ServiceContext{
		Config:    c,
		DB:        db,
		UserModel: models.NewUserModel("user"),
		RoleModel: models.Role{}}
}
