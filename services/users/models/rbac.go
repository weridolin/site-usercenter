package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	CreateUser  int         `gorm:"comment:创建者ID" json:"create_user" yaml:"create_user"`
	Zone        string      `gorm:"uniqueIndex:udx_role;not null;comment:域;size:256;default:site" yaml:"zone" json:"zone"` //角色的作用范围
	Name        string      `gorm:"uniqueIndex:udx_role;not null;comment:角色名;size:256" json:"name" yaml:"name"`
	Users       []*User     `gorm:"many2many:user_roles;ForeignKey:ID;JoinForeignKey:RoleId;References:ID;joinReferences:UserId"`
	Description string      `gorm:"comment:角色描述" json:"description" yaml:"description"`
	Menus       []*Menu     `gorm:"many2many:auth_menu_permission;ForeignKey:ID;JoinForeignKey:RoleId;References:ID;joinReferences:MenuId;"`
	Resources   []*Resource `gorm:"many2many:auth_resource_permission;ForeignKey:ID;JoinForeignKey:RoleId;References:ID;joinReferences:ResourceId;"`
}

func (Role) TableName() string {
	return "auth_role"
}

type UserRoles struct {
	gorm.Model
	UserId uint `gorm:"comment:用户ID"`
	RoleId uint `gorm:"comment:角色ID"`
}

func (UserRoles) TableName() string {
	return "user_roles"
}

func (r *Role) Create(zone, name, description string, userID int, DB *gorm.DB) (*Role, error) {
	new := Role{
		CreateUser:  userID,
		Zone:        zone,
		Name:        name,
		Description: description,
	}
	err := DB.Create(&new).Error
	return &new, err
}

func (r *Role) Delete(id int, UserId int, DB *gorm.DB) error {
	return DB.Where("id = ? and create_user = ?", id, UserId).Delete(&Role{}).Error
}

func (r *Role) Update(id int, data map[string]interface{}, DB *gorm.DB) (*Role, error) {
	role, err := r.QueryById(map[string]interface{}{"id": id}, DB)
	if err != nil {
		return nil, err
	} else {
		err := DB.Model(&role).Updates(data).Error
		return role, err
	}
}

func (t *Role) QueryById(condition interface{}, DB *gorm.DB) (*Role, error) {
	var role *Role
	err := DB.Table("auth_role").Where(condition).First(&role).Error
	return role, err
}

func (r *Role) Query(condition interface{}, DB *gorm.DB) ([]*Role, error) {
	var roles []*Role
	err := DB.Table("auth_role").Preload("Menus").Preload("Resources").Where(condition).Find(&roles).Error
	return roles, err
}

func (r UserRoles) Query(condition interface{}, DB *gorm.DB) ([]*UserRoles, error) {
	var roles []*UserRoles
	err := DB.Table("user_roles").Where(condition).Find(&roles).Error
	return roles, err
}

type Menu struct {
	gorm.Model
	// ServerName string `gorm:"uniqueIndex:udx_menu;not null;comment:服务名"`
	Name      string `gorm:"uniqueIndex:udx_menu;not null;comment:菜单名;size:256" josn:"name" yaml:"name"`
	ParentId  int    `gorm:"comment:父菜单ID" json:"parent_id" yaml:"parent_id"`
	Url       string `gorm:"comment:菜单路径" json:"url" yaml:"url"`
	Component string `gorm:"comment:菜单组件(Vue)" json:"component" yaml:"component"`
	Icon      string `gorm:"comment:菜单图标" json:"icon" yaml:"icon"`
	RouteName string `gorm:"comment:菜单路由名称" json:"route_name" yaml:"route_name"`
	Type      int    `gorm:"comment:菜单类型:0菜单 1按钮;default:0;size:1" json:"type" yaml:"type"`
}

func (Menu) TableName() string {
	return "auth_menu"
}

func (m *Menu) Create(name, url, component, icon, routeName string, parentId, menuType int, DB *gorm.DB) (*Menu, error) {
	new := Menu{
		Name:      name,
		ParentId:  parentId,
		Url:       url,
		Component: component,
		Icon:      icon,
		RouteName: routeName,
		Type:      menuType,
	}
	err := DB.Create(&new).Error
	return &new, err
}

