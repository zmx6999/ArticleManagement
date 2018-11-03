package controllers

import (
	"github.com/astaxie/beego"
	"crypto/sha256"
	"encoding/hex"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) HandleError(msg string,tpl string)  {
	this.Data["errmsg"]=msg
	this.TplName=tpl
}

func (this *BaseController) Sha256Str(x string) string {
	y:=sha256.Sum256([]byte(x))
	return hex.EncodeToString(y[:])
}

func (this *BaseController) ShowLayout()  {
	username:=this.GetSession("username")
	this.Data["username"]=username.(string)
	this.Layout="layout.html"
}