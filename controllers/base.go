package controllers

import (
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"go-admin-api/middleware"
	"go-admin-api/models"
	"strings"
)

const (
	MSG_OK       = 200 // 请求成功码
	MSG_NO_LOGIN = 400 // 未登录
	MSG_ERR      = 401 // 请求公共错误码
	MSG_NO_AUTH  = 402 // 接口未授权
	MSG_IP_LIMIT = 403 // IP限制
)

type BaseController struct {
	beego.Controller
	controllerName string
	actionName     string
	currentRoute   string
	Url            string
	Uri            string
	userId         int64
	userName       string
	nickName       string
	roleIds        []int
	user           *models.User
}

//output result struct
type outputData struct {
	code int64
	msg  interface{}
	data interface{}
}

//Prepare is a pre-processing of all method requests
func (this *BaseController) Prepare() {
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName)
	this.actionName = strings.ToLower(actionName)
	this.currentRoute = this.controllerName + "/" + this.actionName
	this.Data["controllerName"] = this.controllerName
	this.Data["actionName"] = this.actionName
	this.Data["currentRoute"] = this.currentRoute
	this.Data["url"] = this.Ctx.Input.URL()
	this.Data["uri"] = this.Ctx.Input.URI()

	//if action is in whiteList , not need authority authentication
	specialPass := this.isSpecialAuthCheck()
	if !specialPass {
		authPass := this.authorityAuthentication()
		if !authPass {
			this.apiError(MSG_NO_AUTH, "接口未授权", nil)
		}
	}
}

//authority Authentication
func (this *BaseController) authorityAuthentication() bool {
	var roleIds []int
	token := this.Ctx.GetCookie("X-Token")
	jwtkey := beego.AppConfig.String("jwtkey")
	params, ok := middleware.JWTParseToken(token, jwtkey)
	user := params.(jwt.MapClaims)
	if !ok {
		this.apiError(MSG_NO_LOGIN, "未登录,请线登录", nil)
	}
	this.userId = int64(user["id"].(float64))
	this.nickName = user["nickName"].(string)
	this.userName = user["name"].(string)
	for _, v := range user["roleIds"].([]interface{}) {
		roleIds = append(roleIds, int(v.(float64)))
	}
	this.roleIds = roleIds
	for _, v := range roleIds {
		// 对超级管理员不进行权限验证
		admin, _ := beego.AppConfig.Int("admin")
		if v == admin {
			return true
		}
		if ok, _ = middleware.RBAC.RbacCheck(v, this.Data["uri"].(string), this.Ctx.Input.Method()); ok {
			return true
		}
	}
	return false
}

// Special urls are determined in special ways
func (this *BaseController) isSpecialAuthCheck() bool {
	URL := this.Data["url"].(string)
	specialURL := this.specialUrl()
	if _, ok := specialURL[URL]; ok {
		checkWays := specialURL[URL]
		//匿名方法，不需要进行权限检测
		if _, ok := checkWays["anonymous"]; ok {
			return true
		}

		//ip限制下的方法检测
		if _, ok := checkWays["ipLimit"]; ok {
			ipLimit := beego.AppConfig.String("ipLimit")
			ipLimitArray := strings.Split(ipLimit, ",")
			for _, v := range ipLimitArray {
				if v == this.Ctx.Input.IP() {
					return true
				}
			}
			this.apiError(MSG_IP_LIMIT, "接口存在IP限制", nil)
		}
	}
	return false
}

func (this *BaseController) specialUrl() map[string]map[string]bool {
	specialUrl := map[string]map[string]bool{
		"/user/login":      {"anonymous": true},
		"/user/getcaptcha": {"anonymous": true},
		"/user/logout":     {"anonymous": true},
	}
	return specialUrl
}

func (this *BaseController) apiError(code int64, msg interface{}, data interface{}) {
	output := make(map[string]interface{})
	output["code"] = code
	output["msg"] = msg
	output["data"] = data
	this.Data["json"] = output
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) apiSuc(data interface{}) {
	output := make(map[string]interface{})
	output["code"] = MSG_OK
	output["msg"] = "请求成功"
	output["data"] = data
	this.Data["json"] = output
	this.ServeJSON()
}