func (m *Menu) Delete(id int, DB *gorm.DB) error {
	return DB.Where("id = ?", id).Delete(&Menu{}).Error
}

func (m *Menu) Update(id int, data map[string]interface{}, DB *gorm.DB) (*Menu, error) {
	menu, err := m.QueryById(map[string]interface{}{"id": id}, DB)
	if err != nil {
		return nil, err
	} else {
		err := DB.Model(&menu).Updates(data).Error
		return menu, err
	}
}

func (m *Menu) QueryById(condition interface{}, DB *gorm.DB) (*Menu, error) {
	var menu *Menu
	err := DB.Table("auth_menu").Where(condition).First(&menu).Error
	return menu, err
}

type MenuPermission struct {
	gorm.Model
	MenuId uint `gorm:"comment:菜单ID"`
	RoleId uint `gorm:"comment:角色ID"`
}

func (MenuPermission) TableName() string {
	return "auth_menu_permission"
}

func BindMenuRole(menuId, roleId uint, DB *gorm.DB) error {
	return DB.Create(&MenuPermission{
		MenuId: menuId,
		RoleId: roleId,
	}).Error
}

func BatchBindMenuRole(MenuIdList []uint, roleId uint, DB *gorm.DB) error {
	var list []MenuPermission
	for _, v := range MenuIdList {
		list = append(list, MenuPermission{
			MenuId: v,
			RoleId: roleId,
		})
	}
	return DB.Create(&list).Error
}

// 资源api鉴权
type Resource struct {
	gorm.Model
	ServerName  string `gorm:"uniqueIndex:udx_resource;not null;comment:服务名;size:256" json:"server_name" yaml:"server_name"`
	Url         string `gorm:"uniqueIndex:udx_resource;not null;comment:资源路径;size:256" json:"url" yaml:"url"`
	Method      string `gorm:"uniqueIndex:udx_resource;not null;comment:资源方法 ;set(GET,POST,PUT,DELETE);size:256" json:"method" yaml:"method"`
	Version     string `gorm:"comment:资源版本;size:256" json:"version" yaml:"version"`
	Description string `gorm:"comment:资源描述;size:256" json:"description" yaml:"description" `
}

func (Resource) TableName() string {
	return "auth_resource"
}

func (r *Resource) Create(serverName, url, method, version, description string, DB *gorm.DB) (*Resource, error) {
	new := Resource{
		ServerName:  serverName,
		Url:         url,
		Method:      method,
		Version:     version,
		Description: description,
	}
	err := DB.Create(&new).Error
	return &new, err
}

func (r *Resource) Delete(id int, DB *gorm.DB) error {
	return DB.Where("id = ?", id).Delete(&Resource{}).Error
}

func (r *Resource) Update(id int, data map[string]interface{}, DB *gorm.DB) (*Resource, error) {
	resource, err := r.QueryById(map[string]interface{}{"id": id}, DB)
	if err != nil {
		return nil, err
	} else {
		err := DB.Model(&resource).Updates(data).Error
		return resource, err
	}
}

func (r *Resource) QueryById(condition interface{}, DB *gorm.DB) (*Resource, error) {
	var resource *Resource
	err := DB.Table("auth_resource").Where(condition).First(&resource).Error
	return resource, err
}

func (r *Resource) Format() string {
	// {{servername}}/api/{{version}}/{{url}}:{{method}}
	return r.ServerName + "/api/" + r.Version + r.Url + ":" + r.Method
}

type ResourcePermission struct {
	gorm.Model
	ResourceId uint `gorm:"comment:资源ID"`
	RoleId     uint `gorm:"comment:角色ID"`
}

func (ResourcePermission) TableName() string {
	return "auth_resource_permission"
}

func BindResourcePermission(resourceId, roleId uint, DB *gorm.DB) error {
	return DB.Create(&ResourcePermission{
		ResourceId: resourceId,
		RoleId:     roleId,
	}).Error
}

func BatchBindResourcePermission(resourceIdList []uint, roleId uint, DB *gorm.DB) error {
	var list []ResourcePermission
	for _, v := range resourceIdList {
		list = append(list, ResourcePermission{
			ResourceId: v,
			RoleId:     roleId,
		})
	}
	return DB.Create(&list).Error
}
