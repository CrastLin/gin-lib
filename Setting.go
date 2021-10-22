package library
/**
 @auth CrastGin
 @date 2020-10
 */
import (
	"strings"
	"sync"
	"time"
)

var (
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ClientIp     string
	globalData   map[string]interface{}
	mut          sync.RWMutex
)

func init() {
	config := SourceConfig("app")
	RunMode = config.Get("run_mode").MustString("debug")
	HttpPort = config.Get("server.http_port").MustInt(800)
	ReadTimeout = time.Duration(config.Get("server.read_timeout").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(config.Get("server.write_timeout").MustInt(60)) * time.Second
	globalData = make(map[string]interface{}, 0)
}

// Set global config value
func SetGlobal(key string, value interface{}) bool {
	mut.Lock()
	defer mut.Unlock()
	globalData[key] = value
	return true
}

// set global config value if not exists
func SetGlobalOrExists(key string, value interface{}) bool {
	mut.Lock()
	defer mut.Unlock()
	if _, ok := globalData[key]; ok {
		return false
	}
	return SetGlobal(key, value)
}

// get global config value
func GetGlobal(key string) interface{} {
	mut.RLock()
	defer mut.RUnlock()
	keys := strings.Split(key, ".")
	var set interface{}
	var ok bool
	if set, ok = globalData[keys[0]]; !ok {
		return nil
	}
	if da, ok := set.(map[string]interface{}); len(keys) == 2 && ok {
		if set, ok = da[keys[1]]; !ok {
			return nil
		}
	}
	return set
}
