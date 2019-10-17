package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"go-admin-api/middleware"
	"go-common/app/admin/main/macross/model/errors"
)

const (
	VALID  = 1
	ACTIVE = 1
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
	var roles []Role
	var roleRights []RoleRight
	var rights []Right

	var mapRights map[int]Right
	o := orm.NewOrm()
	_, roleErr := o.QueryTable(new(Role)).Filter("valid", VALID).Filter("active", ACTIVE).Limit(-1).All(&roles)

	if roleErr != nil {
		return errors.New("query table Role error:", roleErr)
	}
	for _, v := range roles {
		stdRole := middleware.StdRoleNew()
		Rbac.Add(int(v.Id), stdRole)
	}
	_, errRoleRight := o.QueryTable(new(RoleRight)).Filter("valid", VALID).Filter("active", ACTIVE).Limit(-1).All(&roleRights)
	if errRoleRight != nil {
		return errors.New("query table RoleRight error:", errRoleRight)
	}
	_, errRight := o.QueryTable(new(Right)).Filter("valid", VALID).Filter("active", ACTIVE).Limit(-1).All(&rights)
	if errRight != nil {
		return errors.New("query table Right error:", errRight)
	}
	mapRights = make(map[int]Right)
	for _, v := range rights {
		mapRights[int(v.Id)] = v
	}

	for _, v := range roleRights {
		role, err := Rbac.Get(int(v.RoleId))
		if err != nil {
			return errors.New("rbac get role error", err)
		}
		rightId := int(int(v.RightId))
		role.Add(rightId, middleware.StdRuleNew(mapRights[rightId].Object, mapRights[rightId].Action))
	}
	return nil
}
