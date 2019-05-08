package controllers

import (
	"en2dic/models"
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
)

type BaseController struct {
	beego.Controller
	ctrName    string
	actionName string
	curUser    models.User
}

type BackData struct {
	Ret  int         `json:"ret"`
	Data interface{} `json:"data"`
}

type Heard struct{
	Alg string `json:"alg"`
	Type string `json:"type"`
}

type Auth struct{
	Heard Heard
	Token string
	Signature string
}

func (this *BaseController) Prepare() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	authString := this.Ctx.Input.Header("Authorization")
	kv := strings.Split(authString, ".")
	errStr := "auth string invalid"
	if len(kv) != 2 || kv[0] != "Bearer" {
		BackData := &BackData{}
		BackData.Ret = 203
		BackData.Data = errStr
		this.Data["json"] = BackData
		this.ServeJSON()
		return
	}
	tokenString := kv[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 必要的验证 RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		//// 可选项验证  'aud' claim
		//aud := "https://api.cn.atomintl.com"
		//checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
		//if !checkAud {
		//  return token, errors.New("Invalid audience.")
		//}
		// 必要的验证 'iss' claim
		iss := "http://localhost:8080/"
		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
		if !checkIss {
			return token, errors.New("Invalid issuer.")
		}
		/* 我们的公钥,可以在<a href="https://manage.auth0.com/"
		target="_blank">https://manage.auth0.com/</a> 上下载到对应的封装好的json，里面包括了签名
		*/
		k5c := "xxxx"
		cert := "-----BEGIN CERTIFICATE-----\n" + k5c + "\n-----END CERTIFICATE-----"
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		//result := []byte(cert) // 不是正确的 PUBKEY 格式 都会 报  key is of invalid type
		return result, errors.New("key is of invalid type")
	})
	if err != nil {
		logs.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				errStr = "key is of invalid type"
				this.ServeJSON()
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				errStr = "oken is either expired or not active yet"
				this.Data["json"] = errStr
				this.ServeJSON()
				return
			} else {
				errStr = "Couldn't handle this token"
				this.Data["json"] = errStr
				this.ServeJSON()
				return
			}
		} else {
			errStr = "Couldn't handle this token"
			this.Data["json"] = errStr
			this.ServeJSON()
			return
		}
	}
	if !token.Valid {
		logs.Error("Token invalid:", tokenString)
		errStr = "Token invalid"
		this.Data["json"] = errStr
		this.ServeJSON()
		return
	}
}
