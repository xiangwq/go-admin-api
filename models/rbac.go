package models

import (
	"github.com/astaxie/beego/orm"
	"go-common/app/admin/main/macross/model/errors"
	"time"
)

type AdminRole struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Valid       int8      `json:"valid"`
	Ctime       time.Time `json:"ctime"`
}

type AdminPermission struct {
	Id          int64     `json:"id"`
	Type        int8      `json:"type"`
	ParentId    int64     `json:"parent_id"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	Method      string    `json:"method"`
	Icon        string    `json:"icon"`
	Title       string    `json:"title"`
	Valid       int8      `json:"valid"`
	CreatedUid  int64     `json:"created_uid"`
	Ctime       time.Time `json:"ctime"`
	Utime       time.Time `json:"utime"`
}

type AdminUserRole struct {
	Id     int64 `json:"id"`
	UserId int64 `json:"user_id"`
	RoleId int64 `json:"role_id"`
	Valid  int8  `json:"valid"`
}

type AdminRolePermission struct {
	Id           int64 `json:"id"`
	RoleId       int64 `json:"role_id"`
	PermissionId int64 `json:"permission_id"`
	Valid        int8  `json:"valid"`
}

func init() {
	orm.RegisterModel(new(AdminRole))
	orm.RegisterModel(new(AdminPermission))
	orm.RegisterModel(new(AdminUserRole))
	orm.RegisterModel(new(AdminRolePermission))
}

const (
	Menu = iota
	Interface
	Data
)

func GetRolesByUserId(uid int64) (v []*AdminUserRole, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable(new(AdminUserRole)).Filter("user_id", uid).Filter("valid", VALID).Limit(-1).All(&v)
	if err != nil {
		return nil, errors.New("get roles error:", err)
	}
	return v, nil
}
