package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/mojocn/base64Captcha"
	"go-admin-api/middleware"
	"go-admin-api/models"
)

type User struct {
}

type Login struct {
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	Idkey      string `json:"idkey"`
	VerifyCode string `json:"verify_code"`
}

//  UserController operations for User
type UserController struct {
	BaseController
}

func (this *UserController) Login() {
	var input Login
	var tokenInfo map[string]interface{}
	tokenInfo = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &input)
	captchaOk := middleware.VerifyCaptchaCode(input.Idkey, input.VerifyCode)
	if !captchaOk {
		this.apiError(MSG_ERR, "验证码错误", nil)
	}
	user, err := models.GetUserByPhone(input.Phone)
	fmt.Println(err)
	if err != nil {
		this.apiError(MSG_ERR, "账号或密码错误", nil)
	}
	password := middleware.PassEncrypt(input.Password, beego.AppConfig.String("loginsalt"))
	fmt.Println(password)
	if user.Password != password {
		this.apiError(MSG_ERR, "账号或密码错误", nil)
	}

	roles, errRole := models.GetRolesByUserId(user.Id)
	if errRole != nil {
		this.apiError(MSG_ERR, "获取角色数据错误或当前账号不存在分配的角色", nil)
	}

	roleIds := []int{}
	for _, v := range roles {
		roleIds = append(roleIds, int(v.RoleId))
	}
	tokenInfo["id"] = user.Id
	tokenInfo["name"] = user.Name
	tokenInfo["nickName"] = user.NickName
	tokenInfo["roleIds"] = roleIds
	token := middleware.JWTGenToken(beego.AppConfig.String("jwtkey"), tokenInfo)
	this.apiSuc(token)
}

func (this *UserController) Logout() {
	this.apiSuc(nil)
}

func (this *UserController) GetCaptcha() {
	var captchaData map[string]interface{}
	var input middleware.Captcha
	var config middleware.ConfigCaptcha
	json.Unmarshal(this.Ctx.Input.RequestBody, &input)
	switch input.CaptchaType {
	case middleware.Audio:
		config.Id = input.Id
		config.CaptchaType = input.CaptchaType
		config.ConfigAudio = base64Captcha.ConfigAudio{
			CaptchaLen: 5,
			Language:   "zh",
		}
	case middleware.Digit:
		config.Id = input.Id
		config.CaptchaType = input.CaptchaType
		config.ConfigDigit = base64Captcha.ConfigDigit{
			Height:     47,
			Width:      160,
			MaxSkew:    0.7,
			DotCount:   80,
			CaptchaLen: 5,
		}
	case middleware.Character:
		config.Id = input.Id
		config.CaptchaType = input.CaptchaType
		config.ConfigCharacter = base64Captcha.ConfigCharacter{
			Height: 47,
			Width:  160,
			//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
			Mode:               base64Captcha.CaptchaModeNumber,
			ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
			ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
			IsShowHollowLine:   false,
			IsShowNoiseDot:     false,
			IsShowNoiseText:    false,
			IsShowSlimeLine:    false,
			IsShowSineLine:     false,
			CaptchaLen:         5,
		}
	}
	idkey, data := middleware.GenerateCaptcha(config)

	captchaData = make(map[string]interface{})
	captchaData["idkey"] = idkey
	captchaData["data"] = data
	this.apiSuc(captchaData)
}

func (this *UserController) Info() {
	var userInfo map[string]interface{}
	user, err := models.GetUserById(this.userId)
	if err != nil {
		this.apiError(MSG_NO_LOGIN, "未登录,请线登录", nil)
	}
	userInfo = make(map[string]interface{})
	userInfo["name"] = user.Name
	userInfo["avator"] = user.AvatorPath
	this.apiSuc(userInfo)
}
