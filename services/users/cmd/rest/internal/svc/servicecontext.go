package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/redis/go-redis/v9"
	"github.com/weridolin/site-gateway/services/users/cmd/rest/internal/config"
	"github.com/weridolin/site-gateway/services/users/models"
	"github.com/weridolin/site-gateway/tools"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	UserModel   models.UserModel
	RoleModel   models.Role
	RedisClient *redis.Client
}

func LoadInitData(c config.Config, DB *gorm.DB) {
	// 加载内置权限和角色
	// file := "etc/initdata.yaml"
	dir, _ := os.Getwd()
	file := path.Join(dir, "initData", "default.yaml")
	// fmt.Println("dir:", file)
	dataBytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("读取文件失败：", err)
		return
	}
	// fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	var defaultData struct {
		Menus         []models.Menu                `yaml:"menus"`
		Resources     map[string][]models.Resource `yaml:"resources"`
		Roles         []models.Role                `yaml:"roles"`
		RolesResource []map[string]interface{}     `yaml:"roles_resources"`
	}
	// var defaultData map[interface{}]interface{}
	err = yaml.Unmarshal(dataBytes, &defaultData)
	if err != nil {
		fmt.Println("解析 yaml 文件失败：", err)
		return
	}
	// fmt.Printf("defaultData → %+v\n", defaultData)
	// 初始数据插入数据库
	for _, resources := range defaultData.Resources {
		DB.Create(resources)
	}
	DB.Create(defaultData.Menus)
	DB.Create(defaultData.Roles)

	// 加载admin用户权限
	for _, permission := range defaultData.RolesResource {
		if permission["resource_id"] == "*" {
			var count int64
			DB.Model(&models.Resource{}).Count(&count)
			var resourceList []int
			for i := 1; i <= int(count); i++ {
				resourceList = append(resourceList, i)
			}
			err := models.BatchBindResourcePermission(resourceList, permission["role_id"].(int), DB)
			fmt.Println("err:", err)
		} else {
			var list []models.ResourcePermission
			for _, v := range permission["resource_id"].([]interface{}) {
				list = append(list, models.ResourcePermission{
					ResourceId: v.(int),
					RoleId:     permission["role_id"].(int),
				})
			}
			DB.Create(&list)
		}
	}

	// 缓存权限资源是否需要鉴权到redis,这里不放在内存是为了多个节点情况下不用去做同步
	redis_client := tools.NewRedisClient(c.REDISURI)
	var resourceCacheData []tools.ResourceAuthenticatedItem
	// resourceListJson, _ := json.Marshal(defaultData.Resources)
	// ctx := context.TODO()
	// res := redis_client.Set(ctx, tools.ResourceAuthenticatedCacheKey, resourceListJson, 0)
	// fmt.Printf("set resource auth cache → %+v  result →  %+v \n", resourceListJson, res)
	for _, resources := range defaultData.Resources {
		for _, resource := range resources {
			// ctx := context.TODO()
			// res := redis_client.Set(ctx, tools.ResourceAuthenticatedCacheKey, resource.Authenticated, 0)
			// fmt.Printf("set resource auth cache → %+v  result →  %+v \n", resource.Format(), res)
			resourceCacheData = append(resourceCacheData, tools.ResourceAuthenticatedItem{
				Resource:      resource.Format(),
				Authenticated: resource.Authenticated,
			})
		}
	}
	resourceListJson, _ := json.Marshal(resourceCacheData)
	ctx := context.TODO()
	res := redis_client.Set(ctx, tools.ResourceAuthenticatedCacheKey, resourceListJson, 0)
	fmt.Printf("set resource auth cache → %+v  result →  %+v \n", resourceCacheData, res)

}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(postgres.Open(c.POSTGRESQLURI), &gorm.Config{})

	// db, err := gorm.Open(mysql.Open(c.MySQLDBUri), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{
	// 		// TablePrefix:   "auth_", // 表名前缀，`User` 的表名应该是 `t_users`
	// 		SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
	// 	},
	// })
	if err != nil {
		logx.Error(err)
	}
	//自动同步更新表结构
	fmt.Println("自动同步更新表结构")
	db.AutoMigrate(&models.MenuPermission{})
	db.AutoMigrate(&models.ResourcePermission{})
	db.AutoMigrate(&models.UserRoles{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Role{})

	db.AutoMigrate(&models.Menu{})
	db.AutoMigrate(&models.Resource{})

	LoadInitData(c, db)

	return &ServiceContext{
		Config:      c,
		DB:          db,
		UserModel:   models.NewUserModel("user"),
		RoleModel:   models.Role{},
		RedisClient: tools.NewRedisClient(c.REDISURI),
	}
}
