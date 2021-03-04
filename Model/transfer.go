package Model

import (
	"fmt"
	"time"
)

//定义转账的JSON参数结构体。日志记录表
// JSON 用来为web提供返回页面
// db 用来在数据库中提取传递数据
type TransferModel struct{
	Tid int `db:"tid"`  //交易id号
	From string `json:"from" db:"from"`
	To string  `json:"to" db:"to"`
	Money int `json:"m" db:"money"`
	//-----------转账重发--------------
	Status int `json:"-" db:"status"`
	UpdateTime time.Time `json:"-" db:"updatetime"`
	Isback []uint8 `json:"-" db:"isback"`
	//重点，如果isback在数据库中是 bit类型，那么映射到结构体当中的字段类型就是 []uint8
}

func NewTransfer() *TransferModel{
	return &TransferModel{}
}

func (this *TransferModel) String()string{
	return fmt.Sprintf("%s转给%s,金额:%d\n",this.From,this.To,this.Money)
}



