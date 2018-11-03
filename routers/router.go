package routers

import (
	"github.com/astaxie/beego"
	"181103/controllers"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeExec,Filter)
	beego.InsertFilter("/article_type/*",beego.BeforeExec,Filter)

    beego.Router("/register",&controllers.UserController{},"get:ShowRegister;post:HandleRegister")
	beego.Router("/login",&controllers.UserController{},"get:ShowLogin;post:HandleLogin")
	beego.Router("/logout",&controllers.UserController{},"get:Logout")

	beego.Router("/article/list",&controllers.ArticleController{},"get:ArticleList")
	beego.Router("/article/detail",&controllers.ArticleController{},"get:ArticleDetail")
	beego.Router("/article/add",&controllers.ArticleController{},"get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/article/update",&controllers.ArticleController{},"get:ShowUpdateArticle;post:HandleUpdateArticle")
	beego.Router("/article/delete",&controllers.ArticleController{},"get:DeleteArticle")

    beego.Router("/article_type/add",&controllers.ArticleTypeController{},"get:ShowAddType;post:HandleAddType")
	beego.Router("/article_type/delete",&controllers.ArticleTypeController{},"get:DeleteType")
}

var Filter = func(ctx *context.Context) {
	username:=ctx.Input.Session("username")
	if username==nil {
		ctx.Redirect(302,"/login")
	}
}