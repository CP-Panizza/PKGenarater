package PKGenarater

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"log"
)

//自定义数据库主键生成

/**
idLength:主键生成的长度
table_name:表名，数据库内表的真实名字
pk:表的主键字段
sign: user00000012 中的user
*/
type PKGenarater struct {
	DB        *sql.DB `inject:"db"`
	TableName string
	IdLength  int
	PK        string
	Sign      string

	flag       chan int
	primaryKey string
	firstGet   bool
}

func NewPKGenarater(tableName string, idLength int, PK string, sign string) *PKGenarater {
	pk := &PKGenarater{
		TableName: tableName,
		IdLength:  idLength,
		PK:        PK,
		Sign:      sign,
		flag:      make(chan int, 1),
	}
	pk.flag <- 0
	return pk
}

func (this *PKGenarater) Init() error {
	jugeNullSql := `select count(*) from ` + this.TableName
	var count int
	row := this.DB.QueryRow(jugeNullSql)
	if err := row.Scan(&count); err != nil {
		return err
	}
	//第一次查询，表的主键为空时
	if count == 0 {
		this.firstGet = true
	} else {
		str := strings.Builder{}
		str.WriteString("select max(a.")
		str.WriteString(this.PK)
		str.WriteString(") from ")
		str.WriteString(this.TableName)
		str.WriteString(" as a where a.")
		str.WriteString(this.PK)
		str.WriteString(" like '")
		str.WriteString(this.Sign)
		str.WriteString("%'")
		sql := str.String()
		row := this.DB.QueryRow(sql)
		var result string = ""
		if err := row.Scan(&result); err != nil {
			return err
		}
		this.primaryKey = result
		this.firstGet = false
	}
	log.Println("Success init table '", this.TableName, "' PKGenarater!")
	return nil
}

//获取执行权
func (this *PKGenarater) _getKey() {
	<-this.flag
}

//返回执行权
func (this *PKGenarater) _returnKey() {
	this.flag <- 0
}


func (this *PKGenarater) GetPK() (string, error) {
	this._getKey()
	defer this._returnKey()
	if this.firstGet {
		gpk := this.Sign
		for i := 0; i < (this.IdLength - len(this.Sign) - 1); i++ {
			gpk += "0"
		}
		gpk += "1"
		fmt.Println(gpk)
		this.primaryKey = gpk
		this.firstGet = false
		return gpk, nil
	}
	result := this.primaryKey
	fmt.Println("last val:", result)
	num, err := strconv.Atoi(result[len(this.Sign):])
	if err != nil {
		return "", err
	}
	nextNumber := num + 1
	nextNumberStr := strconv.Itoa(nextNumber)
	gpk := this.Sign
	for i := 0; i < (this.IdLength - len(this.Sign) - len(nextNumberStr)); i++ {
		gpk += "0"
	}
	gpk += nextNumberStr
	this.primaryKey = gpk
	fmt.Println("next val:", gpk)
	return gpk, nil
}


