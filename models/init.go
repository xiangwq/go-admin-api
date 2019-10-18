package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"go-admin-api/middleware"
	"go-common/app/admin/main/macross/model/errors"
)

const (
	VALID = 0
)

func Init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpass")
	dbname := beego.AppConfig.String("dbname")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dsn)
}

func MysqlDsn() string {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpass")
	dbname := beego.AppConfig.String("dbname")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	return dsn
}

func RegisterRbac() error {
	Rbac := middleware.RbacNew()
	var roles []AdminRole
	var rolePermissions []AdminRolePermission
	var permissions []AdminPermission

	var mapPermissions map[int]AdminPermission
	o := orm.NewOrm()
	_, roleErr := o.QueryTable(new(AdminRole)).Filter("valid", VALID).Limit(-1).All(&roles)

	if roleErr != nil {
		return errors.New("query table Role error:", roleErr)
	}
	for _, v := range roles {
		stdRole := middleware.StdRoleNew()
		Rbac.Add(int(v.Id), stdRole)
	}
	_, errRolePermission := o.QueryTable(new(AdminRolePermission)).Filter("valid", VALID).Limit(-1).All(&rolePermissions)
	if errRolePermission != nil {
		return errors.New("query table RoleRight error:", errRolePermission)
	}
	_, errRight := o.QueryTable(new(AdminPermission)).Filter("valid", VALID).Filter("type", Interface).Limit(-1).All(&permissions)
	if errRight != nil {
		return errors.New("query table Right error:", errRight)
	}
	mapPermissions = make(map[int]AdminPermission)
	for _, v := range mapPermissions {
		mapPermissions[int(v.Id)] = v
	}

	for _, v := range rolePermissions {
		role, err := Rbac.Get(int(v.RoleId))
		if err != nil {
			return errors.New("rbac get role error", err)
		}
		rightId := int(int(v.PermissionId))
		role.Add(rightId, middleware.StdRuleNew(mapPermissions[rightId].Url, mapPermissions[rightId].Method))
	}
	return nil
}
