package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/mojocn/base64Captcha"
)

const (
	Digit = iota
	Character
	Audio
)

type Captcha struct {
	Id string 		`json:"id"`
	CaptchaType	 int	`json:"captcha_type"`
}

type ConfigCaptcha struct {
	Captcha
	VerifyValue     string
	ConfigAudio     base64Captcha.ConfigAudio
	ConfigCharacter base64Captcha.ConfigCharacter
	ConfigDigit     base64Captcha.ConfigDigit
}

func PassEncrypt(pass string, salt string) string{
	h := md5.New()
	var buffer bytes.Buffer
	buffer.WriteString(pass)
	buffer.WriteString(salt)
	md5key := buffer.String()
	h.Write([]byte(md5key))
	cipherStr :=h.Sum(nil)
	password := hex.EncodeToString(cipherStr)
	return password
}

// GenerateCaptcha - generate captcha by config,types include (audio,images)
func GenerateCaptcha(params ConfigCaptcha) (string,string){
	var config interface{}
	switch params.CaptchaType {
	case Audio:
		config = params.ConfigAudio
	case Character:
		config = params.ConfigCharacter
	case Digit:
		config = params.ConfigDigit
	}
	idKey, Cap := base64Captcha.GenerateCaptcha(params.Id,config)
	base64Data := base64Captcha.CaptchaWriteToBase64Encoding(Cap)
	return idKey,base64Data
}

// VerifyCaptchaCode - verify captcha code
func VerifyCaptchaCode(idkey string,captchaCode string) bool {
	verifyResult := base64Captcha.VerifyCaptcha(idkey, captchaCode)
	return verifyResult
}