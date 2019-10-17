package middleware

import (
	"errors"
	"regexp"
	"strings"
	"sync"
)

var (
	ErrRoleExist    = errors.New("Role has already existed")
	ErrRoleNotExist = errors.New("Role does not exist")
	// ErrRuleExist = errors.New("Rule has already existed")
	ErrRuleNotExist = errors.New("Rule does not exist")
)

var (
	RBAC *Rbac
)

// rbac
type Rbac struct {
	mutex sync.RWMutex
	roles map[int]*StdRole
}

// 角色
type StdRole struct {
	mutex sync.RWMutex
	rules map[int]*StdRule
}

type StdRule struct {
	object string
	method string
}

type Roles map[int]*StdRole
type Rules map[int]*StdRule

func RbacNew() *Rbac {
	RBAC = &Rbac{
		roles: make(Roles),
	}
	return RBAC
}

func StdRoleNew() *StdRole {
	return &StdRole{
		rules: make(Rules),
	}
}

func StdRuleNew(object, method string) *StdRule {
	return &StdRule{
		object: object,
		method: method,
	}
}

func (rbac *Rbac) Add(id int, role *StdRole) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[id]; ok {
		return ErrRoleExist
	}
	rbac.roles[id] = role
	return nil
}

func (rbac *Rbac) RbacCheck(roleId int, obj string, act string) (authPass bool,err error){
	if _,ok := rbac.roles[roleId]; !ok {
		return false,ErrRoleNotExist
	}
	role := rbac.roles[roleId]
	rules := role.rules
	for _,v := range rules{
		r,_ := regexp.Compile(v.object)
		if((v.method == "*" || strings.ToLower(v.method) == strings.ToLower(act)) && (v.object == obj || r.MatchString(obj))){
			return true,nil
		}
	}
	return false,nil
}

func (rbac *Rbac) Remove(id int) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[id]; ok {
		return ErrRoleNotExist
	}
	delete(rbac.roles, id)
	return nil
}

func (rbac *Rbac) Get(id int) (r *StdRole, err error) {
	rbac.mutex.RLock()
	defer rbac.mutex.RUnlock()
	if _, ok := rbac.roles[id]; !ok {
		err = ErrRoleNotExist
	} else {
		r = rbac.roles[id]
	}
	return
}

func (role *StdRole) Add(id int, rule *StdRule) error {
	role.mutex.Lock()
	defer role.mutex.Unlock()
	role.rules[id] = rule
	return nil
}

func (role *StdRole) Remove(id int) error {
	role.mutex.Lock()
	defer role.mutex.Unlock()
	if _, ok := role.rules[id]; ok {
		return ErrRuleNotExist
	}
	delete(role.rules, id)
	return nil
}

func (role *StdRole) Get(id int) (r *StdRule, err error) {
	role.mutex.RLock()
	defer role.mutex.RUnlock()
	if _, ok := role.rules[id]; !ok {
		err = ErrRuleNotExist
	} else {
		r = role.rules[id]
	}
	return
}
