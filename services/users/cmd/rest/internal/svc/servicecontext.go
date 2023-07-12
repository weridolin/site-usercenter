package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/config"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/weridolin/site-gateway/tools"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	UserModel   models.UserModel
	RoleModel   models.Role
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DBUri), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "auth_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		logx.Error(err)
	}
	//自动同步更新表结构
	db.AutoMigrate(&models.MenuPermission{})
	db.AutoMigrate(&models.ResourcePermission{})
	db.AutoMigrate(&models.UserRoles{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.Menu{})
	db.AutoMigrate(&models.Resource{})

	return &ServiceContext{
		Config:      c,
		DB:          db,
		UserModel:   models.NewUserModel("user"),
		RoleModel:   models.Role{},
		RedisClient: tools.NewRedisClient(),
	}
}
