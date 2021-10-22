package library
/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"strings"
)

// config object func interface
type Confer interface {
	GetAll() gin.H
	GetSection(baseKey string) gin.H
	Get(key string) *ini.Key
	Has(key string) bool
}

// factory init interface
type FactoryInterface interface {
	Source(file string) Confer
}

// config object type struct
type Config struct {
	init      *ini.File
	err       error
}

// check config has exists
func (c *Config) Has(key string) bool {
	if strings.Contains(key, ".") {
		keys := strings.Split(key, ".")
		var section *ini.Section
		var err error
		if section, err = c.init.GetSection(keys[0]); err != nil {
			return false
		}
		if arrLen := len(keys); arrLen < 2 {
			return false
		}
		return section.HasKey(keys[1])
	}
	return c.init.Section("").HasKey(key)
}

// get all base config
func (c *Config) GetSection(baseKey string) gin.H {
	conf, err := c.init.GetSection(baseKey)
	if err != nil {
		panic(fmt.Sprintf("Fail to get base config '%s' , error : %v", baseKey, err))
	}
	var configs = make(gin.H, 5)
	for _, key := range conf.KeyStrings() {
		configs[key] = conf.Key(key).Value()
	}
	return configs
}

// get all config
func (c *Config) GetAll() gin.H {
	var configs = make(gin.H, 0)
	for _, sec := range c.init.SectionStrings() {
		configs[sec] = c.GetSection(sec)
	}
	return configs
}

// get config
func (c *Config) Get(key string) *ini.Key {
	if strings.Contains(key, ".") {
		if !c.Has(key) {
			return &ini.Key{}
		}
		keys := strings.Split(key, ".")
		conf, _ := c.init.GetSection(keys[0])
		return conf.Key(keys[1])
	} else {
		return c.init.Section("").Key(key)
	}
}

type ConfigsFactory struct {
}

// get config file source
func (*ConfigsFactory) Source(file string) Confer {
	factory := &Config{}
	if !strings.Contains(file, "/") {
		file = "config/" + file
	}
	factory.init, factory.err = ini.Load(fmt.Sprintf("%s.ini", file))
	if factory.err != nil {
		panic(fmt.Sprintf("Fail to parse %s.ini , error : %v", file, factory.err))
	}
	return factory
}

// +++++++++++++++++  factory func +++++++++++++++++

// get source path
func callPath(source ...string) string {
	path := "app"
	if len(source) > 0 {
		path = source[0]
	}
	return path
}

// factory init config source
func SourceConfig(source string) Confer {
	factory := &ConfigsFactory{}
	return factory.Source(source)
}

// factory get config func
func GetConfig(key string, source ...string) *ini.Key {
	return SourceConfig(callPath(source...)).Get(key)
}

// factory has config func
func HasConfig(key string, source ...string) bool {
	return SourceConfig(callPath(source...)).Has(key)
}

// factory get all config func
func GetAllConfig(source ...string) gin.H {
	return SourceConfig(callPath(source...)).GetAll()
}

// factory get section func
func GetSectionConfig(baseKey string, source ...string) gin.H {
	return SourceConfig(callPath(source...)).GetSection(baseKey)
}
