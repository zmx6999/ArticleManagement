package main

import (
	_ "181103/routers"
	"github.com/astaxie/beego"
	"time"
)

func main() {
	beego.AddFuncMap("prepage",ShowPrePage)
	beego.AddFuncMap("nextpage",ShowNextPage)
	beego.AddFuncMap("showtime",ShowTime)
	beego.Run()
}

func ShowPrePage(page int) int {
	if page==1 {
		return page
	} else {
		return page-1
	}
}

func ShowNextPage(page int,pageCount int) int {
	if page==pageCount {
		return page
	} else {
		return page+1
	}
}

func ShowTime(t time.Time,timezone string) string {
	loc,_:=time.LoadLocation(timezone)
	return t.In(loc).Format("2006-01-02 15:04:05")
}