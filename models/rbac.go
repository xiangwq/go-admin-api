package models

import (
	"github.com/astaxie/beego/orm"
	"go-common/app/admin/main/macross/model/errors"
	"time"
)

type Role struct {
	Id 			int64
	Name 		string
	Description string
	Valid 		int8
	Active 		int8
	Ctime 		time.Time
}

type Right struct {
	Id 			int64
	Type 		int8
	Description string
	Object 		string
	Action 		string
	Valid 		int8
	Active 		int8
	Ctime 		time.Time
}

type UserRole struct {
	Id 			int64
	Uid  		int64
	RoleId 		int64
	Active 		int8
	Valid 		int8
}

type RoleRight struct {
	Id 			int64
	RoleId		int64
	RightId		int64
	Active 		int8
	Valid 		int8
}

func init() {
	orm.RegisterModel(new(Role))
	orm.RegisterModel(new(Right))
	orm.RegisterModel(new(UserRole))
	orm.RegisterModel(new(RoleRight))
}

const (
	Menu =	iota
	Interface
)

func GetRolesByUserId(uid int64) (v []*UserRole,err error){
	o := orm.NewOrm()
	_, err =  o.QueryTable(new(UserRole)).Filter("uid",uid).Filter("active",ACTIVE).Filter("valid",VALID).Limit(-1).All(&v)
	if err != nil{
		return nil,errors.New("get roles error:",err)
	}
	return v,nil
}
