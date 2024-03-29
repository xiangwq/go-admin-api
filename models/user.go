package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id         int64
	Name       string `orm:"size(128)"`
	NickName   string `orm:"size(128)"`
	AvatorPath string `orm:"size(128)"`
	IdCardNo   string `orm:"size(128)"`
	Phone      string `orm:"size(128)"`
	Password   string `orm:"size(128)"`
	Email      string `orm:"size(128)"`
	Sex        int8
	LoginIp    int
	Lock       int8
	Valid      int8
	Ctime  	   time.Time
	Utime	   time.Time
}


func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.QueryTable(new(User)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}


func GetUserByPhone(phone string) (v *User,err error){
	o := orm.NewOrm()
	v = &User{Phone: phone}
	if err = o.QueryTable(new(User)).Filter("phone", phone).One(v); err == nil {
		return v, nil
	}
	return nil, err
}


// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}


// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
