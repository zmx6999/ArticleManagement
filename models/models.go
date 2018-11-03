package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id int `orm:"pk;auto"`
	Username string `orm:"size(20)"`
	Password string `orm:"size(64)"`

	Articles []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id int `orm:"pk;auto"`
	Title string `orm:"size(100)"`
	Content string `orm:"size(600)"`
	Img string `orm:"size(100)"`
	VisitNum int `orm:"default(0);null"`
	Addtime time.Time

	ArticleType *ArticleType `orm:"rel(fk);on_delete(set_null);null"`
	Users []*User `orm:"rel(m2m)"`
}

type ArticleType struct {
	Id int `orm:"pk;auto"`
	TypeName string `orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"`
}

func init()  {
	orm.RegisterDataBase("default","mysql","root:123456@tcp(localhost:3306)/db181103?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}