package db
/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

type Mssql struct {
	*gorm.DB
}

func (m *Mssql) Connect(opt *Options) *gorm.DB {
	dsn := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s", opt.Host, opt.Port, opt.DbName, opt.User, opt.Password)
	db, err := gorm.Open("mssql", dsn)
	if err != nil {
		panic(fmt.Sprintf("connect mssql server was faield:%s", err.Error()))
	}
	m.DB = db
	return db
}
