package db
/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Mysql struct {
	*gorm.DB
}

type SqlClient interface {
	Connect(opt *Options) *gorm.DB
}

type Options struct {
	Type     string
	User     string
	Password string
	Host     string
	DbName   string
	Port     int
	Prefix   string
	Debug    bool
}

// connect mysql server
func (m *Mysql) Connect(opt *Options) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", opt.User, opt.Password, opt.Host, opt.Port, opt.DbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("connect mysql server was faield:%s", err.Error()))
	}
	m.DB = db
	return db
}
