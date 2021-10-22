package db
/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresSQL struct {
	*gorm.DB
}

// connect postgres server
func (p *PostgresSQL) Connect(opt *Options) *gorm.DB {
	dsn := fmt.Sprintf("host=%s:%d user=%s dbname=%s sslmode=disable password=%s", opt.Host, opt.Port, opt.User, opt.DbName, opt.Password)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("connect postgres server was faield:%s", err.Error()))
	}
	p.DB = db
	return db
}
