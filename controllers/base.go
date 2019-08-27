package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"go-admin-api/middleware"
	"go-admin-api/models"
	"reflect"
	"strings"
)

const(
	MSG_OK 	= 200			// 请求成功码
	MSG_NO_LOGIN = 400  	// 未登录
	MSG_ERR = 401			// 请求公共错误码
)

type BaseController struct {
	beego.Controller
	controllerName 	string
	actionName 		string
	currentRoute 	string
	Url				string
	Uri				string
	userId  		int64
	userName 		string
	nickName      	string
	user 			*models.User
}

//output result struct
type outputData struct {
	code int64
	msg interface{}
	data interface{}
}

//Prepare is a pre-processing of all method requests
func (this *BaseController) Prepare() {
	controllerName,actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName)
	this.actionName = strings.ToLower(actionName)
	this.currentRoute = this.controllerName + "/" + this.actionName
	this.Data["controllerName"] = this.controllerName;
	this.Data["actionName"] = this.actionName;
	this.Data["currentRoute"] = this.currentRoute;
	this.Data["url"] = this.Ctx.Input.URL()
	this.Data["Uri"] = this.Ctx.Input.URI()

	//if action is in whiteList , not need authority authentication
	if(!this.whiteUrlAuthentication()){
		this.authorityAuthentication()
	}
}

//authority Authentication
func (this *BaseController) authorityAuthentication() bool {
	token := this.Ctx.GetCookie("X-Token")
	fmt.Println(token)
	jwtkey := beego.AppConfig.String("jwtkey")
	params, ok := middleware.JWTParseToken(token, jwtkey)
	user := params.(jwt.MapClaims)
	if( !ok ){
		this.apiError(MSG_NO_LOGIN,"未登录,请线登录",nil)
	}
	fmt.Println(reflect.TypeOf(user["id"]))
	this.userId = int64(user["id"].(float64))
	this.nickName = user["nickName"].(string)
	this.userName = user["name"].(string)
	return true
}

func (this *BaseController) whiteUrlAuthentication() bool {
	var whiteUrlList = []string{
		"/user/login",
		"/user/getcaptcha",
	}
	for _,value := range whiteUrlList{
		if this.Data["url"] == value {
			return true
		}
	}
	return false
}

func (this *BaseController) apiError(code int64,msg interface{},data interface{}){
	output := make(map[string]interface{})
	output["code"] = code
	output["msg"] = msg
	output["data"] = data
	this.Data["json"] = output
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) apiSuc(data interface{}){
	output := make(map[string]interface{})
	output["code"] = MSG_OK
	output["msg"] = "请求成功"
	output["data"] = data
	this.Data["json"] = output
	this.ServeJSON()
}