package controllers

import (
	"181103/models"
	"github.com/astaxie/beego/orm"
	"math"
	"errors"
	"strings"
	"path"
	"time"
)

type ArticleController struct {
	BaseController
}

func (this *ArticleController) ArticleList()  {
	var types []models.ArticleType
	o:=orm.NewOrm()
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types

	s:=o.QueryTable("Article")
	typename:=this.GetString("select")
	if typename=="" {
		s=s.RelatedSel("ArticleType")
	} else {
		s=s.RelatedSel("ArticleType").Filter("ArticleType__TypeName",typename)
	}

	totalRows,_:=s.Count()
	pageSize:=2
	pageCount:=int(math.Ceil(float64(totalRows)/float64(pageSize)))
	page,err:=this.GetInt("p")
	if err!=nil {
		page=1
	}
	if page<1 {
		page=1
	}
	if page>pageCount {
		page=pageCount
	}
	var articles []models.Article
	s.Limit(pageSize,pageSize*(page-1)).All(&articles)
	this.Data["totalRows"]=totalRows
	this.Data["pageCount"]=pageCount
	this.Data["page"]=page
	this.Data["data"]=articles
	this.Data["typename"]=typename
	this.ShowLayout()
	this.TplName="article/index.html"
}

func (this *ArticleController) ArticleDetail()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/article/list",302)
		return
	}
	var article models.Article
	o:=orm.NewOrm()
	err=o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id",id).One(&article)
	if err!=nil {
		this.Redirect("/article/list",302)
		return
	}
	article.VisitNum+=1
	o.Update(&article)
	this.Data["data"]=article

	m2m:=o.QueryM2M(&article,"Users")
	var user models.User
	username:=this.GetSession("username")
	user.Username=username.(string)
	o.Read(&user,"Username")
	m2m.Add(user)

	var users []models.User
	o.QueryTable("User").Filter("Articles__Article__Id",id).Distinct().All(&users)
	this.Data["users"]=users
	this.ShowLayout()
	this.TplName="article/content.html"
}

func (this *ArticleController) uploadImage(from string) (string,error) {
	file,head,err:=this.GetFile(from)
	if err!=nil {
		return "",err
	}
	defer file.Close()
	if head.Filename=="" {
		return "",nil
	}
	if head.Size>1024*1024*5 {
		return "",errors.New("image<5M")
	}
	ext:=strings.ToLower(path.Ext(head.Filename))
	if ext!=".jpg" && ext!=".png" && ext!=".jpeg" {
		return "",errors.New("image type error")
	}
	loc,_:=time.LoadLocation("Asia/Saigon")
	filename:="./static/img/"+time.Now().In(loc).Format("20060102150405")+ext
	err=this.SaveToFile(from,filename)
	if err!=nil {
		return "",err
	}
	return strings.TrimPrefix(filename,"."),nil
}

func (this *ArticleController) ShowAddArticle()  {
	var types []models.ArticleType
	o:=orm.NewOrm()
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types
	this.ShowLayout()
	this.TplName="article/add.html"
}

func (this *ArticleController) HandleAddArticle()  {
	var types []models.ArticleType
	o:=orm.NewOrm()
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"]=types

	title:=this.GetString("title")
	content:=this.GetString("content")
	if title=="" || content=="" {
		this.ShowLayout()
		this.HandleError("title and content cannot be empty","article/add.html")
		return
	}

	filename,err:=this.uploadImage("uploadname")
	if err!=nil {
		this.ShowLayout()
		this.HandleError(err.Error(),"article/add.html")
		return
	}
	if filename=="" {
		this.ShowLayout()
		this.HandleError("please upload image","article/add.html")
		return
	}

	var article models.Article
	article.Title=title
	article.Content=content
	article.Img=filename

	typename:=this.GetString("typename")
	var articleType models.ArticleType
	articleType.TypeName=typename
	o.Read(&articleType,"TypeName")
	article.ArticleType=&articleType

	loc,_:=time.LoadLocation("Asia/Saigon")
	article.Addtime=time.Now().In(loc)

	if _,err=o.Insert(&article);err!=nil {
		this.ShowLayout()
		this.HandleError("failed to publish article","article/add.html")
		return
	}
	this.Redirect("/article/list",302)
}

func (this *ArticleController) ShowUpdateArticle() {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/article/list",302)
		return
	}
	var article models.Article
	article.Id=id
	o:=orm.NewOrm()
	err=o.Read(&article)
	if err!=nil {
		this.Redirect("/article/list",302)
		return
	}
	this.Data["data"]=article
	this.ShowLayout()
	this.TplName="article/update.html"
}

func (this *ArticleController) HandleUpdateArticle()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/article/list",302)
		return
	}
	var article models.Article
	article.Id=id
	o:=orm.NewOrm()
	err=o.Read(&article)
	if err!=nil {
		this.Redirect("/article/list",302)
		return
	}
	this.Data["data"]=article

	title:=this.GetString("title")
	content:=this.GetString("content")
	if title=="" || content=="" {
		this.ShowLayout()
		this.HandleError("title and content cannot be empty","article/update.html")
		return
	}

	filename,err:=this.uploadImage("uploadname")
	if err!=nil {
		this.ShowLayout()
		this.HandleError(err.Error(),"article/update.html")
		return
	}

	article.Title=title
	article.Content=content
	if filename!="" {
		article.Img=filename
	}

	loc,_:=time.LoadLocation("Asia/Saigon")
	article.Addtime=time.Now().In(loc)

	if _,err=o.Update(&article);err!=nil {
		this.ShowLayout()
		this.HandleError("failed to edit article","article/update.html")
		return
	}
	this.Redirect("/article/list",302)
}

func (this *ArticleController) DeleteArticle()  {
	id,err:=this.GetInt("id")
	if err!=nil {
		this.Redirect("/article/list",302)
		return
	}

	var article models.Article
	article.Id=id
	o:=orm.NewOrm()
	o.Delete(&article)
	this.Redirect("/article/list",302)
}