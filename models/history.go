package models

import (
	_ "github.com/astaxie/beego/orm"
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
)

func Get(since time.Time, until time.Time) {
	type Source struct {
		DbMax string
		UserMax string
	}

	var sl []Source
	count, err := orm.NewOrm().Raw("SELECT db_max, user_max FROM mysql_slow_query_review_history").QueryRows(&sl)
	if err == nil {
		fmt.Println("Count: ", count)
	} else {
		fmt.Println(err.Error())
	}
}

