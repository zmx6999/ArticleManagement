package controllers

import (
	"181103/models"
	"github.com/astaxie/beego/orm"
	"encoding/base64"
)

type UserController struct {
	BaseController
}

func (this *UserController) ShowRegister() {
	this.TplName="user/register.html"
}

func (this *UserController) HandleRegister()  {
	username:=this.GetString("username")
	password:=this.GetString("password")
	if username=="" || password=="" {
		this.HandleError("username and password cannot be empty","user/register.html")
		return
	}

	var user models.User
	user.Username=username
	user.Password=this.Sha256Str(password)
	o:=orm.NewOrm()
	if _,err:=o.Insert(&user);err!=nil {
		this.HandleError(err.Error(),"user/register.html")
		return
	}

	this.Redirect("/login",302)
}

func (this *UserController) ShowLogin() {
	username:=this.Ctx.GetCookie("username")
	if username!="" {
		tmp,_:=base64.StdEncoding.DecodeString(username)
		this.Data["username"]=string(tmp)
		this.Data["remember"]="checked"
	}
	this.TplName="user/login.html"
}

func (this *UserController) HandleLogin()  {
	username:=this.GetString("username")
	password:=this.GetString("password")
	if username=="" || password=="" {
		this.HandleError("username and password cannot be empty","user/login.html")
		return
	}

	var user models.User
	user.Username=username
	o:=orm.NewOrm()
	if err:=o.Read(&user,"Username");err!=nil {
		this.HandleError("error username","user/login.html")
		return
	}
	if user.Password!=this.Sha256Str(password) {
		this.HandleError("error password","user/login.html")
		return
	}

	remember:=this.GetString("remember")
	if remember!="" {
		this.Ctx.SetCookie("username",base64.StdEncoding.EncodeToString([]byte(username)),100)
	} else {
		this.Ctx.SetCookie("username","",-1)
	}
	this.SetSession("username",username)

	this.Redirect("/article/list",302)
}

func (this *UserController) Logout()  {
	this.DelSession("username")
	this.Redirect("/login",302)
}