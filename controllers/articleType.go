package controllers

import (
	"181103/models"
	"github.com/astaxie/beego/orm"
)

type ArticleTypeController struct {
	BaseController
}

func (this *ArticleTypeController) ShowAddType()  {
	var types []models.ArticleType
	o:=orm.NewOrm()
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types
	this.ShowLayout()
	this.TplName="articleType/addType.html"
}

func (this *ArticleTypeController) HandleAddType() {
	typename:=this.GetString("typename")
	if typename=="" {
		this.ShowLayout()
		this.HandleError("","articleType/addType.html")
		return
	}

	var articleType models.ArticleType
	articleType.TypeName=typename
	o:=orm.NewOrm()
	if _,err:=o.Insert(&articleType);err!=nil {
		this.ShowLayout()
		this.HandleError(err.Error(),"articleType/addType.html")
		return
	}

	this.Redirect("/article_type/add",302)
}

func (this *ArticleTypeController) DeleteType()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/article_type/add",302)
		return
	}

	var articleType models.ArticleType
	articleType.Id=id
	o:=orm.NewOrm()
	o.Delete(&articleType)
	this.Redirect("/article_type/add",302)
}