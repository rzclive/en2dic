package models

import (
	"github.com/astaxie/beego/orm"
)

var (
	WordList map[int]*Word
)

type Word struct {
	Id int    `json:"id" orm:"id"`
	Zh string `json:"zh_cn" orm:"zh_cn"`
	En string `json:"en" orm:"en"`
}

func init() {
	orm.RegisterModel(new(Word))
}

func (this *Word) MutileInsert() {
	o := orm.NewOrm()
	o.Using("")

}
