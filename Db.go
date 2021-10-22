package library

/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"sync"
	"crast/library/db"
)

var (
	// db singleton instance
	dbInstance map[string]*gorm.DB
	// db singleton sync state
	dbOnceSync map[string]*sync.Once
)

// init database config struct
func initOption(config string) *db.Options {
	var dbConfig *db.Options
	if _, err := os.Stat("config/database.ini"); err != nil {
		dbConfig = &db.Options{
			Type:     "mysql",
			User:     "root",
			Password: "",
			Host:     "localhost",
			DbName:   "",
			Port:     3306,
			Prefix:   "",
			Debug:    false,
		}
	} else {
		cfg := SourceConfig("database")
		if config == "default" {
			config = "master_database"
		}
		dbConfig = &db.Options{
			Type:     cfg.Get(config + ".type").MustString("mysql"),
			User:     cfg.Get(config + ".user").MustString("root"),
			Password: cfg.Get(config + ".password").MustString(""),
			Host:     cfg.Get(config + ".host").MustString("localhost"),
			DbName:   cfg.Get(config + ".name").MustString(""),
			Port:     cfg.Get(config + ".port").MustInt(3306),
			Prefix:   cfg.Get(config + ".prefix").MustString(""),
			Debug:    false,
		}
		if debug := cfg.Get(config + ".debug").MustInt(0); debug == 1 {
			dbConfig.Debug = true
		}
	}
	return dbConfig
}

// instance db client
func DbInstance(dbConfigs ...string) *gorm.DB {
	configName := "default"
	if len(dbConfigs) > 0 && dbConfigs[0] != "" {
		configName = dbConfigs[0]
	}
	// init instance and sync
	if dbInstance == nil {
		dbInstance = make(map[string]*gorm.DB)
		dbOnceSync = make(map[string]*sync.Once)
	}

	if client, ok := dbInstance[configName]; !ok || client == nil {
		dbOnceSync[configName] = &sync.Once{}
	}
	dbOnceSync[configName].Do(func() {
		option := initOption(configName)
		var driver db.SqlClient
		switch option.Type {
		case "mysql":
			driver = &db.Mysql{}
			break
		case "mssql":
			driver = &db.Mssql{}
			break
		case "postgres":
			driver = &db.PostgresSQL{}
			break
		default:
			panic(fmt.Sprintf("database config's type:%s was not valid", option.Type))
		}
		Db := driver.Connect(option)
		Db.SingularTable(true)
		if option.Debug {
			Db.LogMode(true)
		}
		if option.Prefix != "" {
			gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
				return option.Prefix + defaultTableName
			}
		}
		dbInstance[configName] = Db
	})
	return dbInstance[configName]
}

// close db
func DbClose(dbConfigs ...string) {
	configName := "default"
	if len(dbConfigs) > 0 && dbConfigs[0] != "" {
		configName = dbConfigs[0]
	}
	if dbInstance[configName] != nil {
		_ = dbInstance[configName].Close()
		dbInstance[configName] = nil
	}
}
